package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/luminaldev/dgate/types"
	"net/http"
	"reflect"
	"time"
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

func (g *Gateway) Connect() error {
	conn, resp, err := websocket.DefaultDialer.Dial(g.GatewayURL, headers)
	if resp.StatusCode == 404 {
		return fmt.Errorf("gateway not found")
	} else if err != nil {
		return err
	}
	g.Connection = conn

	if err = g.hello(); err != nil {
		return err
	}
	if err = g.identify(); err != nil {
		return err
	}
	if err = g.ready(); err != nil {
		return err
	}
	g.startHandler()
	return nil
}

func (g *Gateway) hello() error {
	msg, err := g.readMessage()
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
	g.heartbeatInterval = time.Duration(resp.D.HeartbeatInterval)
	go g.heartbeatSender()
	return nil
}

func (g *Gateway) identify() error {
	payload, err := json.Marshal(types.ResumePayload{
		Op: types.OpcodeResume,
		D: types.ResumePayloadData{
			Token:     g.Selfbot.Token,
			SessionID: g.SessionID,
			Seq:       g.LastSeq,
		},
	})
	if err != nil {
		return err
	}
	if g.canReconnect() {
		err := g.sendMessage(payload)
		if err != nil {
			return err
		}
	} else {
		payload, err := json.Marshal(types.IdentifyPayload{
			Op: types.OpcodeIdentify,
			D: types.IdentifyPayloadData{
				Token:        g.Selfbot.Token,
				Capabilities: CAPABILITIES,
				Properties: types.SuperProperties{
					OS:                     OS,
					Browser:                BROWSER,
					Device:                 DEVICE,
					SystemLocale:           "en-us",
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
		err = g.sendMessage(payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gateway) ready() error {
	msg, err := g.readMessage()
	if err != nil {
		return err
	}
	var event types.DefaultEvent
	err = json.Unmarshal(msg, &event)
	if err != nil {
		return err
	}
	if event.Op == types.OpcodeInvalidSession {
		<-g.CloseChan
		return g.reconnect()
	} else if event.Op != types.OpcodeDispatch {
		return fmt.Errorf("unexpected opcode, expected %d, got %d", types.OpcodeDispatch, event.Op)
	}
	var ready types.ReadyEvent
	if err = json.Unmarshal(msg, &ready); err != nil {
		return err
	}
	g.Selfbot.User = ready.D.User
	g.SessionID = ready.D.SessionID
	g.GatewayURL = ready.D.ResumeGatewayURL
	for _, handler := range g.Handlers.OnReady {
		handler(&ready.D)
	}
	return nil
}

func (g *Gateway) canReconnect() bool {
	return g.SessionID != "" && g.LastSeq != 0 && g.GatewayURL != ""
}
func (g *Gateway) heartbeatSender() {
	ticker := time.NewTicker(g.heartbeatInterval * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := g.sendHeartbeat(); err != nil {
				return
			}
		case <-g.CloseChan:
			return
		default:
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (g *Gateway) sendHeartbeat() error {
	payload, err := json.Marshal(types.DefaultEvent{
		Op: types.OpcodeHeartbeat,
	})
	if err != nil {
		return err
	}
	return g.sendMessage(payload)

}
func (g *Gateway) sendMessage(payload []byte) error {
	err := g.Connection.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		var closeError *websocket.CloseError
		errors.As(err, &closeError)

		switch closeError.Code {
		case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived:
			go g.reset()
			return err
		default:
			closeEvent, ok := types.CloseEventCodes[closeError.Code]
			if ok && closeEvent.Reconnect {
				go g.reconnect()
			}
			return fmt.Errorf("gateway closed with code %d: %s - %s", closeEvent.Code, closeEvent.Description, closeEvent.Explanation)
		}
	}
	return nil
}

func (g *Gateway) readMessage() ([]byte, error) {
	_, msg, err := g.Connection.ReadMessage()
	if err != nil {
		var closeError *websocket.CloseError
		errors.As(err, &closeError)

		switch closeError.Code {
		case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived: // Websocket closed without any close code.
			go g.reset()
			return nil, err
		default:
			if closeEvent, ok := types.CloseEventCodes[closeError.Code]; ok {
				if closeEvent.Reconnect { // If the session is reconnectable.
					go g.reconnect()
				}
				return nil, fmt.Errorf("gateway closed with code %d: %s - %s", closeEvent.Code, closeEvent.Description, closeEvent.Explanation)
			} else {
				return nil, err
			}
		}
	}
	return msg, nil
}

func (g *Gateway) reset() error {
	g.LastSeq = 0
	g.SessionID = ""

	return g.reconnect()
}

func (g *Gateway) reconnect() error {
	return g.Connect()
}

func (g *Gateway) startHandler() {
	for {
		select {
		case <-g.CloseChan:
			return
		default:
			msg, err := g.readMessage()
			if err != nil {
				return
			}
			var def types.DefaultEvent
			if err = json.Unmarshal(msg, &def); err != nil {
				return
			}
			switch def.Op {
			case types.OpcodeDispatch:
				switch def.T {
				case types.EventNameMessageCreate:
					var data types.MessageEvent
					err := json.Unmarshal(msg, &data)
					if err != nil {
						continue
					}
					for _, handler := range g.Handlers.OnMessageCreate {
						handler(&data.D)
					}
				case types.EventNameMessageUpdate:
					var data types.MessageEvent
					err := json.Unmarshal(msg, &data)
					if err != nil {
						continue
					}
					for _, handler := range g.Handlers.OnMessageUpdate {
						handler(&data.D)
					}

				}
			case types.OpcodeHeartbeat:
				g.sendHeartbeat()
			case types.OpcodeHeartbeatACK:
				continue
			case types.OpcodeReconnect:
				g.reconnect()
				for _, handler := range g.Handlers.OnReconnect {
					handler()
				}
			case types.OpcodeInvalidSession:
				g.reconnect()
				for _, handler := range g.Handlers.OnInvalidated {
					handler()
				}
			}
			if !reflect.ValueOf(def.S).IsZero() {
				g.LastSeq = def.S
			}
		}
	}
}
