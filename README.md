[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/log)](https://goreportcard.com/report/github.com/goloop/log) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/log/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/log) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# log

The log module encompasses methods for comprehensive logging, which comprises a variety of logging levels:

- Panic
  Panic typically signifies that something has gone unexpectedly awry. It's primarily utilized to swiftly halt on errors that shouldn't surface during regular operation, or those we aren't equipped to handle smoothly.

- Fatal
  Fatal corresponds to situations that are immensely disastrous for the application. The application is on the verge of terminating to avert any sort of corruption or severe problem, if feasible. Exit code is 1.

- Error
  An error represents a significant issue and depicts the failure of something crucial within an application. Contrary to FATAL, the application itself is not doomed.

- Warn
  This log level implies that an application might be experiencing a problem and an unusual situation has been detected. It's an unexpected and unusual issue, but no real damage is done, and it's uncertain whether the problem will persist or happen again.

- Info
  The messages at this level relate to standard application behavior and milestones. They offer a framework of the events that took place.

- Debug
  This level is meant to provide more detailed, diagnostic information than the INFO level.

- Trace
  This level offers incredibly detailed information, even more so than DEBUG. At this level, every possible detail about the application's behavior should be captured.

## Installation

To install this module use `go get` as:

```
$ go get -u github.com/goloop/log
```

## Quick Start

To use this module import it as:

```go
package main

import (
	"os"

	"github.com/goloop/log"
	"github.com/goloop/log/layout"
	"github.com/goloop/log/level"
)

func main() {
	// Open the file in append mode. Create the file if it doesn't exist.
	// Use appropriate permissions in a production setting.
	file, err := os.OpenFile(
		"errors.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Defer the closing of the file until the function ends.
	defer file.Close()

	// Set the outputs of the log to our file.
	// We can set many different outputs to record
	// individual errors, debug or mixed data.
	log.SetOutputs(
		log.Output{
			Name:      "stdout",
			Writer:    os.Stdout,
			Levels:    level.Debug | level.Info | level.Warn | level.Error,
			Layouts:   layout.Default,
			WithColor: 1, // or trit.True, see github.com/goloop/trit
		},
		log.Output{
			Name:    "file-errors",
			Writer:  file,
			Levels:  level.Warn | level.Error, // only errors and warnings
			Layouts: layout.Default,
		},
	)

	// Now, any log messages will be written to the file.
	log.Errorln("This is a test log message with ERROR.")
	log.Debugln("This is a test log message with DEBUG.")
}
```

