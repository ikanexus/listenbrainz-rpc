/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/ikanexus/listenbrainz-rpc/internal"
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
		listenbrainzUser := viper.GetString("user")
		if verbose, err := cmd.Flags().GetBool("verbose"); err == nil && verbose == true {
			log.SetLevel(log.DebugLevel)
		}
		scrobbler := internal.NewScrobbler(listenbrainzUser, discordAppId)
		return scrobbler.Scrobble()
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

	configDefault := fmt.Sprintf("%s/listenbrainz-rpc.yaml", getXdgHome())
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", configDefault, "config file")

	rootCmd.Flags().StringVar(&discordAppId, "app-id", "1231614541905920113", "Discord App ID")
	rootCmd.Flags().String("user", "", "Listenbrainz Username")
	rootCmd.Flags().BoolP("verbose", "v", false, "Show verbose logging")
}

func getXdgHome() string {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome == "" {
		xdgHome = filepath.Join(home, ".config")
	}
	return xdgHome
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		xdgHome := getXdgHome()

		viper.AddConfigPath(xdgHome)
		viper.SetConfigType("yaml")
		viper.SetConfigName("listenbrainz-rpc")
	}

	viper.SetEnvPrefix("listenbrainz")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
