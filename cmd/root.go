/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ikanexus/listenbrainz-rpc/internal/discord"
	"github.com/ikanexus/listenbrainz-rpc/internal/scrobble"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var discordAppId string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "listenbrainz-rpc",
	Short: "A CLI tool to show what you're watching in ListenBrainz as Discord Activity",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		discord.Login(discordAppId)
		defer discord.Logout()
		listenbrainzUser, err := cmd.Flags().GetString("user")
		cobra.CheckErr(err)
		return scrobble.Scrobble(listenbrainzUser)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/listenbrainz-rpc.yaml)")

	rootCmd.Flags().StringVar(&discordAppId, "app-id", "1231614541905920113", "Discord App ID")
	rootCmd.Flags().String("user", "", "Listenbrainz Username")
}

func getXdgHome(homeDir string) string {
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome == "" {
		xdgHome = filepath.Join(homeDir, ".config")
	}
	return xdgHome
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		xdgHome := getXdgHome(home)

		// Search config in home directory with name ".listenbrainz-rpc" (without extension).
		viper.AddConfigPath(xdgHome)
		viper.SetConfigType("yaml")
		viper.SetConfigName("listenbrainz-rpc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
