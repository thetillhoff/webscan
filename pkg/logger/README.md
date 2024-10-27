# log

## Idea

Provide a unified and simple way for the whole application to print information.
It should be easy to use, and capable of displaying a nice terminal-ui in all contexts.
And it should be thread-safe.

## Usage

```go
import logger

[...]

var log = logger.NewLogger()
// Change settings as needed
log.Info("hello world")
```

## Implicit Configuration
- Environment variable `$NO_COLOR` is respected automatically by disabling color if set to non-empty value.
- Environment variable `$<MYAPP>_NO_COLOR` is respected automatically by disabling color if set to non-empty value. Value for `<MYAPP>` is retrieved by os.Args[0] put to uppercase.
- Environment variable `$TERM` is respected automatically by disabling color and animations if set to `dumb`.


## Variables

- `log.SpinnerUpdateInterval` - Has a default value.
- `log.LogLevel` - Has a default value.

## Methods

### Write to stdout
- `log.Println(message)`

### Change settings
- `log.NoColor()`

### Write log messages with time, level and message to stderr
- `log.Debug(message)`
- `log.Info(message)`
- `log.Warn(message)`
- `log.Error(message)`
- `log.Fatal(message)`

### Write updatable line to stderr, update it, complete it and go on with next line
- `log.Update(message)` & `Complete(message)`
- `log.SpinningUpdate(message)` & `SpinningComplete(message)`
  Uses Braille-chars for the spinner: ['⣾', '⣷', '⣯', '⣟', '⡿', '⢿', '⣻', '⣽'] and '✓' for signalling a completed action.
- `log.SpinningXOfInit(total, message)`, `log.SpinningXOfUpdate()` & `log.SpinningXOfComplete()`

> Make sure to follow `Update` with `Complete`, `SpinningUpdate` with `SpinningComplete`, and `SpinningXOfInit` with `SpinningXOfComplete`, respectively.
> Doing so also enables you to change 'Downloading X...' to 'Downloaded X' in case of `SpinningUpdateComplete`.


## Todo's

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