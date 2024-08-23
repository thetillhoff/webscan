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


## Variables

- `status.SpinnerUpdateInterval` - Has a default value.
- `status.LogLevel` - Has a default value.

## Methods

### Write to stdout
- `status.Println(message)`

### Change settings
- `status.NoColor()`

### Write updatable line to stderr, update it, complete it and go on with next line
- `status.Update(message)` & `Complete(message)`
- `status.SpinningUpdate(message)` & `SpinningComplete(message)`
  Uses Braille-chars for the spinner: ['⣾', '⣷', '⣯', '⣟', '⡿', '⢿', '⣻', '⣽'] and '✓' for signalling a completed action.
- `status.SpinningXOfInit(total, message)`, `status.SpinningXOfUpdate()` & `status.SpinningXOfComplete()`

> Make sure to follow `Update` with `Complete`, `SpinningUpdate` with `SpinningComplete`, and `SpinningXOfInit` with `SpinningXOfComplete`, respectively.
> Doing so also enables you to change 'Downloading X...' to 'Downloaded X' in case of `SpinningUpdateComplete`.


## Todo's

- Braille-spinner for in-progress
  replace with green, yellow, red ['✓'] on completion - for yellow and red show messages in stderr
- Check compatibility with piping, stdin, redirecting to file

- FAIL level is always enabled
  if FAIL is used, it prints
    error code
    error message
    expected remediation
    url for more information / how to raise a bug


- --json prints json-structured output (zerolog?)
  - level
  - timestamp
  - message
- --verbose, -v for INFO level
  Enables SUC, ERR, WRN, INF messages
- -vvv or $DEBUG for DEBUG level
  Enables SUC, ERR, WRN, INF, DBG messages
- --no-color disables all colors
- --quiet/-q disables statusupdates (don't do anything with Update(), SpinningUpdate() or UpdateXOfTotal() or Complete() or SpinningComplete()) on stderr
- nextSpinner() returns the next element of braille spinner
  needs internal variable with currentIndex
- SpinningUpdate(message) -> writes to internal variable "message" and sets "spinnerEnabled" to true
  must trigger Update() every 0.2s or similar with nextSpinner()
- SpinningComplete(message) -> replaces spinner with (green) checkmark, writes final message, cleans previous message
- Update(message) -> writes to internal variable "message" and sets spinnerEnabled to false
- Complete(finalMessage) -> finishes last line by printing final '$message\n', , cleans previous message
- UpdateXOfTotal(message, total) -> calls Update('( X / total )' + message), where both X and total have the same character length

---

- --instant for showing results immediately after they are available instead of finishing the scan first
- make sure ALL_PROXY, HTTP_PROXY, HTTPS_PROXY, NO_PROXY are respected
- make use of $LINES and $COLUMNS to format output