/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/thetillhoff/webscan/pkg/logger"
	"github.com/thetillhoff/webscan/pkg/webscan"
	"github.com/urfave/cli/v2"
)

var version = "dev" // This is just the default. The actual value is injected at compiletime
var verbosity int

func main() {

	app := &cli.App{
		Name:                 "webscan",
		Usage:                "Verifies web things",
		Version:              version,
		HideVersion:          true, // Disable `-v`, `--version` flags
		HideHelpCommand:      true, // Disable `help` subcommand
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Prints the version of webscan",
				Action: func(cCtx *cli.Context) error {
					fmt.Printf("%s\n", cCtx.App.Version)
					return nil
				},
			},
			{
				Name:            "completion",
				Usage:           "Generate the autocompletion script for the specified shell",
				HideHelpCommand: true, // Disable `help` subcommand
				Subcommands: []*cli.Command{
					{
						Name:  "bash",
						Usage: "Generate the autocompletion script for bash",
						Action: func(ctx *cli.Context) error {
							autocomplete_bash()
							return nil
						},
					},
					{
						Name:  "powershell",
						Usage: "Generate the autocompletion script for powershell",
						Action: func(ctx *cli.Context) error {
							autocomplete_powershell()
							return nil
						},
					},
					{
						Name:  "zsh",
						Usage: "Generate the autocompletion script for zsh",
						Action: func(ctx *cli.Context) error {
							autocomplete_zsh()
							return nil
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "ns",
				Value: "",
				Usage: "set custom dns server (default uses system dns)",
			},
			&cli.StringFlag{
				Name:  "dkim-selector",
				Value: "",
				Usage: "set dkim-selector as in <dkim-selector>._domainkey.<your-domain>.<tld>",
			},
			&cli.BoolFlag{
				Name:  "follow",
				Value: false,
				Usage: "follow redirects and CNAME checks",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "increase verbosity to Debug level",
				Count:   &verbosity,
			},
			&cli.BoolFlag{
				Name:  "instant",
				Value: false,
				Usage: "print results immediately after each scan instead of after all scans completed",
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Value: false,
				Usage: "disable coloring of results and logs",
			},
			&cli.BoolFlag{
				Name:  "quiet",
				Value: false,
				Usage: "disable status updates and only prints results",
			},
			&cli.BoolFlag{
				Name:  "all",
				Value: true,
				Usage: "enable all checks",
			},
			&cli.BoolFlag{
				Name:  "dns",
				Value: false,
				Usage: "enable detailed DNS scan",
			},
			&cli.BoolFlag{
				Name:  "ip",
				Value: false,
				Usage: "enable IP analysis",
			},
			&cli.BoolFlag{
				Name:  "port",
				Value: false,
				Usage: "enable detailed port scanning",
			},
			&cli.BoolFlag{
				Name:  "tls",
				Value: false,
				Usage: "enable TLS scan",
			},
			&cli.BoolFlag{
				Name:  "protocol",
				Value: false,
				Usage: "enable HTTP protocol scan",
			},
			&cli.BoolFlag{
				Name:  "header",
				Value: false,
				Usage: "enable HTTP header scan",
			},
			&cli.BoolFlag{
				Name:  "content",
				Value: false,
				Usage: "enable HTTP content scan",
			},
			&cli.BoolFlag{
				Name:  "web",
				Value: false,
				Usage: "enable all HTTP scans except content scan",
			},
			&cli.BoolFlag{
				Name:  "mail",
				Value: false,
				Usage: "enable mail config scan",
			},
			&cli.BoolFlag{
				Name:  "subdomains",
				Value: false,
				Usage: "enable subdomains search",
			},
		},
		Action: func(cCtx *cli.Context) error {
			var (
				err error

				level slog.Level

				engine webscan.Engine2

				writeMutex sync.Mutex
			)

			if cCtx.Args().Len() != 1 {
				log.Fatalln(errors.New("exactly one argument expected"))
			}

			// Logging

			switch verbosity {
			case -1: // If it's not set at all, the number is -1, not 0
				level = slog.LevelWarn
			case 1:
				level = slog.LevelInfo
			case 3:
				level = slog.LevelDebug
			default:
				log.Fatalln(errors.New("invalid verbosity amount"))
			}

			w := os.Stderr

			// set global logger with custom options
			slog.SetDefault(slog.New(
				logger.NewHandler(
					w,
					&writeMutex,
					&slog.HandlerOptions{
						Level: level,
					},
					cCtx.Bool("no-color"),
				),
			))

			engine, err = webscan.NewEngine(
				cCtx.Bool("quiet"),
				cCtx.Bool("noColor"),
				cCtx.String("dnsServer"),
				cCtx.Bool("follow"),
				cCtx.Bool("instant"),
				cCtx.Bool("advancedDnsScan"),
				cCtx.Bool("ipScan"),
				cCtx.Bool("advancedPortScan"),
				cCtx.Bool("tlsScan"),
				cCtx.Bool("httpProtocolScan"),
				cCtx.Bool("httpHeaderScan"),
				cCtx.Bool("htmlContentScan"),
				cCtx.Bool("mailConfigScan"),
				cCtx.Bool("subDomainScan"),
				&writeMutex,
			)
			if err != nil {
				slog.Error("could not create webscan engine with provided parameters", "error", err)
				os.Exit(1)
			}

			if cCtx.Bool("web") { // Enable webscans only
				engine.EnableWebScans()
			}

			if !(cCtx.Bool("advancedDnsScan") ||
				cCtx.Bool("ipScan") ||
				cCtx.Bool("advancedPortScan") ||
				cCtx.Bool("tlsScan") ||
				cCtx.Bool("httpProtocolScan") ||
				cCtx.Bool("httpHeaderScan") ||
				cCtx.Bool("htmlContentScan") ||
				cCtx.Bool("mailConfigScan") ||
				cCtx.Bool("subDomainScan") ||
				cCtx.Bool("web")) {

				engine.EnableAllScans()
			}

			err = engine.Scan(cCtx.Args().First())
			if err != nil {
				slog.Error("could not scan selected target", "error", err.Error())
				os.Exit(2)
			}

			engine.PrintResults()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// cmd.Execute()
}
