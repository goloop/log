package log

import (
	"errors"
	"fmt"
	"math/bits"
)

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

	// The maxLevelConfig is a special flag that indicating the
	// maximum allowed for LevelConfig type.
	maxLevelConfig LevelConfig = (1 << iota) - 1
)

// LevelNames associates human-readable headings with log levels.
var LevelNames = map[LevelFlag]string{
	Panic: "PANIC",
	Fatal: "FATAL",
	Error: "ERROR",
	Warn:  "WARNING",
	Info:  "INFO",
	Debug: "DEBUG",
	Trace: "TRACE",
}

// LevelFlag is the type of single flags of the the LevelConfig.
type LevelFlag uint8

// The IsValid returns true if value contains one of the available flags.
// The custom flags cannot be valid since they should not affect the
// formatting settings. The zero value is an invalid flag too.
func (l *LevelFlag) IsValid() bool {
	return bits.OnesCount(uint(*l)) == 1 && *l <= LevelFlag(maxLevelConfig+1)>>1
}

// LevelConfig type is designed to control the flags responsible
// for adding in the log message additional information as:
// file path, function name and line number.
type LevelConfig LevelFlag

// The Has method returns true if value contains the specified flag.
// Returns false and an error if the value is invalid or an
// invalid flag is specified.
func (l *LevelConfig) Has(flag LevelFlag) (bool, error) {
	switch {
	case !flag.IsValid():
		return false, errors.New("incorrect flag value")
	case !l.IsValid():
		return false, errors.New("the object is damaged")
	}

	return *l&LevelConfig(flag) == LevelConfig(flag), nil // *l&flag != 0, nil
}

// IsValid returns true if value contains zero, one or an
// unique sum of valid LevelFlag flags. The zero value is a valid value.
func (l *LevelConfig) IsValid() bool {
	return *l <= maxLevelConfig
}

// Panic returns true if value contains the Panic flag.
// Returns false and an error if the value is invalid.
func (l *LevelConfig) Panic() (bool, error) {
	return l.Has(Panic)
}

// Fatal returns true if value contains the Fatal flag.
// Returns false and an error if the value is invalid.
func (l *LevelConfig) Fatal() (bool, error) {
	return l.Has(Fatal)
}

// Error returns true if value contains the Error flag.
// Returns false and an error if the value is invalid.
func (l *LevelConfig) Error() (bool, error) {
	return l.Has(Error)
}

// Info returns true if value contains the Info flag.
// Returns false and an error if the value is invalid.
func (l *LevelConfig) Info() (bool, error) {
	return l.Has(Info)
}

// Debug returns true if value contains the Debug flag.
// Returns false and an error if the value is invalid.
func (l *LevelConfig) Debug() (bool, error) {
	return l.Has(Debug)
}

// Trace returns true if value contains the Trace flag.
// Returns false and an error if the value is invalid.
func (l *LevelConfig) Trace() (bool, error) {
	return l.Has(Trace)
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (l *LevelConfig) Set(flags ...LevelFlag) (LevelConfig, error) {
	var r LevelConfig

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); !ok {
			r += LevelConfig(flag)
		}
	}

	*l = r
	return *l, nil
}

// Add adds the specified flags ignores duplicates or flags that value
// already contains. Returns a new value if all is well or old value and
// an error if one or more invalid flags are specified.
func (l *LevelConfig) Add(flags ...LevelFlag) (LevelConfig, error) {
	r := *l

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); !ok {
			r += LevelConfig(flag)
		}
	}

	*l = r
	return *l, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (l *LevelConfig) Delete(flags ...LevelFlag) (LevelConfig, error) {
	r := *l

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); ok {
			r -= LevelConfig(flag)
		}
	}

	*l = r
	return *l, nil
}

// All returns true if all of the specified flags are set.
// Returns false and an error if one or more of the specified flags is invalid.
func (l *LevelConfig) All(flags ...LevelFlag) (bool, error) {
	for _, flag := range flags {
		if ok, err := l.Has(flag); !ok || err != nil {
			return false, err
		}
	}

	return true, nil
}

// Any returns true if at least one of the specified flags is set.
// Returns false and an error if one or more of the specified flags is invalid.
func (l *LevelConfig) Any(flags ...LevelFlag) (bool, error) {
	for _, flag := range flags {
		if ok, err := l.Has(flag); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}
