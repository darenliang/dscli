package cmd

import (
	"bufio"
	"fmt"
	"github.com/darenliang/dscli/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	token          string
	serverID       string
	deleteChannels bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure application",
	RunE:  config,
}

func init() {
	configCmd.Flags().StringVarP(&token, "token", "t", "", "Discord token")
	configCmd.Flags().StringVarP(&serverID, "id", "i", "", "Discord server id to upload files")
	configCmd.Flags().BoolVarP(&deleteChannels, "delete", "d", false, "delete channels from server")

	rootCmd.AddCommand(configCmd)
}

// config command handler
func config(cmd *cobra.Command, args []string) error {
	// parse Discord token
	err := common.SetConfigVal(
		"token",
		cmd.Flag("token").Value.String(),
		`Discord token not provided.
The token will be used to run your account from the CLI app.`,
		"Please enter a Discord token: ",
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

	// get input if any of the other flags aren't set
	if cmd.Flag("token").Value.String() == "" || cmd.Flag("id").Value.String() == "" {
		deleteChannels = false
		fmt.Println("The server must have no existing channels when you first use dscli.")
		fmt.Println()
		color.New(color.FgYellow, color.Bold).Print("Delete all channels in server? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice := strings.ToLower(strings.TrimSpace(input))
		fmt.Println()
		if choice == "y" || choice == "yes" {
			deleteChannels = true
		}
	}

	if deleteChannels {
		session, _, channels, err := common.GetDiscordSession()
		if err != nil {
			return err
		}

		for _, channel := range channels {
			_, err := session.ChannelDelete(channel.ID)
			if err != nil {
				return err
			}
		}

		color.Green("All channels deleted")
	} else {
		color.Yellow(`Channels are not deleted.
Note that you must delete all channels before you can start using dscli.
Not doing so will result in a server that has an invalid state.`)
	}
	fmt.Println()

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
