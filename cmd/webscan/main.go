/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/logger"
	"github.com/thetillhoff/webscan/v3/pkg/webscan"
	"github.com/urfave/cli/v3"
)

var version = "dev" // This is just the default. The actual value is injected at compiletime
var verbosity int

func main() {

	// Version flag: only long form (--version) to avoid conflict with -v (verbose)
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "prints just the version of webscan",
		// No Aliases field = only accepts --version, not -v (which is used for verbose)
	}
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Println(cmd.Root().Version)
	}
	cli.RootCommandHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
USAGE:
	{{if .VisibleFlags}}[global options]{{end}}[target|completion [subcommand]]
{{if .Commands}}
COMMANDS:{{range .Commands}}{{if not .HideHelp}}
	{{join .Names ", "}}{{ "\t"}}{{.Usage}}{{end}}{{end}}{{end}}
{{if .VisibleFlags}}
GLOBAL OPTIONS:{{range .VisibleFlags}}
	{{.}}{{end}}{{end}}
`

	app := cli.Command{
		Name:                   "webscan",
		Usage:                  "Verifies web things",
		Version:                version,
		HideHelpCommand:        true,
		EnableShellCompletion:  true,
		UseShortOptionHandling: true, // Allow not only `-v -v -v`, but also `-vvv`
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "ns",
				Value: "",
				Usage: "set custom dns server (uses system dns by default)",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Value: 5 * time.Second,
				Usage: "timeout for individual network requests (DNS, port, HTTP, RDAP)",
			},
			&cli.StringFlag{
				Name:  "dkim-selector",
				Value: "",
				Usage: "set dkim-selector as in <dkim-selector>._domainkey.<your-domain>.<tld>",
			},
			&cli.BoolFlag{
				Name:  "follow",
				Value: false,
				Usage: "follow CNAME records and HTTP redirects",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "increase verbosity to Debug level (`-v` for info, `-vvv` for debug)",
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
				Name:  "dns",
				Value: false,
				Usage: "focus on detailed DNS scan",
			},
			&cli.BoolFlag{
				Name:  "ip",
				Value: false,
				Usage: "focus on IP analysis",
			},
			&cli.BoolFlag{
				Name:  "port",
				Value: false,
				Usage: "focus on detailed port scanning",
			},
			&cli.BoolFlag{
				Name:  "tls",
				Value: false,
				Usage: "focus on TLS scan",
			},
			&cli.BoolFlag{
				Name:  "protocol",
				Value: false,
				Usage: "focus on HTTP protocol scan",
			},
			&cli.BoolFlag{
				Name:  "header",
				Value: false,
				Usage: "focus on HTTP header scan",
			},
			&cli.BoolFlag{
				Name:  "content",
				Value: false,
				Usage: "focus on HTTP content scan",
			},
			&cli.BoolFlag{
				Name:  "files",
				Value: false,
				Usage: "focus on well-known files scan (robots.txt, security.txt, sensitive files)",
			},
			&cli.BoolFlag{
				Name:  "web",
				Value: false,
				Usage: "focus on all HTTP scans",
			},
			&cli.BoolFlag{
				Name:  "mail",
				Value: false,
				Usage: "focus on mail config scan",
			},
			&cli.BoolFlag{
				Name:  "subdomains",
				Value: false,
				Usage: "focus on subdomains search",
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

			switch {
			case verbosity <= 0:
				level = slog.LevelWarn
			case verbosity == 1:
				level = slog.LevelInfo
			default:
				level = slog.LevelDebug
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

			stdout := io.Writer(os.Stdout)
			statusOut := io.Writer(os.Stderr)
			if cmd.Bool("quiet") {
				statusOut = io.Discard
			}

			engine, err = webscan.NewEngine(
				stdout,
				statusOut,
				cmd.Bool("no-color"),
				cmd.String("ns"),
				cmd.Bool("follow"),
				cmd.Duration("timeout"),
				cmd.Bool("dns"),
				cmd.Bool("ip"),
				cmd.Bool("port"),
				cmd.Bool("tls"),
				cmd.Bool("protocol") || cmd.Bool("web"),
				cmd.Bool("header") || cmd.Bool("web"),
				cmd.Bool("content") || cmd.Bool("web"),
				cmd.Bool("mail"),
				cmd.Bool("subdomains"),
				&writeMutex,
			)
			if err != nil {
				slog.Error("could not create webscan engine with provided parameters", "error", err)
				os.Exit(1)
			}

			webscan.NewScanOptions(
				cmd.Bool("dns"),
				cmd.Bool("ip"),
				cmd.Bool("port"),
				cmd.Bool("tls"),
				cmd.Bool("protocol"),
				cmd.Bool("header"),
				cmd.Bool("content"),
				cmd.Bool("files"),
				cmd.Bool("mail"),
				cmd.Bool("subdomains"),
			).Apply(&engine)

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
