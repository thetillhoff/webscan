package logger

import "log/slog"

func (handler Handler) WithGroup(name string) slog.Handler {
	return &Handler{slogHandler: handler.slogHandler.WithGroup(name), buf: handler.buf, WriteMutex: handler.WriteMutex}
}
