package cmd

import (
	"fmt"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var header = color.New(color.FgYellow, color.Bold, color.Underline)

// quickstartCmd represents the quickstart command
var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "A short quickstart guide on how to use dscli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Welcome to the dscli quickstart guide!

Author: Daren Liang <daren@darenliang.com>
Repo: https://github.com/darenliang/dscli`)
		fmt.Println()
		header.Println("Step 1: Create a Discord bot and get a token")
		fmt.Println(` - Visit this page: https://discord.com/developers/applications
 - Click "New Application"
 - Enter anything as the name
 - Click "Create"
 - Go to "Bot"
 - Click "Add Bot"
 - Click "Yes, do it!"
 - Copy the token by clicking on "Copy"
 - Save this token somewhere safe`)
		fmt.Println()
		header.Println("Step 2: Enable developer mode")
		fmt.Println(` - Go to the Discord client or https://discord.com/app
 - Click on the gear icon to open settings
 - Click on "Appearance"
 - Toggle "Developer Mode" to on`)
		fmt.Println()
		header.Println("Step 3: Create a Discord server and get a server id")
		fmt.Println(` - Go to the Discord client or https://discord.com/app
 - Click on the plus icon on the left panel
 - Click on "Create My Own"
 - Right click server icon on the left panel
 - Click "Copy ID"
 - Save this id somewhere
 - Optional: Remove the general text channel in the server you created`)
		fmt.Println()
		header.Println("Step 4: Configure dscli")
		fmt.Println(` - Run "dscli configure" from command-line ("dscli.exe configure" on Windows)
 - Provide Discord bot token you saved
 - Provide server id you saved`)
		fmt.Println()
		color.New(color.FgGreen, color.Bold).Println("You are ready to use dscli!")
	},
}

func init() {
	rootCmd.AddCommand(quickstartCmd)
}
