package cmd

import (
	"github.com/darenliang/dscli/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// cfgFile is where the config file is located
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dscli",
	Short: "Dscli is a CLI application to store files on Discord",
	Long: `Dscli (Discord store command-line interface) is a CLI application to store
files on Discord.

Dscli probably does something that Discord doesn't like (aka storing huge
amounts of files). Proceed with caution.`,
	Version: common.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dscli.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".dscli")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		color.Red(err.Error())
	}
}
