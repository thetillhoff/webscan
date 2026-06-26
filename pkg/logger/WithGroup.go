package logger

import "log/slog"

func (handler Handler) WithGroup(name string) slog.Handler {
	newHandler := handler
	newHandler.slogHandler = handler.slogHandler.WithGroup(name)
	return &newHandler
}
