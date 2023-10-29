package discord

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/goccy/go-json"
	"github.com/luminaldev/dgate/types"
)

var (
	gatewayURL = "wss://gateway.discord.gg/?encoding=json&v=" + API_VERSION
	headers    = make(http.Header)
)

func init() {
	if len(headers) == 0 {
		headers.Set("Host", "gateway.discord.gg")
		headers.Set("User-Agent", USER_AGENT)
	}
}

type Gateway struct {
	CloseChan         chan struct{}
	Connection        *websocket.Conn
	Selfbot           *Selfbot
	LastSeq           int
	SessionID         string
	heartbeatInterval time.Duration
	GatewayURL        string
	Handlers          Handlers
}

func CreateGateway(selfbot *Selfbot) *Gateway {
	return &Gateway{CloseChan: make(chan struct{}), Selfbot: selfbot, GatewayURL: gatewayURL}
}

func (gateway *Gateway) Connect() error {
	conn, resp, err := websocket.DefaultDialer.Dial(gateway.GatewayURL, headers)

	if resp.StatusCode == 404 {
		return fmt.Errorf("gateway not found")
	} else if err != nil {
		return err
	}

	gateway.Connection = conn

	if err = gateway.hello(); err != nil {
		return err
	}

	if err = gateway.identify(); err != nil {
		return err
	}

	if err = gateway.ready(); err != nil {
		return err
	}

	gateway.startHandler()
	return nil
}

func (gateway *Gateway) hello() error {
	msg, err := gateway.readMessage()

	if err != nil {
		return err
	}

	var resp types.HelloEvent

	if err = json.Unmarshal(msg, &resp); err != nil {
		return err
	}

	if resp.Op != types.OpcodeHello {
		return fmt.Errorf("unexpected opcode, expected %d, got %d", types.OpcodeHello, resp.Op)
	}

	gateway.heartbeatInterval = time.Duration(resp.D.HeartbeatInterval)
	go gateway.heartbeatSender()

	return nil
}

func (gateway *Gateway) identify() error {
	payload, err := json.Marshal(types.ResumePayload{
		Op: types.OpcodeResume,
		D: types.ResumePayloadData{
			Token:     gateway.Selfbot.Token,
			SessionID: gateway.SessionID,
			Seq:       gateway.LastSeq,
		},
	})

	if err != nil {
		return err
	}

	if gateway.canReconnect() {
		err := gateway.sendMessage(payload)

		if err != nil {
			return err
		}
	} else {
		payload, err := json.Marshal(types.IdentifyPayload{
			Op: types.OpcodeIdentify,
			D: types.IdentifyPayloadData{
				Token:        gateway.Selfbot.Token,
				Capabilities: CAPABILITIES,
				Properties: types.SuperProperties{
					OS:                     OS,
					Browser:                BROWSER,
					Device:                 DEVICE,
					SystemLocale:           clientLocale,
					BrowserUserAgent:       USER_AGENT,
					BrowserVersion:         BROWSER_VERSION,
					OSVersion:              OS_VERSION,
					Referrer:               "",
					ReferringDomain:        "",
					ReferrerCurrent:        "",
					ReferringDomainCurrent: "",
					ReleaseChannel:         "stable",
					ClientBuildNumber:      clientBuildNumber,
					ClientEventSource:      nil,
				},
				Presence: types.Presence{
					Status:     STATUS,
					Since:      0,
					Activities: nil,
					Afk:        false,
				},
				Compress: false,
				ClientState: types.ClientState{
					GuildVersions:            types.GuildVersions{},
					HighestLastMessageID:     "0",
					ReadStateVersion:         0,
					UserGuildSettingsVersion: -1,
					UserSettingsVersion:      -1,
					PrivateChannelsVersion:   "0",
					APICodeVersion:           0,
				},
			},
		})

		if err != nil {
			return err
		}

		err = gateway.sendMessage(payload)

		if err != nil {
			return err
		}
	}

	return nil
}

func (gateway *Gateway) ready() error {
	msg, err := gateway.readMessage()

	if err != nil {
		return err
	}

	var event types.DefaultEvent
	err = json.Unmarshal(msg, &event)

	if err != nil {
		return err
	}

	if event.Op == types.OpcodeInvalidSession {
		<-gateway.CloseChan
		return gateway.reconnect()
	} else if event.Op != types.OpcodeDispatch {
		return fmt.Errorf("unexpected opcode, expected %d, got %d", types.OpcodeDispatch, event.Op)
	}

	var ready types.ReadyEvent

	if err = json.Unmarshal(msg, &ready); err != nil {
		return err
	}

	gateway.Selfbot.User = ready.D.User
	gateway.SessionID = ready.D.SessionID
	gateway.GatewayURL = ready.D.ResumeGatewayURL

	for _, handler := range gateway.Handlers.OnReady {
		handler(&ready.D)
	}

	return nil
}

func (gateway *Gateway) canReconnect() bool {
	return gateway.SessionID != "" && gateway.LastSeq != 0 && gateway.GatewayURL != ""
}

func (gateway *Gateway) heartbeatSender() {
	ticker := time.NewTicker(gateway.heartbeatInterval * time.Millisecond) // Every heartbeat interval (sent in milliseconds).
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // On ticker tick.
			if err := gateway.sendHeartbeat(); err != nil {
				return
			}
		case <-gateway.CloseChan: // If a close is signalled.
			return
		default:
			time.Sleep(25 * time.Millisecond) // Wait to prevent CPU overload.
		}
	}
}

func (gateway *Gateway) sendHeartbeat() error {
	payload, err := json.Marshal(types.DefaultEvent{
		Op: types.OpcodeHeartbeat,
		D:  4,
	})

	if err != nil {
		return err
	}

	return gateway.sendMessage(payload)

}
func (gateway *Gateway) sendMessage(payload []byte) error {
	err := gateway.Connection.WriteMessage(websocket.TextMessage, payload)

	if err != nil {
		var closeError *websocket.CloseError
		errors.As(err, &closeError)

		switch closeError.Code {
		case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived:
			go gateway.reset()
			return err
		default:
			closeEvent, ok := types.CloseEventCodes[closeError.Code]

			if ok && closeEvent.Reconnect {
				go gateway.reconnect()
			}

			return fmt.Errorf("gateway closed with code %d: %s - %s", closeEvent.Code, closeEvent.Description, closeEvent.Explanation)
		}
	}
	return nil
}

func (gateway *Gateway) readMessage() ([]byte, error) {
	if gateway.Connection == nil {
		return nil, errors.New("gateway connection is already closed")
	}

	_, msg, err := gateway.Connection.ReadMessage()

	if err != nil {
		closeError := err.(*websocket.CloseError)

		switch closeError.Code {
		case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived: // Websocket closed without any close code.
			go gateway.reset()
			return nil, err
		default:
			if closeEvent, ok := types.CloseEventCodes[closeError.Code]; ok {
				if closeEvent.Reconnect { // If the session is re-connectable.
					go gateway.reconnect()
				} else {
					gateway.Close()

					return nil, fmt.Errorf("gateway closed with code %d: %s - %s", closeEvent.Code, closeEvent.Description, closeEvent.Explanation)
				}
			} else {
				gateway.Close()

				return nil, err
			}
		}
	}

	return msg, nil
}

func (gateway *Gateway) reset() error {
	gateway.LastSeq = 0
	gateway.SessionID = ""

	return gateway.reconnect()
}

func (gateway *Gateway) reconnect() error {
	return gateway.Connect()
}

func (gateway *Gateway) startHandler() {
	for {
		select {
		case <-gateway.CloseChan:
			return
		default:
			msg, err := gateway.readMessage()

			if err != nil {
				fmt.Println(err)
				return
			}

			var event types.DefaultEvent

			if err = json.Unmarshal(msg, &event); err != nil {
				return
			}

			switch event.Op {
			case types.OpcodeDispatch:
				var data types.MessageEvent

				err := json.Unmarshal(msg, &data)

				if err != nil {
					continue
				}

				switch event.T {
				case types.EventNameMessageCreate:
					for _, handler := range gateway.Handlers.OnMessageCreate {
						handler(&data.D)
					}
				case types.EventNameMessageUpdate:
					for _, handler := range gateway.Handlers.OnMessageUpdate {
						handler(&data.D)
					}
				}
			case types.OpcodeHeartbeat:
				gateway.sendHeartbeat()
			case types.OpcodeHeartbeatACK:
				continue
			case types.OpcodeReconnect:
				gateway.reconnect()

				for _, handler := range gateway.Handlers.OnReconnect {
					handler()
				}
			case types.OpcodeInvalidSession:
				gateway.reconnect()

				for _, handler := range gateway.Handlers.OnInvalidated {
					handler()
				}
			}

			if event.S == 0 { // Some payloads, for example the heartbeat ack, don't contribute to the sequence ID.
				gateway.LastSeq = event.S
			}
		}
	}
}

func (gateway *Gateway) Close() {
	gateway.CloseChan <- struct{}{}

	err := gateway.Connection.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "going away"), time.Now().Add(time.Second*10))

	if err != nil {
		if gateway.Connection != nil {
			gateway.Connection.Close()
			gateway.Connection = nil
		}
	}
}
