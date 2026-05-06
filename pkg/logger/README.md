# logger

## Idea

Provide a unified and simple way for the whole application to print log messages.
It should be easy to use, and capable of displaying nicely formatted output in terminal contexts.

## Usage

```go
import "github.com/thetillhoff/webscan/v3/pkg/logger"

handler := logger.NewHandler(os.Stderr, &writeMutex, &slog.HandlerOptions{Level: slog.LevelInfo}, false)
slog.SetDefault(slog.New(handler))
```

## Implicit Configuration

- Environment variable `$NO_COLOR` is respected automatically by disabling color if set to non-empty value.
- Environment variable `$TERM` is respected automatically by disabling color if set to `dumb`.

## Log Levels

- `-v` for Info level
- `-vvv` for Debug level
- Default: Warn level only
