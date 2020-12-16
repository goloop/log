package log

import (
	"errors"
	"fmt"
	"math/bits"
)

const (
	// FilePath flag adding in the log message the path to
	// the go-file where the logging method was called.
	FilePath Formats = 1 << iota

	// FuncName flag adding in the log message the function's name
	// where the logging method was called.
	FuncName

	// LineNumber flag adding in the log message the line number
	// of the go-file where the logging method was called.
	LineNumber

	// The maxFormatsValue is a special flag that indicating the
	// maximum allowed for Formats type.
	maxFormatsValue Formats = (1 << iota) - 1
)

// Formats type is designed to control the flags responsible
// for adding in the log message additional information as:
// file path, function name and line number.
type Formats uint8

// The isValid returns true if value in the Formats type range.
func (f *Formats) isValid(value Formats) bool {
	return value <= maxFormatsValue
}

// The has method returns true if the specified flag is set.
// Returns false and an error if the object is invalid or an
// invalid flag is set. Using None as a flag always returns false.
func (f *Formats) has(flag Formats) (bool, error) {
	switch {
	case !f.isValid(flag) || bits.OnesCount(uint(flag)) != 1:
		msg := "incorrect flag value, a single flag must be set"
		return false, errors.New(msg)
	case !f.isValid(*f):
		msg := "the object is damaged, value out of range"
		return false, errors.New(msg)
	}

	return *f&flag == flag, nil // *f&flag != 0, nil
}

// IsValid returns true if object value is valid.
func (f *Formats) IsValid() bool {
	return f.isValid(*f)
}

// FilePath returns true if the FilePath flag is set.
// Returns false if the flag is not set.
// Returns false and an error if the object has an invalid value.
func (f *Formats) FilePath() (bool, error) {
	return f.has(FilePath)
}

// FuncName returns true if the FuncName flag is set.
// Returns false if the flag is not set.
// Returns false and an error if the object has an invalid value.
func (f *Formats) FuncName() (bool, error) {
	return f.has(FuncName)
}

// LineNumber returns true if the LineNumber flag is set.
// Returns false if the flag is not set.
// Returns false and an error if the object has an invalid value.
func (f *Formats) LineNumber() (bool, error) {
	return f.has(LineNumber)
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (f *Formats) Set(flags ...Formats) (Formats, error) {
	var r Formats

	for _, flag := range flags {
		if !f.isValid(flag) {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		// Add only non-existent flags to avoid going
		// out of range of the Formats type.
		if ok, _ := r.has(flag); !ok {
			r += flag
		}
	}

	*f = r
	return *f, nil
}

// Add adds the specified flags ignores duplicates or flags that already set.
// Returns a new value if all is well or old value and an error if one or
// more invalid flags are specified.
func (f *Formats) Add(flags ...Formats) (Formats, error) {
	var r = *f

	for _, flag := range flags {
		if !f.isValid(flag) {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		// Add only non-existent flags to avoid going
		// out of range of the Formats type.
		if ok, _ := r.has(flag); !ok {
			r += flag
		}
	}

	*f = r
	return *f, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (f *Formats) Delete(flags ...Formats) (Formats, error) {
	var r = *f

	for _, flag := range flags {
		if !f.isValid(flag) {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.has(flag); ok {
			r -= flag
		}
	}

	*f = r
	return *f, nil
}

// All returns true if all of the specified flags are set.
// Returns false and an error if one or more of the specified flags is invalid.
func (f *Formats) All(flags ...Formats) (bool, error) {
	for _, flag := range flags {
		if ok, err := f.has(flag); !ok || err != nil {
			return false, err
		}
	}

	return true, nil
}

// Any returns true if at least one of the specified flags is set.
// Returns false and an error if one or more of the specified flags is invalid.
func (f *Formats) Any(flags ...Formats) (bool, error) {
	for _, flag := range flags {
		if ok, err := f.has(flag); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}
