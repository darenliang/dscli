package common

import (
	"encoding/base32"
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// DecodeFilename decodes filename to the original filename
func DecodeFilename(filename string) (string, error) {
	// go's base32 decoder only likes uppercase characters
	filename = strings.ToUpper(filename)

	// add padding if necessary
	padding := ""
	switch len(filename) % 8 {
	case 2:
		padding = "======"
	case 4:
		padding = "===="
	case 5:
		padding = "==="
	case 7:
		padding = "=="
	}
	filename += padding

	// decode string
	bytes, err := base32.StdEncoding.DecodeString(filename)
	if err != nil {
		return "", errors.New("decode filename: " + err.Error())
	}

	return string(bytes), nil
}

// EncodeFilename encodes filename to Discord-safe channel name
func EncodeFilename(filename string) (string, error) {
	// encode filename, trim right "=" characters and change to lowercase
	encodedFilename := strings.ToLower(strings.TrimRight(base32.StdEncoding.EncodeToString([]byte(filename)), "="))

	// check if filename is invalid
	if len(encodedFilename) < 1 {
		return "", errors.New("filename must contain at least one character")
	}
	if len(encodedFilename) > 100 {
		return "", errors.New(filename + " is too long")
	}

	return encodedFilename, nil
}

// ParseFileMap creates channel name to channel struct mapping
func ParseFileMap(channels []*discordgo.Channel) (map[string]*discordgo.Channel, error) {
	fileMap := make(map[string]*discordgo.Channel)

	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			filename, err := DecodeFilename(channel.Name)
			if err != nil {
				return nil, err
			}
			fileMap[filename] = channel
		}
	}

	return fileMap, nil
}
