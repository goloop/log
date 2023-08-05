package level

import (
	"errors"
	"fmt"
	"math/bits"
)

const (
	// Panic is the panic-type logging level.
	Panic Level = 1 << iota

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

	// The overflowLevelValue is a exceeding the limit of permissible
	// values for the Level.
	overflowLevelValue Level = (1 << iota)

	// Default is the default logging level.
	Default = Panic | Fatal | Error | Warn | Info | Debug | Trace
)

// Labels associates human-readable headings with log levels.
var Labels = map[Level]string{
	Panic: "PANIC",
	Fatal: "FATAL",
	Error: "ERROR",
	Warn:  "WARNING",
	Info:  "INFO",
	Debug: "DEBUG",
	Trace: "TRACE",
}

// ColorLabels associates human-readable headings with log levels.
// See more: https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
var ColorLabels = map[Level]string{
	Panic: fmt.Sprintf("\x1b[5m\x1b[1m\x1b[31m%s\x1b[0m", "PANIC"),
	Fatal: fmt.Sprintf("\x1b[1m\x1b[31m%s\x1b[0m", "FATAL"),
	Error: fmt.Sprintf("\x1b[31m%s\x1b[0m", "ERROR"),
	Warn:  fmt.Sprintf("\x1b[2m\x1b[1m\x1b[31m%s\x1b[0m", "WARNING"),
	Info:  fmt.Sprintf("\x1b[2m\x1b[30m%s\x1b[0m", "INFO"),
	Debug: fmt.Sprintf("\x1b[1m\x1b[32m%s\x1b[0m", "DEBUG"),
	Trace: fmt.Sprintf("\x1b[1m\x1b[33m%s\x1b[0m", "TRACE"),
}

// Level is the type of the level flags.
type Level uint8

// IsSingle returns true if value contains single of the available flag.
func (l *Level) IsSingle() bool {
	return bits.OnesCount(uint(*l)) == 1 &&
		*l <= Level(overflowLevelValue+1)>>1
}

// Contains method returns true if value contains the specified flag.
// Returns false and an error if the value is invalid or an
// invalid flag is specified.
func (l *Level) Contains(flag Level) (bool, error) {
	switch {
	case !flag.IsValid():
		return false, errors.New("incorrect flag value")
	case !l.IsValid():
		return false, errors.New("the object is damaged")
	}

	return *l&flag == flag, nil // *l&flag != 0, nil
}

// IsValid returns true if value contains zero, one or an
// unique sum of valid LevelFlag flags. The zero value is a valid value.
func (l *Level) IsValid() bool {
	// Check if object is zero, which is a valid value.
	if *l == 0 {
		return true
	}

	copy := *l
	// Iterate over all possible values of the constants and
	// check whether they are part of object.
	for level := Level(1); level < overflowLevelValue; level <<= 1 {
		// If layout is part of the object, remove it from object.
		if copy&level == level {
			copy ^= level
		}
	}

	// Check whether all bits of t were "turned off".
	// If t is zero, it means that all bits were matched values
	// of constants, and therefore t is valid.
	return copy == 0
}

// Panic returns true if value contains the Panic flag.
func (l *Level) Panic() bool {
	v, _ := l.Contains(Panic)
	return v
}

// Fatal returns true if value contains the Fatal flag.
func (l *Level) Fatal() bool {
	v, _ := l.Contains(Fatal)
	return v
}

// Error returns true if value contains the Error flag.
func (l *Level) Error() bool {
	v, _ := l.Contains(Error)
	return v
}

// Info returns true if value contains the Info flag.
func (l *Level) Info() bool {
	v, _ := l.Contains(Info)
	return v
}

// Debug returns true if value contains the Debug flag.
func (l *Level) Debug() bool {
	v, _ := l.Contains(Debug)
	return v
}

// Trace returns true if value contains the Trace flag.
func (l *Level) Trace() bool {
	v, _ := l.Contains(Trace)
	return v
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (l *Level) Set(flags ...Level) (Level, error) {
	var r Level

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Contains(flag); !ok {
			r += Level(flag)
		}
	}

	*l = r
	return *l, nil
}

// Add adds the specified flags ignores duplicates or flags that value
// already contains. Returns a new value if all is well or old value and
// an error if one or more invalid flags are specified.
func (l *Level) Add(flags ...Level) (Level, error) {
	r := *l

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Contains(flag); !ok {
			r += Level(flag)
		}
	}

	*l = r
	return *l, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (l *Level) Delete(flags ...Level) (Level, error) {
	r := *l

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Contains(flag); ok {
			r -= Level(flag)
		}
	}

	*l = r
	return *l, nil
}

// All returns true if all of the specified flags are set.
func (l *Level) All(flags ...Level) bool {
	for _, flag := range flags {
		if ok, _ := l.Contains(flag); !ok {
			return false
		}
	}

	return true
}

// Any returns true if at least one of the specified flags is set.
func (l *Level) Any(flags ...Level) bool {
	for _, flag := range flags {
		if ok, _ := l.Contains(flag); ok {
			return ok
		}
	}

	return false
}
