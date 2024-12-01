[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/log)](https://goreportcard.com/report/github.com/goloop/log) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/log/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/log) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# log

A flexible, high-performance logging package for Go applications with support for multiple output formats, logging levels, and concurrent operations.

## Features

- **Multiple Log Levels** with clear semantics:
  - `Panic`: For unrecoverable errors that require immediate attention (calls `panic()`)
  - `Fatal`: For critical errors that prevent application startup/operation (calls `os.Exit(1)`)
  - `Error`: For runtime errors that need investigation but don't stop the application
  - `Warn`: For potentially harmful situations
  - `Info`: For general operational information
  - `Debug`: For detailed system state information
  - `Trace`: For ultra-detailed debugging information

- **Flexible Output Configuration**:
  - Multiple simultaneous outputs (console, files, custom writers)
  - Per-output level filtering
  - Text and JSON formats
  - ANSI color support for terminal output
  - Custom prefix support
  - Configurable timestamps and layouts

- **Thread-Safe Operations**:
  - Safe for concurrent use across goroutines
  - Mutex-protected logging operations

- **Performance Optimized**:
  - Minimal allocations
  - Efficient formatting
  - Level-based filtering at source

## Installation

```bash
go get -u github.com/goloop/log
```

## Quick Start

### Basic Usage

```go
package main

import "github.com/goloop/log"

func main() {
    // Create logger with prefix.
    logger := log.New("APP")

    // Basic logging.
    logger.Info("Application started")
    logger.Debug("Debug information")
    logger.Error("Something went wrong")

    // Formatted logging.
    logger.Infof("User %s logged in", username)

    // With newline.
    logger.Errorln("Failed to connect to database")
}
```

### Advanced Configuration

```go
package main

import (
    "os"
    "github.com/goloop/log"
    "github.com/goloop/log/layout"
    "github.com/goloop/log/level"
)

func main() {
    // Open log file.
    file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Configure multiple outputs.
    log.SetOutputs(
        // Console output with colors.
        log.Output{
            Name:      "console",
            Writer:    os.Stdout,
            Levels:    level.Info | level.Warn | level.Error,
            Layouts:   layout.Default,
            WithColor: 1,
            TextStyle: 1,
        },
        // File output in JSON format.
        log.Output{
            Name:      "file",
            Writer:    file,
            Levels:    level.Error | level.Fatal,
            TextStyle: -1, // JSON format
            WithPrefix: 1,
        },
    )

    // Use the logger
    log.Info("System initialized")
    log.Error("Database connection failed")
}
```

### Custom Layout Configuration

```go
package main

import (
    "github.com/goloop/log"
    "github.com/goloop/log/layout"
)

func main() {
    logger := log.New("APP")

    // Configure custom layout.
    logger.SetOutputs(log.Output{
        Name:    "custom",
        Writer:  os.Stdout,
        Layouts: layout.FullFilePath | layout.FuncName | layout.LineNumber,
    })

    logger.Info("Custom layout message")
}
```

## Output Formats

### Text Format (Default)
```
APP: 2023/12/02 15:04:05 INFO main.go:42 Starting application
```

### JSON Format
```json
{
    "prefix": "APP",
    "timestamp": "2023/12/02 15:04:05",
    "level": "INFO",
    "file": "main.go",
    "line": 42,
    "message": "Starting application"
}
```

## Performance Considerations

- Use appropriate log levels in production (typically Info and above)
- Consider using JSON format only when structured logging is required
- Disable debug/trace levels in production for optimal performance
- Use formatted logging (`Infof`, etc.) only when necessary

## Advanced Features

### Stack Frame Skipping

```go
logger := log.New("APP")
logger.SetSkipStackFrames(2) // skip wrapper functions
```

### Multiple Prefix Support

```go
logger := log.New("APP", "SERVICE", "API")  // Results in "APP-SERVICE-API"
```

### Custom Writers

```go
type CustomWriter struct {
    // implementation
}

func (w *CustomWriter) Write(p []byte) (n int, err error) {
    // custom write logic.
    return len(p), nil
}

logger.SetOutputs(log.Output{
    Name: "custom",
    Writer: &CustomWriter{},
    Levels: level.Info,
})
```

## Managing Outputs

The logger provides several methods to manage outputs:

### Get Current Outputs
```go
// Get all outputs.
outputs := logger.Outputs()

// Get specific outputs by name.
stdoutOutput := logger.Outputs("stdout")
```

### Edit Outputs
```go
// Change output configuration.
logger.EditOutputs(log.Output{
    Name:    "stdout",
    Levels:  level.Error | level.Fatal,  // change levels
    WithColor: 1,                        // enable colors
})

// Disable specific output.
logger.EditOutputs(log.Output{
    Name:    "stdout",
    Enabled: -1,  // or trit.False
})
```

### Delete Outputs
```go
// Remove specific outputs,
logger.DeleteOutputs("stdout", "file")

// Or disable all logging by removing all outputs.
logger.DeleteOutputs(logger.Outputs()...)
```

### Set New Outputs
```go
// Replace all outputs with new ones.
logger.SetOutputs(
    log.Output{
        Name:    "console",
        Writer:  os.Stdout,
        Levels:  level.Info | level.Warn,
    },
    log.Output{
        Name:    "errors",
        Writer:  errorFile,
        Levels:  level.Error | level.Fatal,
    },
)
```

## Why use this logger?

- Flexible configuration
- High performance
- Multiple output support
- Structured logging support
- Thread safety
- Comprehensive logging levels

## Related Projects

- [goloop/g](https://github.com/goloop/g) - Common utilities
- [goloop/trit](https://github.com/goloop/trit) - Three-valued logic


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


