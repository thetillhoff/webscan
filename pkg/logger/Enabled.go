package logger

import (
	"context"
	"log/slog"
)

func (handler Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.slogHandler.Enabled(ctx, level)
}
