package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/jwalton/gchalk"
	"github.com/mattn/go-isatty"
)

var _ slog.Handler = (*Handler)(nil)

// type Options struct {
// 	// Enable source code location (Default: false)
// 	AddSource bool

// 	// Minimum level to log (Default: slog.LevelInfo)
// 	Level slog.Leveler

// 	// Disable color (Default: false)
// 	NoColor bool
// }

type Handler struct {
	slogHandler slog.Handler
	buf         *bytes.Buffer
	WriteMutex  *sync.Mutex
	opts        *slog.HandlerOptions

	isTTY   bool
	noColor bool

	timeFormatter func(string) string

	ErrorFormatter func(string) string
	WarnFormatter  func(string) string
	InfoFormatter  func(string) string
	DebugFormatter func(string) string
}

func NewHandler(w io.Writer, writeMutex *sync.Mutex, opts *slog.HandlerOptions, noColor bool) slog.Handler {
	var handler = Handler{
		slogHandler:    slog.NewTextHandler(w, opts),
		WriteMutex:     writeMutex,
		opts:           opts,
		timeFormatter:  func(s string) string { return s },
		ErrorFormatter: func(s string) string { return s },
		WarnFormatter:  func(s string) string { return s },
		InfoFormatter:  func(s string) string { return s },
		DebugFormatter: func(s string) string { return s },
		isTTY:          false,
		noColor:        false,
	}

	if !noColor { // If  noColor is not already set
		if value, ok := os.LookupEnv("NO_COLOR"); ok && value != "" { // Check for env var "$NO_COLOR"
			slog.Debug("$NO_COLOR is set")
			noColor = true
		} else if value, ok := os.LookupEnv(fmt.Sprintf("%s_NO_COLOR", strings.ToUpper(os.Args[0]))); ok && value != "" { // Check for env var "$MYAPP_NO_COLOR" (where value for MYAPP is retrieved at runtime - which means it doesn't work with `go run`)
			slog.Debug(fmt.Sprintf("$%s_NO_COLOR is set", strings.ToUpper(os.Args[0])))
			noColor = true
		} else {
			slog.Debug(fmt.Sprintf("Neither $NO_COLOR nor $%s_NO_COLOR are set", strings.ToUpper(os.Args[0])))
		}
	}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		handler.isTTY = true
	} else if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		handler.isTTY = true
	} else {
		handler.isTTY = false
	}

	if value, ok := os.LookupEnv("TERM"); ok && value == "dumb" { // Check for env var "$TERM"
		slog.Debug("$TERM is dumb")
		noColor = true        // Override
		handler.isTTY = false // Override
	}

	if handler.isTTY {
		slog.Debug("handler.isTTY is true")
	} else {
		slog.Debug("handler.isTTY is false")
	}

	if noColor {

		slog.Debug("handler.noColor is true")

		handler.timeFormatter = func(s string) string { return s }

		handler.ErrorFormatter = func(s string) string { return s }
		handler.WarnFormatter = func(s string) string { return s }
		handler.InfoFormatter = func(s string) string { return s }
		handler.DebugFormatter = func(s string) string { return s }

	} else { // Color enabled

		slog.Debug("handler.noColor is false")

		handler.timeFormatter = func(s string) string { return gchalk.Dim(s) }

		handler.ErrorFormatter = func(s string) string { return gchalk.Red(s) }
		handler.WarnFormatter = func(s string) string { return gchalk.Yellow(s) }
		handler.InfoFormatter = func(s string) string { return s }
		handler.DebugFormatter = func(s string) string { return gchalk.Dim(s) }

	}

	return handler
}
