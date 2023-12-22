package types

import "time"

type Opcode struct {
	Op int `json:"op"`
}

type HelloEvent struct {
	Opcode
	D struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}

type ResumePayload struct {
	Op int               `json:"op"`
	D  ResumePayloadData `json:"d"`
}

type ResumePayloadData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

type IdentifyPayload struct {
	Op int64               `json:"op"`
	D  IdentifyPayloadData `json:"d"`
}

type IdentifyPayloadData struct {
	Token        string          `json:"token"`
	Capabilities int64           `json:"capabilities"`
	Properties   SuperProperties `json:"properties"`
	Presence     Presence        `json:"presence"`
	Compress     bool            `json:"compress"`
	ClientState  ClientState     `json:"client_state"`
}

type ClientState struct {
	GuildVersions            GuildVersions `json:"guild_versions"`
	HighestLastMessageID     string        `json:"highest_last_message_id"`
	ReadStateVersion         int64         `json:"read_state_version"`
	UserGuildSettingsVersion int64         `json:"user_guild_settings_version"`
	UserSettingsVersion      int64         `json:"user_settings_version"`
	PrivateChannelsVersion   string        `json:"private_channels_version"`
	APICodeVersion           int64         `json:"api_code_version"`
}

type GuildVersions struct {
}

type Presence struct {
	Status     string `json:"status"`
	Since      int64  `json:"since"`
	Activities []any  `json:"activities"`
	Afk        bool   `json:"afk"`
}

type SuperProperties struct {
	OS                     string `json:"os"`
	Browser                string `json:"browser"`
	Device                 string `json:"device"`
	SystemLocale           string `json:"system_locale"`
	BrowserUserAgent       string `json:"browser_user_agent"`
	BrowserVersion         string `json:"browser_version"`
	OSVersion              string `json:"os_version"`
	Referrer               string `json:"referrer"`
	ReferringDomain        string `json:"referring_domain"`
	ReferrerCurrent        string `json:"referrer_current"`
	ReferringDomainCurrent string `json:"referring_domain_current"`
	ReleaseChannel         string `json:"release_channel"`
	ClientBuildNumber      string `json:"client_build_number"`
	ClientEventSource      any    `json:"client_event_source"`
}

// https://discord.com/developers/docs/topics/gateway-events#payload-structure
type DefaultEvent struct {
	Op int    `json:"op"`
	T  string `json:"t,omitempty"`
	S  int    `json:"s,omitempty"`
	D  any    `json:"d,omitempty"`
}

type ReadyEvent struct {
	Op int64          `json:"op"`
	D  ReadyEventData `json:"d"`
}

type ReadyEventData struct {
	Version          int     `json:"v"`
	User             User    `json:"user"`
	Guilds           []Guild `json:"guilds"`
	SessionID        string  `json:"session_id"`
	ResumeGatewayURL string  `json:"resume_gateway_url"`
}

type MessageEvent struct {
	Op int              `json:"op"`
	D  MessageEventData `json:"d"`
}
type MemberEvent struct {
	Op int                        `json:"op"`
	D  GuildMembersChunkEventData `json:"d"`
}

type MessageEventData struct {
	// Data is in different struct because it needs to be recursive
	MessageData
	ReferencedMessage MessageData `json:"referenced_message"`
}

type GuildMembersChunkEventData struct {
	Presences []struct {
		User struct {
			Id string `json:"id"`
		} `json:"user"`
		Status       string `json:"status"`
		ClientStatus struct {
			Web string `json:"web"`
		} `json:"client_status"`
		Broadcast  interface{}   `json:"broadcast"`
		Activities []interface{} `json:"activities"`
	} `json:"presences"`
	NotFound []interface{} `json:"not_found"`
	Members  []struct {
		User struct {
			Username             string      `json:"username"`
			PublicFlags          int         `json:"public_flags"`
			Id                   string      `json:"id"`
			GlobalName           string      `json:"global_name"`
			DisplayName          string      `json:"display_name"`
			Discriminator        string      `json:"discriminator"`
			Bot                  bool        `json:"bot"`
			AvatarDecorationData interface{} `json:"avatar_decoration_data"`
			Avatar               string      `json:"avatar"`
		} `json:"user"`
		Roles                      []interface{} `json:"roles"`
		PremiumSince               interface{}   `json:"premium_since"`
		Pending                    bool          `json:"pending"`
		Nick                       interface{}   `json:"nick"`
		Mute                       bool          `json:"mute"`
		JoinedAt                   time.Time     `json:"joined_at"`
		Flags                      int           `json:"flags"`
		Deaf                       bool          `json:"deaf"`
		CommunicationDisabledUntil interface{}   `json:"communication_disabled_until"`
		Avatar                     interface{}   `json:"avatar"`
	} `json:"members"`
	GuildId    string `json:"guild_id"`
	ChunkIndex int    `json:"chunk_index"`
	ChunkCount int    `json:"chunk_count"`
}
type MessageData struct {
	Type             int                `json:"type,omitempty"`
	Content          string             `json:"content,omitempty"`
	ChannelID        string             `json:"channel_id,omitempty"`
	Embeds           []Embed            `json:"embeds,omitempty"`
	Reactions        []Reaction         `json:"reactions,omitempty"`
	Author           User               `json:"author,omitempty"`
	GuildID          string             `json:"guild_id,omitempty"`
	MessageId        string             `json:"id,omitempty"`
	Components       []MessageComponent `json:"components,omitempty"`
	Attachments      []Attachment       `json:"attachments,omitempty"`
	Flags            int                `json:"flags,omitempty"`
	MessageReference MessageReference   `json:"message_reference,omitempty"`
}

type Attachment struct {
	ID                 string `json:"id"`
	Filename           string `json:"filename"`
	Size               int    `json:"size"`
	URL                string `json:"url"`
	ProxyURL           string `json:"proxy_url"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	ContentType        string `json:"content_type"`
	Placeholder        string `json:"placeholder"`
	PlaceholderVersion int    `json:"placeholder_version"`
}
type MessageReference struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id"`
}
