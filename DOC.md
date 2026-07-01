# log — reference

The full reference for the `log` package: the mental model, outputs and their
configuration, levels and layouts, every logging method, the `slog` bridge and
practical recipes.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Loggers and the default logger](#loggers-and-the-default-logger)
- [Logging methods](#logging-methods)
- [Levels](#levels)
- [Outputs](#outputs)
- [Layouts and formatting](#layouts-and-formatting)
- [Text and JSON output](#text-and-json-output)
- [Managing outputs](#managing-outputs)
- [Ad-hoc writers](#ad-hoc-writers)
- [Structured logging with slog](#structured-logging-with-slog)
- [Conditional logging and error handling](#conditional-logging-and-error-handling)
- [Stack frames and prefixes](#stack-frames-and-prefixes)
- [Recipes and tips](#recipes-and-tips)

## Mental model

`log` is a leveled, multi-output logger. A single log call fans out to every
configured **output**; each output decides — independently — which levels it
accepts, whether it renders text or JSON, whether it uses colour, and how the
message prefix is laid out.

Two ideas shape the API:

1. **One message, many outputs.** You do not create a separate logger per
   destination. You attach several `Output` values to one logger — say a
   coloured console for `Info+` and a JSON file for `Error+` — and every call is
   routed to the ones that want it.
2. **Work is skipped when nobody needs it.** Level filtering happens at the
   source, and stack-frame capture and formatting are done only when at least
   one output requires them, so silent levels are cheap.

The logger is safe for concurrent use, and the hot path pools buffers to keep
allocations low.

```go
import (
    "github.com/goloop/log/v2"
    "github.com/goloop/log/v2/level"
    "github.com/goloop/log/v2/layout"
)
```

## Loggers and the default logger

```go
func New(prefixes ...string) *Logger
func (logger *Logger) Copy() *Logger

func Log() *Logger
func SetDefault(logger *Logger)
```

`New` creates a logger; any prefixes are trimmed and joined with hyphens
(`New("APP", "API")` → `APP-API`) and prepended to each message. `Copy` clones a
logger with its outputs so you can tweak it without touching the original.

The package also keeps a default logger used by the package-level functions
(`log.Info`, `log.Errorf`, …). `Log` returns it; `SetDefault` replaces it (handy
in tests or to install a pre-configured logger at startup).

```go
logger := log.New("APP")
logger.Info("started")

log.Info("uses the default logger")
```

## Logging methods

Each level has three forms, on both the `*Logger` and the package:

| Form | Example | Behaviour |
|------|---------|-----------|
| plain    | `Info(a ...any)`             | like `fmt.Print` |
| `…f`     | `Infof(format string, a ...any)` | like `fmt.Printf` |
| `…ln`    | `Infoln(a ...any)`           | like `fmt.Println` |

The levels, from most to least severe, are `Panic`, `Fatal`, `Error`, `Warn`,
`Info`, `Debug`, `Trace`:

- `Panic*` logs, then calls `panic()`.
- `Fatal*` logs, then calls `os.Exit(1)`.
- The rest log and return.

```go
logger.Info("Application started")
logger.Infof("User %s logged in", user)
logger.Errorln("Failed to connect to database")
```

There is also an `F`-prefixed family (`Finfo`, `Ferrorf`, …) — see
[Ad-hoc writers](#ad-hoc-writers).

## Levels

`level.Level` is a bit-flag set from the `level` subpackage:

```go
const (
    Panic level.Level = 1 << iota
    Fatal
    Error
    Warn
    Info
    Debug
    Trace
)
var Default = Panic | Fatal | Error | Warn | Info | Debug | Trace
```

Combine them with `|` to say which levels an output accepts:

```go
Levels: level.Info | level.Warn | level.Error
```

`Enabled(l)` reports whether any output would accept level `l` — use it to guard
expensive arguments (see [Conditional logging](#conditional-logging-and-error-handling)).

## Outputs

An `Output` is a destination plus its rendering rules. `Name` and `Writer` are
mandatory; the rest have sensible defaults.

| Field | Type | Purpose |
|-------|------|---------|
| `Name`            | `string`         | unique identifier used by `Outputs`/`EditOutputs`/`DeleteOutputs` |
| `Writer`          | `io.Writer`      | where bytes go (`os.Stdout`, a file, a custom writer) |
| `Levels`          | `level.Level`    | which levels this output accepts |
| `Layouts`         | `layout.Layout`  | which caller-context blocks to include in the prefix |
| `Space`           | `string`         | separator between prefix blocks |
| `WithPrefix`      | `trit.Trit`      | show the logger prefix (default on) |
| `WithColor`       | `trit.Trit`      | ANSI colour per level (text only, UNIX-like; default off) |
| `Enabled`         | `trit.Trit`      | enable/disable the output (default on) |
| `TextStyle`       | `trit.Trit`      | text (`true`) vs JSON (`false`); default text |
| `TimestampFormat` | `string`         | `time.Format` layout for the timestamp |
| `LevelFormat`     | `string`         | wrapper around the level name, e.g. `"[%s]"` |

The `trit.Trit` fields use three-valued logic (from the `trit` package): a value
`> 0` is true, `< 0` is false, and `0` means "default" (or "leave unchanged" in
edit mode). You can pass raw `1`/`-1` or `trit.True`/`trit.False`.

Two ready-made outputs are provided: `log.Stdout` and `log.Stderr`.

```go
log.SetOutputs(
    log.Output{
        Name:      "console",
        Writer:    os.Stdout,
        Levels:    level.Info | level.Warn | level.Error,
        Layouts:   layout.Default,
        WithColor: 1,
        TextStyle: 1,
    },
    log.Output{
        Name:       "file",
        Writer:     file,
        Levels:     level.Error | level.Fatal,
        TextStyle:  -1, // JSON
        WithPrefix: 1,
    },
)
```

## Layouts and formatting

`layout.Layout` is a bit-flag set controlling which caller-context blocks appear
in the message prefix:

```go
const (
    FullFilePath layout.Layout = 1 << iota
    ShortFilePath
    FuncName
    FuncAddress
    LineNumber
)
var Default = ShortFilePath | FuncName | LineNumber
```

```go
Layouts: layout.FullFilePath | layout.FuncName | layout.LineNumber
```

`TimestampFormat` and `LevelFormat` further shape the prefix; `Space` sets the
separator between blocks.

## Text and JSON output

An output renders text when `TextStyle` is true (the default) and JSON when it
is false.

Text:

```
APP: 2023/12/02 15:04:05 INFO main.go:42 Starting application
```

JSON (empty fields are omitted):

```json
{
    "prefix": "APP",
    "level": "INFO",
    "timestamp": "2023/12/02 15:04:05",
    "message": "Starting application",
    "filePath": "/home/user/app/main.go",
    "lineNumber": 42,
    "funcName": "main"
}
```

The JSON keys are `prefix`, `level`, `timestamp`, `message`, `filePath`,
`lineNumber`, `funcName` and `funcAddress`.

## Managing outputs

```go
func (logger *Logger) SetOutputs(outputs ...Output) error
func (logger *Logger) EditOutputs(outputs ...Output) error
func (logger *Logger) DeleteOutputs(names ...string)
func (logger *Logger) Outputs(names ...string) []Output
```

`SetOutputs` replaces the whole set. `EditOutputs` changes named outputs in
place — only the fields you set are applied (a `trit` field left at `0` is
untouched), so you can flip colour or levels without redefining the writer.
`DeleteOutputs` removes outputs by name; `Outputs` returns all of them, or just
the named ones.

```go
logger.EditOutputs(log.Output{Name: "console", Levels: level.Error | level.Fatal})
logger.EditOutputs(log.Output{Name: "console", Enabled: -1}) // disable
logger.DeleteOutputs("file")
```

## Ad-hoc writers

The `F`-prefixed methods (`Finfo`, `Ferrorf`, `Fdebugln`, …) write to the
configured outputs **and additionally** to the writer passed as the first
argument, without changing the logger's configuration:

```go
var buf bytes.Buffer
logger.Finfo(&buf, "captured here and in the configured outputs")
```

This is useful for capturing a specific message in a test or request-scoped
buffer while still logging normally.

## Structured logging with slog

The logger can back the standard library's `log/slog`:

```go
func NewSlog(prefixes ...string) *slog.Logger
func (logger *Logger) Handler() slog.Handler
```

```go
slogger := log.NewSlog("APP")
slogger.Info("user logged in", "user", "bob", "id", 42)

// Or attach the handler to an existing logger.
logger := log.New("APP")
slogger = slog.New(logger.Handler())
```

slog levels map onto the logger levels (Debug, Info, Warn, Error). Record
attributes — including those added via `With`/`WithGroup` — become typed JSON
fields in JSON outputs and `key=value` pairs in text outputs.

## Conditional logging and error handling

```go
func (logger *Logger) Enabled(l level.Level) bool
func (logger *Logger) SetErrorHandler(handler func(o Output, n int, err error))
```

Guard expensive arguments with `Enabled` so nothing is computed for a level no
output wants:

```go
if logger.Enabled(level.Debug) {
    logger.Debug(expensiveDump())
}
```

Writes are best-effort by default and write errors are ignored. Register a
handler to observe them — for example to alert on a failing file or network
output:

```go
logger.SetErrorHandler(func(o log.Output, n int, err error) {
    fmt.Fprintf(os.Stderr, "log output %q failed: %v\n", o.Name, err)
})
```

## Stack frames and prefixes

```go
func (logger *Logger) SetSkipStackFrames(skip int) int
func (logger *Logger) SkipStackFrames() int
func (logger *Logger) SetPrefix(prefix string) string
func (logger *Logger) Prefix() string
```

When you wrap the logger behind your own helper, the reported file/line points
at the wrapper. `SetSkipStackFrames` tells the logger how many frames to skip so
the caller's location is reported instead. `SetPrefix`/`Prefix` read or change
the prefix after construction.

```go
logger.SetSkipStackFrames(2) // skip wrapper functions
```

## Recipes and tips

**Console + file split.** Attach a coloured text console for `Info+` and a JSON
file for `Error+`; one call feeds both, each filtered independently.

**Production levels.** Keep `Debug`/`Trace` out of production outputs so the
formatting and stack capture for those levels are skipped entirely; combine with
`Enabled` to avoid building expensive debug payloads.

**Per-request capture.** Use the `F`-family to tee a message into a
request-scoped buffer while normal logging continues.

**Toggle without redefining.** `EditOutputs` applies only the fields you set —
flip `WithColor`, change `Levels`, or disable via `Enabled: -1` without
respecifying the writer.

**Bridge existing slog code.** If your app already logs via `slog`, install
`logger.Handler()` so its records flow through the same outputs, colours and
formats.
