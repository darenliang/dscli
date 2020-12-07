package cmd

import (
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// quickstartCmd represents the quickstart command
var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "A short quickstart guide on how to use dscli",
	RunE: func(cmd *cobra.Command, args []string) error {
		return browser.OpenURL("https://github.com/darenliang/dscli/blob/master/quickstart/README.md")
	},
}

func init() {
	rootCmd.AddCommand(quickstartCmd)
}
