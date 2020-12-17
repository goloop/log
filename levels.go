package log

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

	// The maxLevelsValue is a special flag that indicating the
	// maximum allowed for Levels type.
	maxLevelsValue Levels = (1 << iota) - 1
)

// The LevelNames associates human-readable headings with log levels.
var LevelNames = map[Level]string{
	Panic: "PANIC",
	Fatal: "FATAL",
	Error: "ERROR",
	Warn:  "WARNING",
	Info:  "INFO",
	Debug: "DEBUG",
	Trace: "TRACE",
}

// Level is the type of single flags of the the Levels.
type Level uint8

// The IsValid returns true if value contains one of the available flags.
// The custom flags cannot be valid since they should not affect the
// formatting settings. The zero value is an invalid flag too.
func (l *Level) IsValid() bool {
	return bits.OnesCount(uint(*l)) == 1 && *l <= Level(maxLevelsValue+1)>>1
}

// Levels type is designed to control the flags responsible
// for adding in the log message additional information as:
// file path, function name and line number.
type Levels Level

// The Has method returns true if value contains the specified flag.
// Returns false and an error if the value is invalid or an
// invalid flag is specified.
func (l *Levels) Has(flag Level) (bool, error) {
	switch {
	case !flag.IsValid():
		return false, errors.New("incorrect flag value")
	case !l.IsValid():
		return false, errors.New("the object is damaged")
	}

	return *l&Levels(flag) == Levels(flag), nil // *l&flag != 0, nil
}

// IsValid returns true if value contains zero, one or an
// unique sum of valid Level flags. The zero value is a valid value.
func (l *Levels) IsValid() bool {
	return *l <= maxLevelsValue
}

// Panic returns true if value contains the Panic flag.
// Returns false and an error if the value is invalid.
func (l *Levels) Panic() (bool, error) {
	return l.Has(Panic)
}

// Fatal returns true if value contains the Fatal flag.
// Returns false and an error if the value is invalid.
func (l *Levels) Fatal() (bool, error) {
	return l.Has(Fatal)
}

// Error returns true if value contains the Error flag.
// Returns false and an error if the value is invalid.
func (l *Levels) Error() (bool, error) {
	return l.Has(Error)
}

// Info returns true if value contains the Info flag.
// Returns false and an error if the value is invalid.
func (l *Levels) Info() (bool, error) {
	return l.Has(Info)
}

// Debug returns true if value contains the Debug flag.
// Returns false and an error if the value is invalid.
func (l *Levels) Debug() (bool, error) {
	return l.Has(Debug)
}

// Trace returns true if value contains the Trace flag.
// Returns false and an error if the value is invalid.
func (l *Levels) Trace() (bool, error) {
	return l.Has(Trace)
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (l *Levels) Set(flags ...Level) (Levels, error) {
	var r Levels

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); !ok {
			r += Levels(flag)
		}
	}

	*l = r
	return *l, nil
}

// Add adds the specified flags ignores duplicates or flags that value
// already contains. Returns a new value if all is well or old value and
// an error if one or more invalid flags are specified.
func (l *Levels) Add(flags ...Level) (Levels, error) {
	var r = *l

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); !ok {
			r += Levels(flag)
		}
	}

	*l = r
	return *l, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (l *Levels) Delete(flags ...Level) (Levels, error) {
	var r = *l

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); ok {
			r -= Levels(flag)
		}
	}

	*l = r
	return *l, nil
}

// All returns true if all of the specified flags are set.
// Returns false and an error if one or more of the specified flags is invalid.
func (l *Levels) All(flags ...Level) (bool, error) {
	for _, flag := range flags {
		if ok, err := l.Has(flag); !ok || err != nil {
			return false, err
		}
	}

	return true, nil
}

// Any returns true if at least one of the specified flags is set.
// Returns false and an error if one or more of the specified flags is invalid.
func (l *Levels) Any(flags ...Level) (bool, error) {
	for _, flag := range flags {
		if ok, err := l.Has(flag); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}
