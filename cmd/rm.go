package cmd

import (
	"errors"
	"github.com/darenliang/dscli/common"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:        "rm <remote file>",
	Example:    "rm example.txt",
	SuggestFor: []string{"remove", "delete"},
	Short:      "Remove file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires one argument")
		}
		return nil
	},
	RunE: rm,
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

// rm command handler
func rm(cmd *cobra.Command, args []string) error {
	session, _, channels, err := common.GetDiscordSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fileMap, err := common.ParseFileMap(channels)
	if err != nil {
		return err
	}

	filename := args[0]

	// old file exists
	if channel, ok := fileMap[filename]; ok {
		// remove file (via channel delete)
		_, err = session.ChannelDelete(channel.ID)
		if err != nil {
			return errors.New("cannot delete file: " + err.Error())
		}
		return nil
	}

	return errors.New(filename + " not found")
}
