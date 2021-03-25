package common

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
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
