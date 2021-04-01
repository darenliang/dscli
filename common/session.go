package common

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"strconv"
)

// GetDiscordSession gets the Discord session, guild and channels
func GetDiscordSession() (*discordgo.Session, *discordgo.Guild, []*discordgo.Channel, error) {
	token, ok := viper.Get("token").(string)
	if !ok {
		return nil, nil, nil, errors.New("invalid Discord user token config")
	}

	session, err := discordgo.New(token)
	if err != nil {
		return nil, nil, nil, errors.New("discord session: " + err.Error())
	}

	// monkey patch client timeout to hang indefinitely
	session.Client.Timeout = 0

	err = session.Open()
	if err != nil {
		return nil, nil, nil, errors.New("discord session: " + err.Error())
	}

	id, ok := viper.Get("id").(string)
	if !ok {
		session.Close()
		return nil, nil, nil, errors.New("invalid server id config")
	}

	guild, err := session.Guild(id)
	if err != nil {
		session.Close()
		return nil, nil, nil, errors.New("discord server: " + err.Error())
	}

	channels, err := session.GuildChannels(guild.ID)
	if err != nil {
		session.Close()
		return nil, nil, nil, errors.New("discord channels: " + err.Error())
	}

	return session, guild, channels, nil
}

// GetMaxFileSizeUpload gets the maximum file size that can be uploaded
func GetMaxFileSizeUpload(session *discordgo.Session, guild *discordgo.Guild) (int, error) {
	user, err := session.User("@me")
	if err != nil {
		return 0, err
	}

	maxSizeUser := MaxDiscordFileSizeDefault

	// update max size if user is not a bot
	if !user.Bot {
		switch user.PremiumType {
		case 0: // do nothing
		case 1:
			maxSizeUser = MaxDiscordFileSizeNitroClassic
		case 2:
			maxSizeUser = MaxDiscordFileSizeNitro
		default:
			return 0, errors.New("invalid user premium type: " + strconv.Itoa(user.PremiumType))
		}
	} else {
		// return default if user is a bot
		return MaxDiscordFileSizeDefault, nil
	}

	maxSizeGuild := MaxDiscordFileSizeDefault

	// update max size if guild is a premium tier
	switch guild.PremiumTier {
	case 0, 1: // do nothing
	case 2:
		maxSizeGuild = MaxDiscordFileSizeNitroClassic
	case 3:
		maxSizeGuild = MaxDiscordFileSizeNitro
	default:
		return 0, errors.New("invalid guild premium tier: " + strconv.Itoa(int(guild.PremiumTier)))
	}

	if maxSizeUser > maxSizeGuild {
		return maxSizeUser, nil
	}

	return maxSizeGuild, nil
}
