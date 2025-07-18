package status

import (
	"fmt"
	"log/slog"
)

func (status *Status) Println(message string) {

	if _, err := fmt.Fprintf(status.out, "%s\n", message); err != nil {
		slog.Debug("status: Error writing to output", "error", err)
	}

}
