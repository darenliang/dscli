package cmd

import (
	"github.com/darenliang/dscli/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var botToken string
var serverID string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure application",
	RunE:  config,
}

func init() {
	configCmd.Flags().StringVarP(&botToken, "token", "t", "", "Discord bot token")
	configCmd.Flags().StringVarP(&serverID, "id", "i", "", "Discord server id to upload files")

	rootCmd.AddCommand(configCmd)
}

// config command handler
func config(cmd *cobra.Command, args []string) error {
	// parse discord bot token
	err := common.SetConfigVal(
		"token",
		cmd.Flag("token").Value.String(),
		`Discord bot token not provided.
The token will be used to run the Discord bot from the CLI app.`,
		"Please enter a Discord bot token: ",
	)
	if err != nil {
		return err
	}

	// parse server id
	err = common.SetConfigVal(
		"id",
		cmd.Flag("id").Value.String(),
		`Server ID not provided.
The server ID will be used to write files to a Discord server.`,
		"Please enter a server ID: ",
	)
	if err != nil {
		return err
	}

	// write config
	if viper.SafeWriteConfig() != nil {
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
	}

	color.Green("Config set.")
	return nil
}
