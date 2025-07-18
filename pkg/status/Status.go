package status

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/mattn/go-isatty"
)

type Status struct {
	isTTY bool

	writeMutex *sync.Mutex
	out        io.Writer

	// Sets the speed of the spinner.
	// Default is 100ms.
	SpinnerUpdateInterval time.Duration
	spinnerChars          []string
	spinnerIndex          int
	spinnerCompleteChar   string
	spinnerMessage        string
	spinnerStop           chan struct{}
	spinning              bool
	spinnerCharFormatter  func(s string) string

	spinningXOfTotal   int
	spinningXOfCurrent int
	spinningXOfMessage string
}

func NewStatus(noColor bool, writeMutex *sync.Mutex, out io.Writer) Status { // TODO use opts object (pointer) for quiet and noColor
	var (
		status = Status{
			writeMutex: writeMutex,
			out:        out,

			// Sets the speed of the spinner.
			// Default is 100ms.
			SpinnerUpdateInterval: 100 * time.Millisecond,

			spinnerChars:         []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"},
			spinnerCharFormatter: func(s string) string { return s },
			spinnerIndex:         0,
			spinning:             false,
		}
	)

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
		status.isTTY = true
	} else if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		status.isTTY = true
	} else {
		status.isTTY = false
	}

	if value, ok := os.LookupEnv("TERM"); ok && value == "dumb" { // Check for env var "$TERM"
		slog.Debug("$TERM is dumb")
		noColor = true       // Override
		status.isTTY = false // Override
	}

	if status.isTTY {
		slog.Debug("status.isTTY is true")
	} else {
		slog.Debug("status.isTTY is false")
	}

	if noColor {

		slog.Debug("status.noColor is true")

		status.spinnerCompleteChar = "✓"

	} else {

		slog.Debug("status.noColor is false")

		status.spinnerCharFormatter = func(s string) string {
			return gchalk.Dim(s)
		}

		status.spinnerCompleteChar = gchalk.Green("✓")

	}

	return status
}
