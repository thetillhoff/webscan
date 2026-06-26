/*
Copyright © 2023-2026 Till Hoffmann <till@thetillhoff.de>

Web frontend for webscan security scanner
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v5/pkg/logger"
	"github.com/thetillhoff/webscan/v5/pkg/webserver"
)

var version = "dev" // This is just the default. The actual value is injected at buildtime

func main() {
	// Parse command-line flags
	port := flag.String("port", "8080", "Port to listen on")
	dnsServer := flag.String("dns", "", "Custom DNS server (default: system DNS)")
	follow := flag.Bool("follow", false, "Follow CNAMEs and HTTP redirects")
	requestTimeout := flag.Duration("timeout", 5*time.Second, "Timeout for individual network requests (DNS, port, HTTP, RDAP)")
	noColor := flag.Bool("no-color", false, "Disable colored output")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	scanTimeout := flag.Duration("scan-timeout", 30*time.Second, "Maximum duration per scan job")
	maxConcurrentScans := flag.Int("max-concurrent-scans", 1, "Maximum number of scan workers")
	maxRequestBytes := flag.Int64("max-request-bytes", 4096, "Maximum request body size in bytes")
	redisAddr := flag.String("redis-addr", "127.0.0.1:6379", "Redis server address")
	redisPassword := flag.String("redis-password", "", "Redis password")
	redisDB := flag.Int("redis-db", 0, "Redis database index")
	maxQueueSize := flag.Int("max-queue-size", 1000, "Maximum number of queued jobs")
	versionFlag := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	// Set up logging
	var writeMutex sync.Mutex

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}

	// Set global logger with custom options
	slog.SetDefault(slog.New(
		logger.NewHandler(
			os.Stderr,
			&writeMutex,
			&slog.HandlerOptions{
				Level: level,
			},
			*noColor,
		),
	))

	// Args - if provided, serve specific target directly
	var target string
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	// Parse blocklists from env vars (one entry per line)
	domainBlocklist := parseBlocklist(os.Getenv("DOMAIN_BLOCKLIST"))
	ipBlocklist := parseBlocklist(os.Getenv("IP_BLOCKLIST"))
	// Private/loopback/link-local targets are blocked by default (SSRF guard).
	// Set ALLOW_PRIVATE_TARGETS=1 to scan internal networks deliberately.
	allowPrivateTargets := os.Getenv("ALLOW_PRIVATE_TARGETS") == "1"

	// Create and run server
	server, err := webserver.NewServer(
		version,
		*noColor,
		*dnsServer,
		*follow,
		*requestTimeout,
		*port,
		&writeMutex,
		*scanTimeout,
		*maxConcurrentScans,
		*maxRequestBytes,
		*redisAddr,
		*redisPassword,
		*redisDB,
		*maxQueueSize,
		domainBlocklist,
		ipBlocklist,
		allowPrivateTargets,
	)
	if err != nil {
		slog.Error("Could not create webscan web server", "error", err)
		os.Exit(1)
	}

	slog.Info("✅ webscan web server ready")
	if target != "" {
		slog.Info("Will scan target automatically", "target", target)
	}

	// Run server
	err = server.Run()
	if err != nil {
		slog.Error("Web server shut down with error", "error", err)
		os.Exit(2)
	}

	log.Println("Web server stopped")
}

func parseBlocklist(raw string) []string {
	var entries []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			entries = append(entries, line)
		}
	}
	return entries
}
