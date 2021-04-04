package cmd

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/dscli/common"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
)

// dlCmd represents the dl command
var dlCmd = &cobra.Command{
	Use:        "dl <remote file> <local file>",
	Example:    "dl test.txt test.txt",
	SuggestFor: []string{"download"},
	Short:      "Download file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one argument")
		}
		return nil
	},
	RunE: dl,
}

func init() {
	rootCmd.AddCommand(dlCmd)
}

// dl command handler
func dl(cmd *cobra.Command, args []string) error {
	remote := args[0] // remote filename

	var local string // local filename
	if len(args) == 1 {
		local = remote
	} else {
		local = args[1]
	}

	session, _, channels, err := common.GetDiscordSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fileMap, err := common.ParseFileMap(channels)
	if err != nil {
		return err
	}

	// check if remote filename exists
	channel, ok := fileMap[remote]
	if !ok {
		return errors.New(remote + " doesn't exist")
	}

	// create local file
	localFile, err := os.Create(local)
	if err != nil {
		return err
	}
	defer localFile.Close()

	var msgs []*discordgo.Message
	var bar *progressbar.ProgressBar
	currMsgID := "0"
	first := true

	// do-while
	for ok := true; ok; ok = len(msgs) == common.MaxDiscordMessageRequest {
		// query messages after currMsgID
		msgs, err = session.ChannelMessages(
			channel.ID,
			common.MaxDiscordMessageRequest,
			"", currMsgID, "",
		)
		if err != nil {
			return err
		}

		// reverse iterate through msgs
		for i := len(msgs) - 1; i >= 0; i-- {
			currMsgID = msgs[i].ID

			if len(msgs[i].Attachments) < 1 {
				continue
			}

			// on first time, init progressbar with file size
			if first {
				sizeStr := msgs[i].Attachments[0].Filename

				size, err := strconv.Atoi(sizeStr)
				if err != nil {
					return err
				}

				bar = progressbar.DefaultBytes(
					int64(size),
					"Downloading "+remote,
				)

				first = false
			}

			// get file response
			resp, err := common.HttpClient.Get(msgs[i].Attachments[0].URL)
			if err != nil {
				return err
			}

			// write chunk to local file
			_, err = io.Copy(io.MultiWriter(localFile, bar), resp.Body)
			if err != nil {
				return err
			}

			// close response
			err = resp.Body.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
