package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/ikanexus/listenbrainz-rpc/internal"
	altsrc "github.com/urfave/cli-altsrc/v3"
	"github.com/urfave/cli-altsrc/v3/yaml"
	"github.com/urfave/cli/v3"
)

func Execute() {
	var configFile string

	cmd := &cli.Command{
		Name:  "listenbrainz-rpc",
		Usage: "A CLI tool to show what you're watching in ListenBrainz as Discord Activity",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "config file",
				Value:       filepath.Join(xdg.ConfigHome, "listenbrainz-rpc.yaml"),
				Sources:     cli.EnvVars("LISTENBRAINZ_CONFIG"),
				Destination: &configFile,
			},
			&cli.StringFlag{
				Name:    "app-id",
				Aliases: []string{"a"},
				Usage:   "Discord App ID",
				Value:   "1232457767726485545",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("LISTENBRAINZ_APP_ID"),
					yaml.YAML("app-id", altsrc.NewStringPtrSourcer(&configFile)),
				),
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "Listenbrainz Username",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("LISTENBRAINZ_USER"),
					yaml.YAML("user", altsrc.NewStringPtrSourcer(&configFile)),
				),
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show verbose logging",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("LISTENBRAINZ_VERBOSE"),
					yaml.YAML("verbose", altsrc.NewStringPtrSourcer(&configFile)),
				),
				Action: func(ctx context.Context, cmd *cli.Command, v bool) error {
					if v {
						log.SetLevel(log.DebugLevel)
					}
					return nil
				},
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			_, _ = tea.LogToFileWith("listenbrainz.log", "", log.Default())
			return ctx, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := internal.Config{
				AppID:   cmd.String("app-id"),
				User:    cmd.String("user"),
				Verbose: cmd.Bool("verbose"),
			}

			m := internal.NewModel(cfg)
			if _, err := tea.NewProgram(m).Run(); err != nil {
				return err
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
