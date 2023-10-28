package discord

type Selfbot struct {
	Token string
	User  User
}

type Hello struct {
	Op int `json:"op"`
	D  struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}
type Resume struct {
	Op int        `json:"op"`
	D  ResumeData `json:"d"`
}
type ResumeData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

type Identify struct {
	Op int64        `json:"op"`
	D  IdentifyData `json:"d"`
}

type IdentifyData struct {
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
	Status     string        `json:"status"`
	Since      int64         `json:"since"`
	Activities []interface{} `json:"activities"`
	Afk        bool          `json:"afk"`
}

type SuperProperties struct {
	OS                     string      `json:"os"`
	Browser                string      `json:"browser"`
	Device                 string      `json:"device"`
	SystemLocale           string      `json:"system_locale"`
	BrowserUserAgent       string      `json:"browser_user_agent"`
	BrowserVersion         string      `json:"browser_version"`
	OSVersion              string      `json:"os_version"`
	Referrer               string      `json:"referrer"`
	ReferringDomain        string      `json:"referring_domain"`
	ReferrerCurrent        string      `json:"referrer_current"`
	ReferringDomainCurrent string      `json:"referring_domain_current"`
	ReleaseChannel         string      `json:"release_channel"`
	ClientBuildNumber      int64       `json:"client_build_number"`
	ClientEventSource      interface{} `json:"client_event_source"`
}

type DefaultEvent struct {
	Op int    `json:"op"`
	T  string `json:"t,omitempty"`
}
type Ready struct {
	Op int64     `json:"op"`
	D  ReadyData `json:"d"`
}

type ReadyData struct {
	Version          int     `json:"v"`
	User             User    `json:"user"`
	Guilds           []Guild `json:"guilds"`
	SessionID        string  `json:"session_id"`
	ResumeGatewayURL string  `json:"resume_gateway_url"`
}
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot,omitempty"`
	System        bool   `json:"system,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
	Banner        string `json:"banner,omitempty"`
	AccentColor   int    `json:"accent_color,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Verified      bool   `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         uint64 `json:"flag,omitempty"`
	PremiumType   uint64 `json:"premium_type,omitempty"`
	PublicFlags   uint64 `json:"public_flag,omitempty"`
}

type Guild struct {
	ID                          string        `json:"id"`
	Name                        string        `json:"name"`
	Icon                        string        `json:"icon"`
	IconHash                    string        `json:"icon_hash,omitempty"`
	Splash                      string        `json:"splash"`
	DiscoverySplash             string        `json:"discovery_splash"`
	Owner                       bool          `json:"owner,omitempty"`
	OwnerID                     string        `json:"owner_id"`
	Permissions                 string        `json:"permissions,omitempty"`
	Region                      string        `json:"region,omitempty"`
	AfkChannelID                string        `json:"afk_channel_id"`
	AfkTimeout                  int           `json:"afk_timeout"`
	WidgetEnabled               bool          `json:"widget_enabled,omitempty"`
	WidgetChannelID             string        `json:"widget_channel_id,omitempty"`
	VerificationLevel           uint64        `json:"verification_level"`
	DefaultMessageNotifications uint64        `json:"default_message_notifications"`
	ExplicitContentFilter       uint64        `json:"explicit_content_filter"`
	Roles                       []Role        `json:"roles"`
	Emojis                      []Emoji       `json:"emojis"`
	Features                    []string      `json:"features"`
	MFALevel                    uint64        `json:"mfa_level"`
	ApplicationID               string        `json:"application_id"`
	SystemChannelID             string        `json:"system_channel_id"`
	SystemChannelFlags          uint64        `json:"system_channel_flags"`
	RulesChannelID              string        `json:"rules_channel_id"`
	MaxPresences                int           `json:"max_presences,omitempty"`
	MaxMembers                  int           `json:"max_members,omitempty"`
	VanityUrl                   string        `json:"vanity_url_code"`
	Description                 string        `json:"description"`
	Banner                      string        `json:"banner"`
	PremiumTier                 uint64        `json:"premium_tier"`
	PremiumSubscriptionCount    int           `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string        `json:"preferred_locale"`
	PublicUpdatesChannelID      string        `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int           `json:"max_video_channel_users,omitempty"`
	ApproximateMemberCount      int           `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int           `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               WelcomeScreen `json:"welcome_screen,omitempty"`
	NSFWLevel                   uint64        `json:"nsfw_level"`
	Stickers                    []Sticker     `json:"stickers,omitempty"`
	PremiumProgressBarEnabled   bool          `json:"premium_progress_bar_enabled"`

	// Unavailable Guild Object
	// https://discord.com/developers/docs/resources/guild#unavailable-guild-object
	Unavailable bool `json:"unavailable,omitempty"`
}

// Role Object
// https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Color        int      `json:"color"`
	Hoist        bool     `json:"hoist"`
	Icon         string   `json:"icon,omitempty"`
	UnicodeEmoji string   `json:"unicode_emoji,omitempty"`
	Position     int      `json:"position"`
	Permissions  string   `json:"permissions"`
	Managed      bool     `json:"managed"`
	Mentionable  bool     `json:"mentionable"`
	Tags         RoleTags `json:"tags,omitempty"`
}

// Role Tags Structure
// https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID             string `json:"bot_id,omitempty"`
	IntegrationID     string `json:"integration_id,omitempty"`
	PremiumSubscriber bool   `json:"premium_subscriber,omitempty"`
}

type Emoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	User          User     `json:"user,omitempty"`
	RequireColons bool     `json:"require_colons,omitempty"`
	Managed       bool     `json:"managed,omitempty"`
	Animated      bool     `json:"animated,omitempty"`
	Available     bool     `json:"available,omitempty"`
}

type WelcomeScreen struct {
	Description           string                 `json:"description"`
	WelcomeScreenChannels []WelcomeScreenChannel `json:"welcome_channels"`
}

// Welcome Screen Channel Structure
// https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-channel-structure
type WelcomeScreenChannel struct {
	ChannelID   string `json:"channel_id"`
	Description string `json:"description"`
	EmojiID     string `json:"emoji_id"`
	EmojiName   string `json:"emoji_name"`
}

// Sticker Structure
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-structure
type Sticker struct {
	ID          string `json:"id"`
	PackID      string `json:"pack_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Asset       string `json:"asset,omitempty"`
	Type        uint64 `json:"type"`
	FormatType  uint64 `json:"format_type"`
	Available   bool   `json:"available,omitempty"`
	GuildID     string `json:"guild_id,omitempty"`
	User        User   `json:"user,omitempty"`
	SortValue   int    `json:"sort_value,omitempty"`
}

type Message struct {
	Op int         `json:"op"`
	D  MessageData `json:"d"`
}

type MessageData struct {
	Content    string             `json:"content,omitempty"`
	ChannelID  string             `json:"channel_id,omitempty"`
	Embeds     []Embed            `json:"embeds,omitempty"`
	Reactions  []Reaction         `json:"reactions,omitempty"`
	Author     User               `json:"author,omitempty"`
	GuildID    string             `json:"guild_id,omitempty"`
	MessageId  string             `json:"id,omitempty"`
	Components []MessageComponent `json:"components,omitempty"`
	Flags      int                `json:"flags,omitempty"`
}
type Embed struct {
	Title string `json:"title,omitempty"`

	// The type of embed. Always EmbedTypeRich for webhook embeds.
	Type        string             `json:"type,omitempty"`
	Description string             `json:"description,omitempty"`
	URL         string             `json:"url,omitempty"`
	Image       *MessageEmbedImage `json:"image,omitempty"`

	// The color code of the embed.
	Color     int                    `json:"color,omitempty"`
	Footer    EmbedFooter            `json:"footer,omitempty"`
	Thumbnail *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Provider  EmbedProvider          `json:"provider,omitempty"`
	Author    EmbedAuthor            `json:"author,omitempty"`
	Fields    []EmbedField           `json:"fields,omitempty"`
}
type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}
type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Reaction struct {
	Emojis Emoji `json:"emoji,omitempty"`
	Count  int   `json:"count,omitempty"`
}
type MessageComponent struct {
	Type    int       `json:"type"`
	Buttons []Buttons `json:"components"`
}
type Buttons struct {
	Type     int         `json:"type,omitempty"`
	Style    int         `json:"style,omitempty"`
	Label    string      `json:"label,omitempty"`
	CustomID string      `json:"custom_id,omitempty"`
	Hash     string      `json:"hash,omitempty"`
	Emoji    ButtonEmoji `json:"emoji,omitempty"`
	Disabled bool        `json:"disabled,omitempty"`
}

type ButtonEmoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}
