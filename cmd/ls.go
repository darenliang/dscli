package cmd

import (
	"fmt"
	"github.com/darenliang/dscli/common"
	"github.com/spf13/cobra"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"strconv"
	"time"
)

var listView bool

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:        "ls",
	SuggestFor: []string{"list"},
	Short:      "List files",
	RunE:       ls,
}

func init() {
	lsCmd.Flags().BoolVarP(&listView, "list", "l", false, "show list of files with filesize and date created timestamp in the format \"<filesize> <unix timestamp> <filename>\"")

	rootCmd.AddCommand(lsCmd)
}

// ls command handler
func ls(cmd *cobra.Command, args []string) error {
	session, _, channels, err := common.GetDiscordSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fileMap, err := common.ParseFileMap(channels)
	if err != nil {
		return err
	}

	var files []string
	for filename := range fileMap {
		files = append(files, filename)
	}

	collate.New(language.English).SortStrings(files)

	if !listView {
		common.PrintFiles(files)
	} else {
		for _, filename := range files {
			channel := fileMap[filename]
			// ignore error to prevent dscli from locking up
			// obtain file creation date from last pin timestamp in a channel
			timestamp, err := channel.LastPinTimestamp.Parse()
			if err != nil {
				timestamp = time.Unix(0, 0)
			}
			filesize, err := strconv.ParseUint(channel.Topic, 10, 64)
			if err != nil {
				filesize = 0
			}
			fmt.Printf("%d %d %s\n", filesize, timestamp.Unix(), filename)
		}
	}

	return nil
}
