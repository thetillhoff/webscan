package logger

import "log/slog"

func (handler Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{slogHandler: handler.slogHandler.WithAttrs(attrs), buf: handler.buf, WriteMutex: handler.WriteMutex}
}
