/*
Copyright © 2023 Till Hoffmann <till@thetillhoff.de>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thetillhoff/webscan/pkg/webscan"
)

var cfgFile string
var dnsServer string
var dkimSelector string
var all bool
var follow bool
var verbose bool

var detailedDnsScan bool
var ipScan bool
var detailedPortScan bool
var tlsScan bool
var webScans bool
var httpProtocolScan bool
var httpHeaderScan bool
var httpContentScan bool
var mailConfigScan bool
var subdomainScan bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webscan",
	Short: "Verifies web things",
	Args:  cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			url = args[0] // Only one arg is allowed

			engine webscan.Engine
		)

		if dnsServer != "" {
			engine = webscan.EngineWithCustomDns(url, dnsServer)
		} else {
			engine = webscan.DefaultEngine(url)
		}
		engine.DkimSelector = dkimSelector
		engine.FollowRedirects = follow // follow is always optional
		engine.Verbose = verbose        // verbose is always optional

		if all {
			engine.DetailedDnsScan = true
			engine.IpScan = true
			engine.DetailedPortScan = true
			engine.TlsScan = true
			engine.HttpProtocolScan = true
			engine.HttpHeaderScan = true
			engine.HttpContentScan = true
			// engine.MailConfigScan = true // TODO
			// engine.SubdomainScan = true // TODO
		} else {
			engine.DetailedDnsScan = detailedDnsScan
			engine.IpScan = ipScan
			engine.DetailedPortScan = detailedPortScan
			engine.TlsScan = tlsScan
			engine.HttpProtocolScan = httpProtocolScan || webScans
			engine.HttpHeaderScan = httpHeaderScan || webScans
			engine.HttpContentScan = httpContentScan || webScans
			engine.MailConfigScan = mailConfigScan
			engine.SubdomainScan = subdomainScan
		}

		engine, err = engine.Scan(url)
		if err != nil {
			log.Fatalln(err)
		}

		engine.PrintScanResults()

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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "increase verbosity")

	rootCmd.PersistentFlags().BoolVarP(&detailedDnsScan, "dns", "d", false, "enable detailed DNS scan")
	rootCmd.PersistentFlags().BoolVarP(&ipScan, "ip", "i", false, "enable IP analysis")
	rootCmd.PersistentFlags().BoolVar(&detailedPortScan, "port", false, "enable detailed port scanning")
	rootCmd.PersistentFlags().BoolVarP(&tlsScan, "tls", "t", false, "enable TLS scan")
	rootCmd.PersistentFlags().BoolVar(&httpProtocolScan, "protocol", false, "enable HTTP protocol scan")
	rootCmd.PersistentFlags().BoolVar(&httpHeaderScan, "header", false, "enable HTTP header scan")
	rootCmd.PersistentFlags().BoolVarP(&webScans, "web", "w", false, "enable all HTTP scans except content scans")
	rootCmd.PersistentFlags().BoolVarP(&httpContentScan, "content", "c", false, "enable HTTP content scan")
	rootCmd.PersistentFlags().BoolVarP(&mailConfigScan, "mail", "m", false, "enable mail config scan")
	rootCmd.PersistentFlags().BoolVarP(&subdomainScan, "subdomains", "s", false, "enable subdomains search")

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
