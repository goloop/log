package log

import (
	"errors"
	"fmt"
	"math/bits"
)

const (
	// Panic is the panic-type logging level.
	Panic Levels = 1 << iota

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

	// The maxLevelsValue is a special flag that indicating the
	// maximum allowed for Levels type.
	maxLevelsValue Levels = (1 << iota) - 1
)

// The levelsCaptions associates human-readable headings with log levels.
var levelCaptions = map[Levels]string{
	Panic: "PANIC",
	Fatal: "FATAL",
	Error: "ERROR",
	Warn:  "WARNING",
	Info:  "INFO",
	Debug: "DEBUG",
	Trace: "TRACE",
}

// Levels type is designed to control the flags responsible
// for activation of log levels.
type Levels uint8

// The isValid returns true if value in the Levels type range.
func (l *Levels) isValid(value Levels) bool {
	return value <= maxLevelsValue
}

// The has method returns true if the specified flag is set.
// Returns false and an error if the object is invalid or an
// invalid flag is set. Using None as a flag always returns false.
func (l *Levels) has(flag Levels) (bool, error) {
	switch {
	case !l.isValid(flag) || bits.OnesCount(uint(flag)) != 1:
		msg := "incorrect flag value, a single flag must be set"
		return false, errors.New(msg)
	case !l.isValid(*l):
		msg := "the object is damaged, value out of range"
		return false, errors.New(msg)
	}

	return *l&flag == flag, nil // *l&flag != 0, nil
}

// IsValid returns true if object value is valid.
func (l *Levels) IsValid() bool {
	return l.isValid(*l)
}

// Panic returns true if the Panic flag is set.
// Returns false if the flag is not set or false
// and an error if the object has an invalid value.
func (l *Levels) Panic() (bool, error) {
	return l.has(Panic)
}

// Fatal returns true if the Fatal flag is set.
// Returns false if the flag is not set or false
// and an error if the object has an invalid value.
func (l *Levels) Fatal() (bool, error) {
	return l.has(Fatal)
}

// Error returns true if the Error flag is set.
// Returns false if the flag is not set or false
// and an error if the object has an invalid value.
func (l *Levels) Error() (bool, error) {
	return l.has(Error)
}

// Warn returns true if the Warn flag is set.
// Returns false if the flag is not set or false
// and an error if the object has an invalid value.
func (l *Levels) Warn() (bool, error) {
	return l.has(Warn)
}

// Info returns true if the Info flag is set.
// Returns false if the flag is not set or false
// and an error if the object has an invalid value.
func (l *Levels) Info() (bool, error) {
	return l.has(Info)
}

// Debug returns true if the Debug flag is set.
// Returns false if the flag is not set or false
// and an error if the object has an invalid value.
func (l *Levels) Debug() (bool, error) {
	return l.has(Debug)
}

// Trace returns true if the Trace flag is set.
// Returns false if the flag is not set or false
// and an error if the object has an invalid value.
func (l *Levels) Trace() (bool, error) {
	return l.has(Trace)
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (l *Levels) Set(flags ...Levels) (Levels, error) {
	var r Levels

	for _, flag := range flags {
		if !l.isValid(flag) {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		// Add only non-existent flags to avoid going
		// out of range of the Levels type.
		if ok, _ := r.has(flag); !ok {
			r += flag
		}
	}

	*l = r
	return *l, nil
}

// Add adds the specified flags ignores duplicates or flags that already set.
// Returns a new value if all is well or old value and an error if one or
// more invalid flags are specified.
func (l *Levels) Add(flags ...Levels) (Levels, error) {
	var r = *l

	for _, flag := range flags {
		if !l.isValid(flag) {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		// Add only non-existent flags to avoid going
		// out of range of the Levels type.
		if ok, _ := r.has(flag); !ok {
			r += flag
		}
	}

	*l = r
	return *l, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (l *Levels) Delete(flags ...Levels) (Levels, error) {
	var r = *l

	for _, flag := range flags {
		if !l.isValid(flag) {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.has(flag); ok {
			r -= flag
		}
	}

	*l = r
	return *l, nil
}

// All returns true if all of the specified flags are set.
// Returns false and an error if one or more of the specified flags is invalid.
func (l *Levels) All(flags ...Levels) (bool, error) {
	for _, flag := range flags {
		if ok, err := l.has(flag); !ok || err != nil {
			return false, err
		}
	}

	return true, nil
}

// Any returns true if at least one of the specified flags is set.
// Returns false and an error if one or more of the specified flags is invalid.
func (l *Levels) Any(flags ...Levels) (bool, error) {
	for _, flag := range flags {
		if ok, err := l.has(flag); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}

// ----------------------------------------------------------

/*
const (
	// Panic is the panic-type logging level.
	Panic Level = "PANIC" // 1

	// Fatal is the fatal-type logging level.
	Fatal Level = "FATAL" // 2

	// Error is the error-type logging level.
	Error Level = "ERROR" // 4

	// Warn is the warning-type logging level.
	Warn Level = "WARNING" // 8

	// Info is the information-type logging level.
	Info Level = "INFO" // 16

	// Debug is the debug-type logging level.
	Debug Level = "DEBUG" // 32

	// Trace is the trace-type logging level.
	Trace Level = "TRACE" // 64
)

// Level is string type alias designed to control log levels.
type Level string

type Levels []Level
*/

/*
// The levelMap is log-level map, intended for storing the states
// of the logging levels (active or not).
type levelMap map[Level]bool

// Set sets active log levels. Returns a list of only active log levels.
func (l *levelMap) Set(levels ...Level) []Level {
	var result = make(map[Level]bool, len(levels))
	for _, level := range levels {
		result[level] = true
	}

	*l = result
	return getActiveLevels(*l)
}

// Add adds new levels to the list of active logging levels.
func (l *levelMap) Add(levels ...Level) []Level {
	for _, level := range levels {
		(*l)[level] = true
	}

	return getActiveLevels(*l)
}

// Delete removes the specified logging levels from
// the list of active logging levels.
func (l *levelMap) Delete(levels ...Level) []Level {
	for _, level := range levels {
		(*l)[level] = false
	}

	return getActiveLevels(*l)
}

// All returns true if all logging levels are supported.
func (l *levelMap) All(levels ...Level) bool {
	for _, level := range levels {
		if v := (*l)[level]; !v {
			return false
		}
	}

	return true
}

// Any returns true if any logging level is supported.
func (l *levelMap) Any(levels ...Level) bool {
	for _, level := range levels {
		if v := (*l)[level]; v {
			return true
		}
	}

	return false
}
*/
