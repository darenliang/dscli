package cmd

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/dscli/common"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
)

var dlDebug bool

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
	dlCmd.Flags().BoolVarP(&dlDebug, "debug", "d", false, "debug mode: <total bytes> <bytes downloaded>")

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

	// get remote file size
	remoteSize, err := strconv.ParseUint(channel.Topic, 10, 64)
	if err != nil {
		return errors.New(remote + " has a corrupted file size: " + channel.Topic)
	}

	// create local file
	localFile, err := os.Create(local)
	if err != nil {
		return err
	}
	defer localFile.Close()

	var msgs []*discordgo.Message

	var bar *progressbar.ProgressBar
	if !dlDebug {
		bar = progressbar.DefaultBytes(
			int64(remoteSize),
			"Downloading "+remote,
		)
	}

	// used for progress bar and debug modes
	filesize := 0
	progress := 0

	currMsgID := "0"

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

			// get file response
			resp, err := common.HttpClient.Get(msgs[i].Attachments[0].URL)
			if err != nil {
				return err
			}

			// write chunk to local file
			if !dlDebug {
				_, err = io.Copy(io.MultiWriter(localFile, bar), resp.Body)
				if err != nil {
					return err
				}
			} else {
				written, err := io.Copy(localFile, resp.Body)
				if err != nil {
					return err
				}
				progress += int(written)
				fmt.Printf("%d %d \n", filesize, progress)
			}

			// close response
			resp.Body.Close()
		}
	}

	return nil
}
