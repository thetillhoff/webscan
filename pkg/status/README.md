# status

## Idea

Provide a unified and simple way for the whole application to print information.
It should be easy to use, and capable of displaying a nice terminal-ui in all contexts.

Important: It needs to be compatible with logging, so logs and status messages are atomic.

## Usage

```go
status := status.NewStatus(false, false)

status.Update("working on a")
time.Sleep(time.Second)
status.Update("working on b")
time.Sleep(time.Second)
status.Update("working on c")
time.Sleep(time.Second)
status.Complete("complete")

status.SpinningUpdate("working on a")
time.Sleep(3 * time.Second)
status.SpinningUpdate("working on b")
time.Sleep(3 * time.Second)
status.SpinningUpdate("working on c")
time.Sleep(3 * time.Second)
status.SpinningComplete("complete")

status.Println("PRINT")
```

## Implicit Configuration

- Environment variable `$NO_COLOR` is respected automatically by disabling color if set to non-empty value.
- Environment variable `$<MYAPP>_NO_COLOR` is respected automatically by disabling color if set to non-empty value. Value for `<MYAPP>` is retrieved by os.Args[0] put to uppercase.
- Environment variable `$TERM` is respected automatically by disabling color and animations if set to `dumb`.

## Methods

### Write updatable line to stderr, update it, complete it and go on with next line

- `status.Update(message)` & `Complete(message)`
- `status.SpinningUpdate(message)` & `SpinningComplete(message)`
  Uses Braille-chars for the spinner: ['⣾', '⣷', '⣯', '⣟', '⡿', '⢿', '⣻', '⣽'] and '✓' for signalling a completed action.
- `status.SpinningXOfInit(total, message)`, `status.SpinningXOfUpdate()` & `status.SpinningXOfComplete()`

> Make sure to follow `Update` with `Complete`, `SpinningUpdate` with `SpinningComplete`, and `SpinningXOfInit` with `SpinningXOfComplete`, respectively.
