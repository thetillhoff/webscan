package status

import (
	"fmt"
)

// Make sure to run Complete() afterwards
func (status *Status) Update(message string) {

	if !status.quiet && status.isTTY {

		status.writeMutex.Lock()
		fmt.Print("\r\033[K" + message) // '\r' return to start of line, '\033[K' delete line -> ASCII escape characters
		status.writeMutex.Unlock()

	}
}
