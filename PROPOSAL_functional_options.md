# Proposal: Unified Functional Options via Generics

## Problem

Six scan packages each declare an identical pattern independently:

```go
// repeated in dnsScan, portScan, ipScan, httpHeaderScan, htmlContentScan, httpProtocolScan

type scanConfig struct { ... }

type ConfigOption func(*scanConfig)   // ← identical structure, different type

func WithFoo(v T) ConfigOption {
    return func(sc *scanConfig) { sc.foo = v }
}

func Scan(..., options ...ConfigOption) (Result, error) {
    config := &scanConfig{}
    for _, option := range options {  // ← identical apply loop, six times
        option(config)
    }
    ...
}
```

Every new scan package copies this boilerplate and every new option adds a `WithXxx` function. There is no shared definition of what a `ConfigOption` is — it is just a convention.

## Proposed Solution

Use Go generics to define the option type and apply helper once in `pkg/types`:

```go
// pkg/types/option.go

// Option is a functional option that configures a value of type C.
type Option[C any] func(*C)

// ApplyOptions applies all options to cfg in order.
func ApplyOptions[C any](cfg *C, opts []Option[C]) {
    for _, opt := range opts {
        opt(cfg)
    }
}
```

Each scan package then:

1. Replaces `type ConfigOption func(*scanConfig)` with a type alias:
   ```go
   type ConfigOption = types.Option[scanConfig]
   ```

2. Replaces the apply loop with a single call:
   ```go
   config := &scanConfig{}
   types.ApplyOptions(config, options)
   ```

3. `WithXxx` functions keep their current signatures unchanged since `ConfigOption` is still exported and fully named.

### Example after change (`dnsScan`):

```go
type scanConfig struct {
    nameserver      string
    advanced        bool
    followRedirects bool
    ...
}

type ConfigOption = types.Option[scanConfig]

func WithAdvanced(advanced bool) ConfigOption {
    return func(sc *scanConfig) { sc.advanced = advanced }
}

func Scan(target types.Target, status *status.Status, options ...ConfigOption) (Result, error) {
    config := &scanConfig{}
    types.ApplyOptions(config, options)
    ...
}
```

No changes to call sites (`Scan(..., dnsScan.WithAdvanced(true))`).

## Trade-offs

| | Current | Proposed |
|---|---|---|
| `ConfigOption` definition | 1 line × 6 packages | shared in `types` + 1-line alias × 6 |
| Apply loop | 3 lines × 6 packages | 1 line × 6 packages |
| Discoverability | `ConfigOption` defined near usage | requires knowing about `types.Option` |
| Future packages | copy-paste required | just alias + `ApplyOptions` call |
| Caller API | unchanged | unchanged |
| Go version required | any | 1.18+ (generics) — already met: go.mod says 1.26 |

The savings per package are modest (~4 lines). The main value is having a **single canonical definition** of what a functional option is, so future scan packages follow the same pattern without copy-pasting, and the apply loop can never diverge (e.g., applying options in reverse, skipping nil checks, etc.).

## Scope

Affected packages: `dnsScan`, `portScan`, `ipScan`, `httpHeaderScan`, `htmlContentScan`, `httpProtocolScan`, `seoScan`.

Unaffected: all call sites, all `With*` function signatures, all `Scan` function signatures.
