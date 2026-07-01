[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/log)](https://goreportcard.com/report/github.com/goloop/log) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/log/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/log/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# log

`log` is a leveled, multi-output logging package for Go. A single log call fans
out to every configured **output**, and each output independently decides which
levels it accepts, whether it renders text or JSON, whether it uses colour, and
how the message prefix is laid out.

It is safe for concurrent use, and work is skipped when nobody needs it: level
filtering happens at the source, and stack-frame capture and formatting run only
when at least one output requires them, so silent levels stay cheap.

## Features

- **Seven levels** — `Panic`, `Fatal`, `Error`, `Warn`, `Info`, `Debug`,
  `Trace`, each with plain / `…f` / `…ln` forms.
- **Multiple outputs** — console, files and custom writers at once, each with
  per-output level filtering, text or JSON, ANSI colour and a configurable
  prefix layout.
- **Thread-safe** — mutex-protected, with buffer pooling on the hot path.
- **`slog` bridge** — back the standard `log/slog` with `NewSlog` or
  `logger.Handler()`.
- **Ad-hoc writers** — the `F`-family also tees a message to a one-off writer.
- **Observability** — `Enabled` guards expensive work; `SetErrorHandler`
  surfaces failing outputs.

## Installation

```bash
go get -u github.com/goloop/log/v2
```

```go
import "github.com/goloop/log/v2"
```

Requires Go 1.24 or newer.

## Quick start

```go
package main

import "github.com/goloop/log/v2"

func main() {
    logger := log.New("APP")

    logger.Info("Application started")
    logger.Infof("User %s logged in", "bob")
    logger.Errorln("Failed to connect to database")
}
```

Fan out to a coloured console and a JSON file, each filtered by level:

```go
import (
    "github.com/goloop/log/v2"
    "github.com/goloop/log/v2/layout"
    "github.com/goloop/log/v2/level"
)

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
        Name:      "file",
        Writer:    file,
        Levels:    level.Error | level.Fatal,
        TextStyle: -1, // JSON
    },
)

log.Info("System initialized")
log.Error("Database connection failed")
```

Back the standard `log/slog`:

```go
slogger := log.NewSlog("APP")
slogger.Info("user logged in", "user", "bob", "id", 42)
```

## Documentation

- Full reference and recipes: [DOC.md](DOC.md) · [DOC.UK.md](DOC.UK.md)
- Package API: [pkg.go.dev/github.com/goloop/log/v2](https://pkg.go.dev/github.com/goloop/log/v2)
- Changes between versions: [CHANGELOG.md](CHANGELOG.md)

## Contributing

Contributions are welcome. Please run `go test ./...`, `go vet ./...` and
`gofmt -l .` before submitting a pull request.

## License

`log` is released under the MIT License. See [LICENSE](LICENSE).
