package cmd

import (
	"fmt"
	"github.com/darenliang/dscli/common"
	"github.com/spf13/cobra"
	"sort"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:        "ls",
	SuggestFor: []string{"list"},
	Short:      "List files",
	RunE:       ls,
}

func init() {
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

	sort.Strings(files)

	for _, str := range files {
		fmt.Println(str)
	}

	return nil
}
