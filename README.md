[//]: # (!!!Don't modify the README.md, use `make readme` to generate it!!!)


[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/log)](https://goreportcard.com/report/github.com/goloop/log) [![License](https://img.shields.io/badge/license-BSD-blue)](https://github.com/goloop/log/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/log)

*Version: 0.0.10*


# log

The log module implements methods for logging code including various levels of
logging:

    - FATAL
      Fatal represents truly catastrophic situations, as far as
      application is concerned. An application is about to abort
      to prevent some kind of corruption or serious problem,
      if possible. Exit code - 1.

    - ERROR
      An error is a serious issue and represents the failure of
      something important going on in an application. Unlike FATAL,
      the application itself isn't going down the tubes.

    - WARN
      It's log level to indicate that an application might have a
      problem and that theare detected an unusual situation.
      It's unexpected and unusual problem, but no real harm done,
      and it's not known whether the issue will persist or recur.

    - INFO
      This level's messages correspond to normal application
      behavior and milestones. They provide the skeleton of what
      happened.

    - DEBUG
      This level must to include more granular, diagnostic
      information then INFO level.

    - TRACE
      This is really fine-grained information-finer even than DEBUG.
      At this level should capture every detail you possibly can about
      the application's behavior.

## Installation

To install this module use `go get` as:

    $ go get -u github.com/goloop/log

## Quick Start

To use this module import it as:

    package main

    import (
        "github.com/goloop/log"
    )

    type App struct {
        Log *log.Log
    }

    func main() {
        var app = &App{}
        app.Log, _ = log.New()

        app.Log.Levels.Delete(log.TRACE)

        app.Log.Debugln("This information will be shown on the screen")
        app.Log.Tracef("%s\n%s\n", "Trace level was deactivated,",
            "this message willn't be displayed")
    }


## Usage

    const SKIP = 4

SKIP default stack offset values.

    const Timestamp = "01.02.2006 15:04:05"

Timestamp is the format for displaying the time stamp in the log message.

#### type Level

    type Level string


Level identifies the logging level.

    const (
    	FATAL Level = "FATAL"
    	ERROR Level = "ERROR"
    	WARN  Level = "WARNING"
    	INFO  Level = "INFO"
    	DEBUG Level = "DEBUG"
    	TRACE Level = "TRACE"
    )

Allowed log level constants.

#### type Levels

    type Levels map[Level]bool


Levels contains active log levels.

#### func (*Levels) Add

    func (l *Levels) Add(levels ...Level) []Level

Add adds new levels to the list of active logging levels.

#### func (*Levels) All

    func (l *Levels) All(levels ...Level) bool

All returns true if all logging levels are supported.

#### func (*Levels) Any

    func (l *Levels) Any(levels ...Level) bool

Any returns true if any logging level is supported.

#### func (*Levels) Delete

    func (l *Levels) Delete(levels ...Level) []Level

Delete removes the specified logging levels from the list of active logging
levels.

#### func (*Levels) Set

    func (l *Levels) Set(levels ...Level) []Level

Set sets active log levels.

#### type Log

    type Log struct {
    	// Writer is the message receiver object (os.Stdout by default).
    	Writer io.Writer

    	// Timestamp is the format for displaying the
    	// time stamp in the log message.
    	Timestamp string

    	// Levels map of available levels.
    	Levels Levels

    	// ShowFilePath if true appends the full path to the go-file,
    	// the logging method was called.
    	ShowFilePath bool

    	// ShowFuncName if true, appends the function name where the
    	// logging method was called.
    	ShowFuncName bool

    	// ShowFileLine if true appends the line number in the go-file
    	// where the logging method was called.
    	ShowFileLine bool

    	// FatalStatusCode the exit code when calling the Fatal method.
    	// Default - 1. If the code is <= 0, the forced exit will not occur.
    	FatalStatusCode int
    }


Log this is the logging object.

#### func  New

    func New(levels ...Level) (*Log, error)

New returns new Log object. Takes zero or more log levels as arguments. If
logging levels are not specified, all possible logging levels will be activated,
otherwise only the specified logging levels will be activated.

#### func (*Log) Copy

    func (l *Log) Copy() *Log

Copy returns copy of the log object.

#### func (*Log) Debug

    func (l *Log) Debug(a ...interface{}) (n int, err error)

Debug creates message with DEBUG level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Debugf

    func (l *Log) Debugf(format string, a ...interface{}) (n int, err error)

Debugf creates message with DEBUG level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Debugln

    func (l *Log) Debugln(a ...interface{}) (n int, err error)

Debugln creates message with DEBUG, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Error

    func (l *Log) Error(a ...interface{}) (n int, err error)

Error creates message with ERROR level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Errorf

    func (l *Log) Errorf(format string, a ...interface{}) (n int, err error)

Errorf creates message with ERROR level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Errorln

    func (l *Log) Errorln(a ...interface{}) (n int, err error)

Errorln creates message with ERROR, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Fatal

    func (l *Log) Fatal(a ...interface{})

Fatal creates message with FATAL level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. Performs forced exit from the program with status - 1.

#### func (*Log) Fatalf

    func (l *Log) Fatalf(format string, a ...interface{})

Fatalf creates message with FATAL level, according to a format specifier and
writes to log.Writer. Performs forced exit from the program with status - 1.

#### func (*Log) Fatalln

    func (l *Log) Fatalln(a ...interface{})

Fatalln creates message with FATAL, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. Performs forced exit from the program with status - 1.

#### func (*Log) Fdebug

    func (l *Log) Fdebug(w io.Writer, a ...interface{}) (n int, err error)

Fdebug creates message with DEBUG level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Fdebugf

    func (l *Log) Fdebugf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Fdebugf creates message with DEBUG level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fdebugln

    func (l *Log) Fdebugln(w io.Writer, a ...interface{}) (n int, err error)

Fdebugln creates message with DEBUG level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ferror

    func (l *Log) Ferror(w io.Writer, a ...interface{}) (n int, err error)

Ferror creates message with ERROR level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Ferrorf

    func (l *Log) Ferrorf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Ferrorf creates message with ERROR level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ferrorln

    func (l *Log) Ferrorln(w io.Writer, a ...interface{}) (n int, err error)

Ferrorln creates message with ERROR level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ffatal

    func (l *Log) Ffatal(w io.Writer, a ...interface{})

Ffatal creates message with FATAL level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. Performs forced exit from the program with status - 1.

#### func (*Log) Ffatalf

    func (l *Log) Ffatalf(w io.Writer, format string, a ...interface{})

Ffatalf creates message with FATAL level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered. Performs forced exit from the program with status - 1.

#### func (*Log) Ffatalln

    func (l *Log) Ffatalln(w io.Writer, a ...interface{})

Ffatalln creates message with FATAL level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. Performs forced exit from the program with status - 1.

#### func (*Log) Finfo

    func (l *Log) Finfo(w io.Writer, a ...interface{}) (n int, err error)

Finfo creates message with INFO level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Finfof

    func (l *Log) Finfof(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Finfof creates message with INFO level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Finfoln

    func (l *Log) Finfoln(w io.Writer, a ...interface{}) (n int, err error)

Finfoln creates message with INFO level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Format

    func (l *Log) Format(showFilePath, showFuncName, showFileLine bool)

Format sets the message prefix display configuration flags for display: file
path, function name and file line.

#### func (*Log) Ftrace

    func (l *Log) Ftrace(w io.Writer, a ...interface{}) (n int, err error)

Ftrace creates message with TRACE level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Ftracef

    func (l *Log) Ftracef(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Ftracef creates message with TRACE level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ftraceln

    func (l *Log) Ftraceln(w io.Writer, a ...interface{}) (n int, err error)

Ftraceln creates message with TRACE level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fwarn

    func (l *Log) Fwarn(w io.Writer, a ...interface{}) (n int, err error)

Fwarn creates message with WARN level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Fwarnf

    func (l *Log) Fwarnf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Fwarnf creates message with WARN level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fwarnln

    func (l *Log) Fwarnln(w io.Writer, a ...interface{}) (n int, err error)

Fwarnln creates message with WARN level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Info

    func (l *Log) Info(a ...interface{}) (n int, err error)

Info creates message with INFO level, using the default formats for its operands
and writes to log.Writer. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Infof

    func (l *Log) Infof(format string, a ...interface{}) (n int, err error)

Infof creates message with INFO level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Infoln

    func (l *Log) Infoln(a ...interface{}) (n int, err error)

Infoln creates message with INFO, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Trace

    func (l *Log) Trace(a ...interface{}) (n int, err error)

Trace creates message with TRACE level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Tracef

    func (l *Log) Tracef(format string, a ...interface{}) (n int, err error)

Tracef creates message with TRACE level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Traceln

    func (l *Log) Traceln(a ...interface{}) (n int, err error)

Traceln creates message with TRACE, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Warn

    func (l *Log) Warn(a ...interface{}) (n int, err error)

Warn creates message with WARN level, using the default formats for its operands
and writes to log.Writer. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Warnf

    func (l *Log) Warnf(format string, a ...interface{}) (n int, err error)

Warnf creates message with WARN level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Warnln

    func (l *Log) Warnln(a ...interface{}) (n int, err error)

Warnln creates message with WARN, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### type Trace

    type Trace struct {
    	FileLine int
    	FuncName string
    	FilePath string
    }


Trace contains the top-level trace information where the logging method was
called.
