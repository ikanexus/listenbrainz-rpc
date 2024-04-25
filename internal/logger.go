package internal

import (
	"bufio"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	charm_log "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger = NewLogger()

func NewLogger() *charm_log.Logger {

	file, err := os.Open("listenbrainz.log")
	cobra.CheckErr(err)

	level := charm_log.InfoLevel
	if viper.GetBool("verbose") {
		level = charm_log.DebugLevel
	}
	logger := charm_log.NewWithOptions(bufio.NewWriter(file), charm_log.Options{
		Formatter: charm_log.JSONFormatter,
		Level:     level,
	})
	f, err := tea.LogToFile("listenbrainz1.log", "debug")
	if err != nil {
		cobra.CheckErr(err)
	}
	defer f.Close()
	defer file.Close()
	return logger
}
