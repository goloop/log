[//]: # (!!!Don't modify the README.md, use `make readme` to generate it!!!)


[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/log)](https://goreportcard.com/report/github.com/goloop/log) [![License](https://img.shields.io/badge/license-BSD-blue)](https://github.com/goloop/log/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/log)

*Version: 0.0.12*


# log

The log module implements methods for logging code including various levels of
logging:

    - Panic
      A panic typically means something went unexpectedly wrong.
      Mostly it use to fail fast on errors that shouldn't occur during
      normal operation, or that we aren't prepared to handle gracefully.

    - Fatal
      Fatal represents truly catastrophic situations, as far as
      application is concerned. An application is about to abort
      to prevent some kind of corruption or serious problem,
      if possible. Exit code - 1.

    - Error
      An error is a serious issue and represents the failure of
      something important going on in an application. Unlike FATAL,
      the application itself isn't going down the tubes.

    - Warn
      It's log level to indicate that an application might have a
      problem and that theare detected an unusual situation.
      It's unexpected and unusual problem, but no real harm done,
      and it's not known whether the issue will persist or recur.

    - Info
      This level's messages correspond to normal application
      behavior and milestones. They provide the skeleton of what
      happened.

    - Debug
      This level must to include more granular, diagnostic
      information then INFO level.

    - Trace
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

        app.Log.Config.Levels.Delete(log.TRACE)

        app.Log.Debugln("This information will be shown on the screen")
        app.Log.Tracef("%s\n%s\n", "Trace level was deactivated,",
            "this message willn't be displayed!")
    }


## Usage

    const (

    	// None means nothing.
    	None = 0
    )


    const TimestampFormat = "01.02.2006 15:04:05"

TimestampFormt is default date and time format for a timestamp.

    var LevelNames = map[LevelFlag]string{
    	Panic: "PANIC",
    	Fatal: "FATAL",
    	Error: "ERROR",
    	Warn:  "WARNING",
    	Info:  "INFO",
    	Debug: "DEBUG",
    	Trace: "TRACE",
    }

The LevelFlagNames associates human-readable headings with log levels.

#### type Config

    type Config struct {
    	// TimestampFormat defines the time and date format for the
    	// timestamp in the log message.
    	TimestampFormat string

    	// Formats is the flag-holder where flags responsible for
    	// formatting the log message prefix.
    	Formats FormatConfig

    	// Levels is the flag-holder where flags responsible for
    	// levels of the logging: Panic, Fatal, Error, Warn, Info etc.
    	Levels LevelConfig

    	// FatalStatusCode is an exit code when calling the Fatal method.
    	// Default - 1. If the code is <= 0, the forced exit will not occur.
    	FatalStatusCode int
    }


Config is the type of logging configurations: message display parameters, log
levels, etc.

#### func (Config) FatalAllowed

    func (c Config) FatalAllowed() bool

FatalAllowed reutrns ture if the exit code for Fatal methot not equal zero.

#### type FormatConfig

    type FormatConfig FormatFlag


FormatConfig type is designed to control the flags responsible for adding in the
log message additional information as: file path, function name and line number.

#### func (*FormatConfig) Add

    func (f *FormatConfig) Add(flags ...FormatFlag) (FormatConfig, error)

Add adds the specified flags ignores duplicates or flags that value already
contains. Returns a new value if all is well or old value and an error if one or
more invalid flags are specified.

#### func (*FormatConfig) All

    func (f *FormatConfig) All(flags ...FormatFlag) (bool, error)

All returns true if all of the specified flags are set. Returns false and an
error if one or more of the specified flags is invalid.

#### func (*FormatConfig) Any

    func (f *FormatConfig) Any(flags ...FormatFlag) (bool, error)

Any returns true if at least one of the specified flags is set. Returns false
and an error if one or more of the specified flags is invalid.

#### func (*FormatConfig) Delete

    func (f *FormatConfig) Delete(flags ...FormatFlag) (FormatConfig, error)

Delete deletes the specified flags ignores duplicates or flags that were not
set. Returns a new value if all is well or old value and an error if one or more
invalid flags are specified.

#### func (*FormatConfig) FilePath

    func (f *FormatConfig) FilePath() (bool, error)

FilePath returns true if value contains the FilePath flag. Returns false and an
error if the value is invalid.

#### func (*FormatConfig) FuncName

    func (f *FormatConfig) FuncName() (bool, error)

FuncName returns true if value contains the FuncName flag. Returns false and an
error if the value is invalid.

#### func (*FormatConfig) Has

    func (f *FormatConfig) Has(flag FormatFlag) (bool, error)

The Has method returns true if value contains the specified flag. Returns false
and an error if the value is invalid or an invalid flag is specified.

#### func (*FormatConfig) IsValid

    func (f *FormatConfig) IsValid() bool

IsValid returns true if value contains zero, one or an unique sum of valid
FormatFlag flags. The zero value is a valid value.

#### func (*FormatConfig) LineNumber

    func (f *FormatConfig) LineNumber() (bool, error)

LineNumber returns true if value contains the LineNumber flag. Returns false and
an error if the value is invalid.

#### func (*FormatConfig) Set

    func (f *FormatConfig) Set(flags ...FormatFlag) (FormatConfig, error)

Set sets the specified flags ignores duplicates. The flags that were set
previously will be discarded. Returns a new value if all is well or old value
and an error if one or more invalid flags are specified.

#### type FormatFlag

    type FormatFlag uint8


FormatFlag is the type of single flags of the the FormatConfig.

    const (
    	// FilePath flag adding in the log message the path to
    	// the go-file where the logging method was called.
    	FilePath FormatFlag = 1 << iota

    	// FuncName flag adding in the log message the function's name
    	// where the logging method was called.
    	FuncName

    	// LineNumber flag adding in the log message the line number
    	// of the go-file where the logging method was called.
    	LineNumber
    )


#### func (*FormatFlag) IsValid

    func (f *FormatFlag) IsValid() bool

The IsValid returns true if value contains one of the available flags. The
custom flags cannot be valid since they should not affect the formatting
settings. The zero value is an invalid flag too.

#### type LevelConfig

    type LevelConfig LevelFlag


LevelConfig type is designed to control the flags responsible for adding in the
log message additional information as: file path, function name and line number.

#### func (*LevelConfig) Add

    func (l *LevelConfig) Add(flags ...LevelFlag) (LevelConfig, error)

Add adds the specified flags ignores duplicates or flags that value already
contains. Returns a new value if all is well or old value and an error if one or
more invalid flags are specified.

#### func (*LevelConfig) All

    func (l *LevelConfig) All(flags ...LevelFlag) (bool, error)

All returns true if all of the specified flags are set. Returns false and an
error if one or more of the specified flags is invalid.

#### func (*LevelConfig) Any

    func (l *LevelConfig) Any(flags ...LevelFlag) (bool, error)

Any returns true if at least one of the specified flags is set. Returns false
and an error if one or more of the specified flags is invalid.

#### func (*LevelConfig) Debug

    func (l *LevelConfig) Debug() (bool, error)

Debug returns true if value contains the Debug flag. Returns false and an error
if the value is invalid.

#### func (*LevelConfig) Delete

    func (l *LevelConfig) Delete(flags ...LevelFlag) (LevelConfig, error)

Delete deletes the specified flags ignores duplicates or flags that were not
set. Returns a new value if all is well or old value and an error if one or more
invalid flags are specified.

#### func (*LevelConfig) Error

    func (l *LevelConfig) Error() (bool, error)

Error returns true if value contains the Error flag. Returns false and an error
if the value is invalid.

#### func (*LevelConfig) Fatal

    func (l *LevelConfig) Fatal() (bool, error)

Fatal returns true if value contains the Fatal flag. Returns false and an error
if the value is invalid.

#### func (*LevelConfig) Has

    func (l *LevelConfig) Has(flag LevelFlag) (bool, error)

The Has method returns true if value contains the specified flag. Returns false
and an error if the value is invalid or an invalid flag is specified.

#### func (*LevelConfig) Info

    func (l *LevelConfig) Info() (bool, error)

Info returns true if value contains the Info flag. Returns false and an error if
the value is invalid.

#### func (*LevelConfig) IsValid

    func (l *LevelConfig) IsValid() bool

IsValid returns true if value contains zero, one or an unique sum of valid
LevelFlag flags. The zero value is a valid value.

#### func (*LevelConfig) Panic

    func (l *LevelConfig) Panic() (bool, error)

Panic returns true if value contains the Panic flag. Returns false and an error
if the value is invalid.

#### func (*LevelConfig) Set

    func (l *LevelConfig) Set(flags ...LevelFlag) (LevelConfig, error)

Set sets the specified flags ignores duplicates. The flags that were set
previously will be discarded. Returns a new value if all is well or old value
and an error if one or more invalid flags are specified.

#### func (*LevelConfig) Trace

    func (l *LevelConfig) Trace() (bool, error)

Trace returns true if value contains the Trace flag. Returns false and an error
if the value is invalid.

#### type LevelFlag

    type LevelFlag uint8


LevelFlag is the type of single flags of the the LevelConfig.

    const (
    	// Panic is the panic-type logging level.
    	Panic LevelFlag = 1 << iota

    	// Fatal is the fatal-type logging level.
    	Fatal

    	// Error is the error-type logging level.
    	Error

    	// Warn is the warning-type logging level.
    	Warn

    	// Info is the information-type logging level.
    	Info

    	// Debug is the debug-type logging level.
    	Debug

    	// Trace is the trace-type logging level.
    	Trace
    )


#### func (*LevelFlag) IsValid

    func (l *LevelFlag) IsValid() bool

The IsValid returns true if value contains one of the available flags. The
custom flags cannot be valid since they should not affect the formatting
settings. The zero value is an invalid flag too.

#### type Log

    type Log struct {

    	// Writer is the message receiver object (os.Stdout by default).
    	Writer io.Writer

    	Config *Config
    }


Log is the logger object.

#### func  New

    func New(flags ...LevelFlag) (*Log, error)

New returns new Log object. Accepts zero or more log-level flags as arguments.
If logging levels are not specified, all possible log-levels will be activated.

#### func (*Log) Copy

    func (l *Log) Copy() *Log

Copy returns copy of the log object.

#### func (*Log) Debug

    func (l *Log) Debug(a ...interface{}) (n int, err error)

Debug creates message with Debug level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Debugf

    func (l *Log) Debugf(format string, a ...interface{}) (n int, err error)

Debugf creates message with Debug level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Debugln

    func (l *Log) Debugln(a ...interface{}) (n int, err error)

Debugln creates message with Debug, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Error

    func (l *Log) Error(a ...interface{}) (n int, err error)

Error creates message with Error level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Errorf

    func (l *Log) Errorf(format string, a ...interface{}) (n int, err error)

Errorf creates message with Error level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Errorln

    func (l *Log) Errorln(a ...interface{}) (n int, err error)

Errorln creates message with Error, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Fatal

    func (l *Log) Fatal(a ...interface{}) (n int, err error)

Fatal creates message with Fatal level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fatalf

    func (l *Log) Fatalf(format string, a ...interface{}) (n int, err error)

Fatalf creates message with Fatal level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fatalln

    func (l *Log) Fatalln(a ...interface{}) (n int, err error)

Fatalln creates message with Fatal, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Fdebug

    func (l *Log) Fdebug(w io.Writer, a ...interface{}) (n int, err error)

Fdebug creates message with Debug level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Fdebugf

    func (l *Log) Fdebugf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Fdebugf creates message with Debug level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fdebugln

    func (l *Log) Fdebugln(w io.Writer, a ...interface{}) (n int, err error)

Fdebugln creates message with Debug level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ferror

    func (l *Log) Ferror(w io.Writer, a ...interface{}) (n int, err error)

Ferror creates message with Error level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Ferrorf

    func (l *Log) Ferrorf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Ferrorf creates message with Error level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ferrorln

    func (l *Log) Ferrorln(w io.Writer, a ...interface{}) (n int, err error)

Ferrorln creates message with Error level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ffatal

    func (l *Log) Ffatal(w io.Writer, a ...interface{}) (n int, err error)

Ffatal creates message with Fatal level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Ffatalf

    func (l *Log) Ffatalf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Ffatalf creates message with Fatal level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ffatalln

    func (l *Log) Ffatalln(w io.Writer, a ...interface{}) (n int, err error)

Ffatalln creates message with Fatal level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Finfo

    func (l *Log) Finfo(w io.Writer, a ...interface{}) (n int, err error)

Finfo creates message with Info level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Finfof

    func (l *Log) Finfof(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Finfof creates message with Info level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Finfoln

    func (l *Log) Finfoln(w io.Writer, a ...interface{}) (n int, err error)

Finfoln creates message with Info level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fpanic

    func (l *Log) Fpanic(w io.Writer, a ...interface{}) (n int, err error)

Fpanic creates message with Panic level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Fpanicf

    func (l *Log) Fpanicf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Fpanicf creates message with Panic level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fpanicln

    func (l *Log) Fpanicln(w io.Writer, a ...interface{}) (n int, err error)

Fpanicln creates message with Panic level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ftrace

    func (l *Log) Ftrace(w io.Writer, a ...interface{}) (n int, err error)

Ftrace creates message with Trace level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Ftracef

    func (l *Log) Ftracef(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Ftracef creates message with Trace level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Ftraceln

    func (l *Log) Ftraceln(w io.Writer, a ...interface{}) (n int, err error)

Ftraceln creates message with Trace level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fwarn

    func (l *Log) Fwarn(w io.Writer, a ...interface{}) (n int, err error)

Fwarn creates message with Warn level, using the default formats for its
operands and writes to w. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Fwarnf

    func (l *Log) Fwarnf(w io.Writer, format string,
    	a ...interface{}) (n int, err error)

Fwarnf creates message with Warn level, according to a format specifier and
writes to w. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Fwarnln

    func (l *Log) Fwarnln(w io.Writer, a ...interface{}) (n int, err error)

Fwarnln creates message with Warn level, using the default formats for its
operands and writes to w. Spaces are always added between operands and a newline
is appended. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Info

    func (l *Log) Info(a ...interface{}) (n int, err error)

Info creates message with Info level, using the default formats for its operands
and writes to log.Writer. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Infof

    func (l *Log) Infof(format string, a ...interface{}) (n int, err error)

Infof creates message with Info level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Infoln

    func (l *Log) Infoln(a ...interface{}) (n int, err error)

Infoln creates message with Info, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Panic

    func (l *Log) Panic(a ...interface{}) (n int, err error)

Panic creates message with Panic level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Panicf

    func (l *Log) Panicf(format string, a ...interface{}) (n int, err error)

Panicf creates message with Panic level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Panicln

    func (l *Log) Panicln(a ...interface{}) (n int, err error)

Panicln creates message with Panic, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Trace

    func (l *Log) Trace(a ...interface{}) (n int, err error)

Trace creates message with Trace level, using the default formats for its
operands and writes to log.Writer. Spaces are added between operands when
neither is a string. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Tracef

    func (l *Log) Tracef(format string, a ...interface{}) (n int, err error)

Tracef creates message with Trace level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Traceln

    func (l *Log) Traceln(a ...interface{}) (n int, err error)

Traceln creates message with Trace, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### func (*Log) Warn

    func (l *Log) Warn(a ...interface{}) (n int, err error)

Warn creates message with Warn level, using the default formats for its operands
and writes to log.Writer. Spaces are added between operands when neither is a
string. It returns the number of bytes written and any write error encountered.

#### func (*Log) Warnf

    func (l *Log) Warnf(format string, a ...interface{}) (n int, err error)

Warnf creates message with Warn level, according to a format specifier and
writes to log.Writer. It returns the number of bytes written and any write error
encountered.

#### func (*Log) Warnln

    func (l *Log) Warnln(a ...interface{}) (n int, err error)

Warnln creates message with Warn, level using the default formats for its
operands and writes to log.Writer. Spaces are always added between operands and
a newline is appended. It returns the number of bytes written and any write
error encountered.

#### type StackSlice

    type StackSlice struct {
    	FileLine int
    	FuncName string
    	FilePath string
    }


StackSlice contains the top-level trace information where the logging method was
called.
