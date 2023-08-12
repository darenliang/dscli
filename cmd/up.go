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

var (
	upDebug  bool
	upResume bool
)

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
	upCmd.Flags().BoolVarP(&upResume, "resume", "r", false, "resume upload")

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

	// check if max Discord channel limit is reached
	if !upResume && len(fileMap) >= common.MaxDiscordChannels {
		return errors.New(
			fmt.Sprintf(
				"max Discord channel limit of %d is reached",
				common.MaxDiscordChannels,
			),
		)
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
	if _, ok := fileMap[remote]; ok && !upResume {
		return errors.New(remote + " already exists on Discord")
	} else if !ok && upResume {
		return errors.New(remote + " does not exist on Discord")
	}

	// get size of local file
	stat, err := localFile.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	sizeStr := strconv.FormatInt(size, 10)

	// get max Discord file size
	maxDiscordFileSize, err := common.GetMaxFileSizeUpload(session, guild)
	if err != nil {
		return err
	}

	var channel *discordgo.Channel
	blockNumber := 0

	if upResume {
		channel = fileMap[remote]

		if channel.Topic != sizeStr {
			return errors.New("remote file size does not match local file size")
		}

		msgs, err := session.ChannelMessages(channel.ID, 1, "", "0", "")
		if err != nil {
			return err
		}

		if len(msgs) == 0 || len(msgs[0].Attachments) == 0 {
			return errors.New("cannot infer block size")
		}

		if msgs[0].Attachments[0].Size > maxDiscordFileSize {
			return errors.New(fmt.Sprintf(
				"inferred block size %d is larger than the largest permitted block size %d",
				msgs[0].Attachments[0].Size,
				maxDiscordFileSize,
			))
		}

		maxDiscordFileSize = msgs[0].Attachments[0].Size

		msgs, err = session.ChannelMessages(channel.ID, 2, "", "", "")
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			if len(msg.Attachments) == 0 {
				continue
			}
			if msg.Attachments[0].Size != maxDiscordFileSize {
				return errors.New("complete upload inferred from incomplete last block")
			}
			blockNumber, err = strconv.Atoi(msg.Attachments[0].Filename)
			if err != nil {
				return err
			}
			break
		}

		if int64(blockNumber)*int64(maxDiscordFileSize) == size {
			return errors.New("upload is already complete")
		}

	} else {
		// encode remote filename
		encodedRemote, err := common.EncodeFilename(remote)
		if err != nil {
			return err
		}

		// create channel for file
		channel, err = session.GuildChannelCreate(guild.ID, encodedRemote, discordgo.ChannelTypeGuildText)
		if err != nil {
			return errors.New("cannot create remote file: " + err.Error())
		}

		// set channel topic to filesize
		channelSettings := &discordgo.ChannelEdit{
			Topic: sizeStr,
		}
		// ignore if errored since it is not critical
		_, _ = session.ChannelEdit(channel.ID, channelSettings)
	}

	// seek to block number
	_, err = localFile.Seek(int64(blockNumber)*int64(maxDiscordFileSize), io.SeekStart)
	if err != nil {
		return err
	}

	// setup buffer with max Discord file size
	buf := make([]byte, maxDiscordFileSize)

	var bar *progressbar.ProgressBar

	if !upDebug {
		// init progress bar
		bar = progressbar.DefaultBytes(
			size,
			"Uploading "+localBase,
		)
		err := bar.Add(blockNumber * maxDiscordFileSize)
		if err != nil {
			return err
		}
	}

	first := true

	for {
		blockNumber += 1

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
					Name:   strconv.Itoa(blockNumber),
					Reader: reader,
				},
			},
		}

		// send chunk
		// retry 5 times because internet can be flaky and Discord sometimes
		// likes to drop connections
		var message *discordgo.Message
		maxUploadTries := 10
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

		if first && !upResume {
			// if pin fails, ignore
			// the reason why we pin is because Discord exposes the timestamp
			// of the last pin in channel which is useful for obtaining the
			// file's creation date
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
