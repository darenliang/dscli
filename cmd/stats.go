package cmd

import (
	"fmt"
	"github.com/darenliang/dscli/common"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"strconv"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Statistics for dscli",
	RunE:  stats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func stats(cmd *cobra.Command, args []string) error {
	session, guild, channels, err := common.GetDiscordSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fileMap, err := common.ParseFileMap(channels)
	if err != nil {
		return err
	}

	fmt.Printf("Number of files stored: %d/500\n", len(fileMap))

	// calculate total size
	var totalSize uint64
	for _, channel := range fileMap {
		filesize, err := strconv.ParseUint(channel.Topic, 10, 64)
		if err == nil {
			totalSize += filesize
		}
	}

	fmt.Printf("Total size stored: %s\n", humanize.Bytes(totalSize))

	maxFileSizeUpload, err := common.GetMaxFileSizeUpload(session, guild)
	if err != nil {
		return err
	}

	fmt.Printf("Upload chunk size: %s\n", humanize.Bytes(uint64(maxFileSizeUpload)))

	return nil
}
