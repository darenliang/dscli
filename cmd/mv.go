package cmd

import (
	"errors"
	"github.com/darenliang/dscli/common"
	"github.com/spf13/cobra"
)

// mvCmd represents the mv command
var mvCmd = &cobra.Command{
	Use:        "mv source dest",
	Example:    "mv old_file.txt new_file.txt",
	SuggestFor: []string{"move", "rename"},
	Short:      "Move file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires two arguments")
		}
		return nil
	},
	RunE: mv,
}

func init() {
	rootCmd.AddCommand(mvCmd)
}

// mv command handler
func mv(cmd *cobra.Command, args []string) error {
	oldFilename := args[0] // old filename
	newFilename := args[1] // new filename

	// encode new filename
	encodedNewFilename, err := common.EncodeFilename(newFilename)
	if err != nil {
		return err
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

	// old file exists
	if channel, ok := fileMap[oldFilename]; ok {
		// filename with new filename exists
		if _, exists := fileMap[newFilename]; exists {
			return errors.New(newFilename + " already exists on Discord")
		}
		// rename file (via channel rename)
		_, err = session.ChannelEdit(channel.ID, encodedNewFilename)
		if err != nil {
			return errors.New("cannot move file: " + err.Error())
		}
		return nil
	}

	return errors.New(oldFilename + " not found")
}
