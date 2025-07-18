package status

import (
	"fmt"
	"log/slog"
)

// Make sure to run Complete() afterwards
func (status *Status) Update(message string) {

	if status.isTTY {

		status.writeMutex.Lock()
		if _, err := fmt.Fprint(status.out, "\r\033[K"+message); err != nil { // '\r' return to start of line, '\033[K' delete line -> ASCII escape characters
			slog.Debug("status: Error writing to output", "error", err)
		}
		status.writeMutex.Unlock()

	}
}
