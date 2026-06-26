package logger

import "log/slog"

func (handler Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := handler
	newHandler.slogHandler = handler.slogHandler.WithAttrs(attrs)
	return &newHandler
}
