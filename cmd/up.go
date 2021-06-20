package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/dscli/common"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

var upDebug bool

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:        "up <local file> <remote file>",
	Example:    "up test.txt test.txt",
	SuggestFor: []string{"upload"},
	Short:      "Upload file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one argument")
		}
		return nil
	},
	RunE: up,
}

func init() {
	upCmd.Flags().BoolVarP(&upDebug, "debug", "d", false, "debug mode: <total bytes> <bytes uploaded>")

	rootCmd.AddCommand(upCmd)
}

// up command handler
func up(cmd *cobra.Command, args []string) error {
	session, guild, channels, err := common.GetDiscordSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fileMap, err := common.ParseFileMap(channels)
	if err != nil {
		return err
	}

	local := args[0] // local filename

	// open local file to upload
	localFile, err := os.Open(local)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, localBase := filepath.Split(local)

	var remote string // remote filename
	if len(args) == 1 {
		remote = localBase
	} else {
		remote = args[1]
	}

	// remote filename already exists
	if _, ok := fileMap[remote]; ok {
		return errors.New(remote + " already exists on Discord")
	}

	// encode remote filename
	encodedRemote, err := common.EncodeFilename(remote)
	if err != nil {
		return err
	}

	// create channel for file
	channel, err := session.GuildChannelCreate(guild.ID, encodedRemote, discordgo.ChannelTypeGuildText)
	if err != nil {
		return errors.New("cannot create remote file: " + err.Error())
	}

	// get max Discord file size
	maxDiscordFileSize, err := common.GetMaxFileSizeUpload(session, guild)
	if err != nil {
		return err
	}

	// setup buffer with max Discord file size
	buf := make([]byte, maxDiscordFileSize)

	// get size of local file
	stat, err := localFile.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	sizeStr := strconv.Itoa(int(size))

	// set channel topic to filesize
	channelSettings := &discordgo.ChannelEdit{
		Topic: sizeStr,
	}
	// ignore if errored since it is not critical
	_, _ = session.ChannelEditComplex(channel.ID, channelSettings)

	var bar *progressbar.ProgressBar

	if !upDebug {
		// init progress bar
		bar = progressbar.DefaultBytes(
			size,
			"Uploading "+localBase,
		)
	}

	first := true
	part := 1

	for {
		// read chunk
		length, err := localFile.Read(buf)
		if err != nil {
			return err
		}

		var reader io.Reader

		if !upDebug {
			reader = io.TeeReader(bytes.NewReader(buf[:length]), bar)
		} else {
			reader = bytes.NewReader(buf[:length])
		}

		msg := &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   strconv.Itoa(part),
					Reader: reader,
				},
			},
		}

		part += 1

		// send chunk
		// retry 5 times because internet can be flaky and Discord sometimes
		// likes to drop connections
		var message *discordgo.Message
		maxUploadTries := 5
		for i := 0; i < maxUploadTries; i++ {
			message, err = session.ChannelMessageSendComplex(channel.ID, msg)
			if err != nil {
				if i == maxUploadTries-1 {
					return err
				} else {
					continue
				}
			}
			break
		}

		if upDebug {
			offset, err := localFile.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
			fmt.Printf("%d %d \n", size, offset)
		}

		if first {
			// if pin fails, ignore
			_ = session.ChannelMessagePin(message.ChannelID, message.ID)
			first = false
		}

		// exit loop if EOF
		if length < maxDiscordFileSize {
			break
		}
	}

	return nil
}
