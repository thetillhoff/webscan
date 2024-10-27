package logger

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

func (handler Handler) Handle(ctx context.Context, r slog.Record) error {
	var (
		buf bytes.Buffer

		formatter func(string) string
		firstAttr = true
	)

	buf.WriteString("\r\033[K") // Reset current last line

	buf.WriteString(handler.timeFormatter(r.Time.Format(time.TimeOnly))) // Append formatted time // TODO Dimming

	switch r.Level { // Append level-specific formatted shortCode

	case slog.LevelDebug:
		formatter = handler.DebugFormatter
		buf.WriteString(formatter("  " + shortCodeDebug + "  "))
		buf.WriteString(formatter(r.Message)) // Append message

	case slog.LevelInfo:
		formatter = handler.InfoFormatter
		buf.WriteString(formatter("  " + shortCodeInfo + "  "))
		buf.WriteString(formatter(r.Message)) // Append message

	case slog.LevelWarn:
		formatter = handler.WarnFormatter
		buf.WriteString(formatter("  " + shortCodeWarn + "  "))
		buf.WriteString(formatter(r.Message)) // Append message

	case slog.LevelError:
		formatter = handler.ErrorFormatter
		buf.WriteString(formatter("  " + shortCodeError + "  "))
		buf.WriteString(formatter(r.Message)) // Append message
	}

	r.Attrs(func(attr slog.Attr) bool { // Add attributes

		if firstAttr {

			buf.WriteString(formatter(": " + attr.Key + "=\"" + attr.Value.String() + "\"")) // Append message
			firstAttr = false

		} else {

			buf.WriteString(formatter(", " + attr.Key + "=\"" + attr.Value.String() + "\"")) // Append message

		}

		return true
	})

	buf.WriteString("\n") // Commit by newline

	handler.WriteMutex.Lock()         // Begin atomic write
	defer handler.WriteMutex.Unlock() // Make sure atomic write ends

	fmt.Fprint(os.Stderr, buf.String()) // Write

	return nil
}
