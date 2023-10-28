package discord

const (
	API_VERSION     = "9"
	BROWSER         = "Firefox"
	BROWSER_VERSION = "111.0"
	CAPABILITIES    = 4093
	DEVICE          = "" // Discord's official client sends an empty string.
	OS              = "Windows"
	OS_VERSION      = "10"
	STATUS          = "offline" // https://discord.com/developers/docs/topics/gateway-events#update-presence-status-types
	USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/" + BROWSER_VERSION
	BUILD_NUMBER    = 240884
)
