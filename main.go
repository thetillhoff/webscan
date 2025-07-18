/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/thetillhoff/webscan/v3/pkg/logger"
	"github.com/thetillhoff/webscan/v3/pkg/webscan"
	"github.com/urfave/cli/v3"
)

var version = "dev" // This is just the default. The actual value is injected at compiletime
var verbosity int

func main() {

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "prints just the version of webscan",
	}
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Println(cmd.Root().Version)
	}
	cli.RootCommandHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
 USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	{{if len .Authors}}
 AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
 COMMANDS:
 {{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
 GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{if .Copyright }}
 COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
 VERSION:
	{{.Version}}
	{{end}}
 `

	app := cli.Command{
		Name:                   "webscan",
		Usage:                  "Verifies web things",
		Version:                version,
		HideHelpCommand:        true,
		EnableShellCompletion:  true,
		UseShortOptionHandling: true, // Allow not only `-v -v -v`, but also `-vvv`
		Commands: []*cli.Command{
			{
				Name:            "completion",
				Usage:           "Generate the autocompletion script for the specified shell",
				HideHelpCommand: true, // Disable `help` subcommand
				Commands: []*cli.Command{
					{
						Name:  "bash",
						Usage: "Generate the autocompletion script for bash",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							autocomplete_bash()
							return nil
						},
					},
					{
						Name:  "powershell",
						Usage: "Generate the autocompletion script for powershell",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							autocomplete_powershell()
							return nil
						},
					},
					{
						Name:  "zsh",
						Usage: "Generate the autocompletion script for zsh",
						Action: func(ctx context.Context, cmd *cli.Command) error {
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
				Config: cli.BoolConfig{
					Count: &verbosity,
				},
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
				Name:  "follow",
				Value: false,
				Usage: "follow CNAME and HTTP redirects",
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
				Usage: "enable all HTTP scans",
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var (
				err error

				level slog.Level

				engine webscan.Engine

				writeMutex sync.Mutex
			)

			// Logging

			switch verbosity {
			case 0: // If it's not set at all, the number is -1, not 0
				level = slog.LevelWarn
			case 1:
				level = slog.LevelInfo
			case 3:
				level = slog.LevelDebug
			default:
				log.Fatalln(errors.New("invalid amount of verbosity flags"))
			}

			// set global logger with custom options
			slog.SetDefault(slog.New(
				logger.NewHandler(
					os.Stderr,
					&writeMutex,
					&slog.HandlerOptions{
						Level: level,
					},
					cmd.Bool("no-color"),
				),
			))

			// Args

			if cmd.Args().Len() != 1 {
				log.Fatalln(errors.New("exactly one argument expected"))
			}

			var stdout io.Writer
			if cmd.Bool("quiet") {
				stdout = io.Discard
			} else {
				stdout = os.Stdout
			}

			engine, err = webscan.NewEngine(
				stdout,
				cmd.Bool("no-color"),
				cmd.String("dnsServer"),
				cmd.Bool("follow"),
				cmd.Bool("advancedDnsScan"),
				cmd.Bool("ipScan"),
				cmd.Bool("advancedPortScan"),
				cmd.Bool("tlsScan"),
				cmd.Bool("httpProtocolScan"),
				cmd.Bool("httpHeaderScan"),
				cmd.Bool("htmlContentScan"),
				cmd.Bool("mailConfigScan"),
				cmd.Bool("subDomainScan"),
				&writeMutex,
			)
			if err != nil {
				slog.Error("could not create webscan engine with provided parameters", "error", err)
				os.Exit(1)
			}

			if cmd.Bool("dns") { // Enable advanced dns scans
				engine.EnableDetailedDnsScan()
			}

			if cmd.Bool("ip") { // Enable ip scans
				engine.EnableIpScan()
			}

			if cmd.Bool("port") { // Enable detailed port scans
				engine.EnableDetailedPortScan()
			}

			if cmd.Bool("tls") { // Enable tls scans
				engine.EnableTlsScan()
			}

			if cmd.Bool("protocol") { // Enable http protocol scans
				engine.EnableHttpProtocolScan()
			}

			if cmd.Bool("header") { // Enable http header scans
				engine.EnableHttpHeaderScan()
			}

			if cmd.Bool("content") { // Enable http content scans
				engine.EnableHtmlContentScan()
			}

			if cmd.Bool("mail") { // Enable mail scans
				engine.EnableMailConfigScan()
			}

			if cmd.Bool("subdomains") { // Enable subdomain scans
				engine.EnableSubdomainScan()
			}

			if cmd.Bool("web") { // Enable webscans only
				engine.EnableWebScans()
			}

			engine.EnableAllScansIfNoneAreExplicitlySet()

			err = engine.Scan(cmd.Args().First())
			if err != nil {
				slog.Error("could not scan selected target", "error", err.Error())
				os.Exit(2)
			}

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

	// cmd.Execute()
}
