package common

// All constants below maybe subject to change.

const (
	// Version is the current dscli app version
	Version = "1.13.0"

	// MaxDiscordFileSizeDefault represents the maximum file size in bytes that
	// an attachment can be for non nitro users
	MaxDiscordFileSizeDefault = 26214400

	// MaxDiscordFileSizeNitroClassic represents the maximum file size in bytes that
	// an attachment can be for nitro classic users
	MaxDiscordFileSizeNitroClassic = 52428800

	// MaxDiscordFileSizeNitro represents the maximum file size in bytes that
	// an attachment can be for nitro users
	MaxDiscordFileSizeNitro = 524288000

	// MaxDiscordMessageRequest is the maximum number of messages you can query
	MaxDiscordMessageRequest = 100

	// MaxDiscordChannels is the maximum number of channels you can have
	MaxDiscordChannels = 500
)
