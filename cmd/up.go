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

	// get max discord file size
	maxDiscordFileSize, err := common.GetMaxFileSizeUpload(session, guild)
	if err != nil {
		return nil
	}

	// setup buffer with max discord file size
	buf := make([]byte, maxDiscordFileSize)

	// get size of local file
	stat, err := localFile.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	sizeStr := strconv.Itoa(int(size))

	var bar *progressbar.ProgressBar

	if !upDebug {
		// init progress bar
		bar = progressbar.DefaultBytes(
			size,
			"Uploading "+localBase,
		)
	}

	first := true
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
					Name:   sizeStr,
					Reader: reader,
				},
			},
		}

		// send chunk
		message, err := session.ChannelMessageSendComplex(channel.ID, msg)
		if err != nil {
			return err
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
