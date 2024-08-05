/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thetillhoff/webscan/pkg/logger"
	"github.com/thetillhoff/webscan/pkg/webscan"
)

var cfgFile string
var dnsServer string
var dkimSelector string
var all bool
var follow bool
var verboseInfo bool
var verboseDebug bool

var instant bool
var noColor bool
var quiet bool

var advancedDnsScan bool
var ipScan bool
var advancedPortScan bool
var tlsScan bool
var webScans bool
var httpProtocolScan bool
var httpHeaderScan bool
var htmlContentScan bool
var mailConfigScan bool
var subDomainScan bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webscan",
	Short: "Verifies web things",
	Args:  cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err    error
			target = args[0] // Only one arg is allowed

			level slog.Level

			engine webscan.Engine2

			writeMutex sync.Mutex
		)

		if verboseDebug { // Most verbose
			level = slog.LevelDebug
		} else if verboseInfo { // Verbose
			level = slog.LevelInfo
		} else { // Default
			level = slog.LevelWarn // TODO Warn vs Error... for example if a request fails and the app continues anyway...
		}

		// Logging

		w := os.Stderr

		// set global logger with custom options
		slog.SetDefault(slog.New(
			logger.NewHandler(
				w,
				&writeMutex,
				&slog.HandlerOptions{
					Level: level,
				},
				noColor,
			),
		))

		engine, err = webscan.NewEngine(
			quiet,
			noColor,
			dnsServer,
			target,
			follow,
			instant,
			advancedDnsScan,
			ipScan,
			advancedPortScan,
			tlsScan,
			httpProtocolScan,
			httpHeaderScan,
			htmlContentScan,
			mailConfigScan,
			subDomainScan,
			&writeMutex,
		)
		if err != nil {
			slog.Error("could not create webscan engine with provided parameters", "error", err)
			os.Exit(1)
		}

		// TODO
		// engine.DkimSelector = dkimSelector
		// engine.FollowRedirects = follow // follow is always optional
		// engine.Verbose = verbose        // verbose is always optional

		if webScans { // Enable webscans only
			engine.EnableWebScans()
		}

		if !(advancedDnsScan ||
			ipScan ||
			advancedPortScan ||
			tlsScan ||
			httpProtocolScan ||
			httpHeaderScan ||
			htmlContentScan ||
			mailConfigScan ||
			subDomainScan) {
			engine.EnableAllScans() // Enable all scans by default
		}

		err = engine.Scan(target)
		if err != nil {
			slog.Error("could not scan selected target", "error", err.Error())
			os.Exit(2)
		}

		engine.PrintResults()

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.webscan.yaml)")

	rootCmd.PersistentFlags().StringVar(&dnsServer, "ns", "", "custom dns server (default uses system dns)")
	rootCmd.PersistentFlags().StringVar(&dkimSelector, "dkim-selector", "", "set dkim-selector as in <dkim-selector>._domainkey.<your-domain>.<tld>")
	rootCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "enables all checks")
	rootCmd.PersistentFlags().BoolVarP(&follow, "follow", "f", false, "enables following redirects and CNAME checks")
	// rootCmd.PersistentFlags().BoolVarP(&verboseInfo, "verbose", "v", false, "increase verbosity to Info level")
	// TODO due to neglecting of cobra, viper, pflags, a change to https://cli.urfave.org/v2/examples/flags/ is recommended
	// By doing that, the goal would be that verboseDebug can be enabled with `-vvv` and verboseInfo with `-v`
	rootCmd.PersistentFlags().BoolVarP(&verboseDebug, "debug", "v", false, "increase verbosity to Debug level")

	rootCmd.PersistentFlags().BoolVar(&instant, "instant", false, "print results immediately after each scan instead of after all scans completed")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disables coloring of results and logs")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "disables status updates and only prints results")

	rootCmd.PersistentFlags().BoolVarP(&advancedDnsScan, "dns", "d", false, "enable detailed DNS scan")
	rootCmd.PersistentFlags().BoolVarP(&ipScan, "ip", "i", false, "enable IP analysis")
	rootCmd.PersistentFlags().BoolVar(&advancedPortScan, "port", false, "enable detailed port scanning")
	rootCmd.PersistentFlags().BoolVarP(&tlsScan, "tls", "t", false, "enable TLS scan")
	rootCmd.PersistentFlags().BoolVar(&httpProtocolScan, "protocol", false, "enable HTTP protocol scan")
	rootCmd.PersistentFlags().BoolVar(&httpHeaderScan, "header", false, "enable HTTP header scan")
	rootCmd.PersistentFlags().BoolVarP(&webScans, "web", "w", false, "enable all HTTP scans except content scans")
	rootCmd.PersistentFlags().BoolVarP(&htmlContentScan, "content", "c", false, "enable HTTP content scan")
	rootCmd.PersistentFlags().BoolVarP(&mailConfigScan, "mail", "m", false, "enable mail config scan")
	rootCmd.PersistentFlags().BoolVarP(&subDomainScan, "subdomains", "s", false, "enable subdomains search")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".webscan" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".webscan")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
