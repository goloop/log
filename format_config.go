package log

import (
	"errors"
	"fmt"
	"math/bits"
)

const (
	// FullPathFormat flag adding in the log message the path to
	// the go-file where the logging method was called.
	FullPathFormat FormatFlag = 1 << iota

	// ShortPathFormat flag adding in the log message the short path
	// to the go-file where the logging method was called.
	ShortPathFormat

	// FuncNameFormat flag adding in the log message the function's name
	// where the logging method was called.
	FuncNameFormat

	// LineNumberFormat flag adding in the log message the line number
	// of the go-file where the logging method was called.
	LineNumberFormat

	// The maxFormatConfig is a special flag that indicating the
	// maximum allowed for FormatFlag type.
	maxFormatConfig FormatConfig = (1 << iota) - 1

	// DefaultFormat is the default format for the log message.
	DefaultFormat = ShortPathFormat | FuncNameFormat | LineNumberFormat
)

// FormatFlag is the type of single flags of the the FormatConfig.
type FormatFlag uint8

// IsValid returns true if value contains one of the available flags.
// The custom flags cannot be valid since they should not affect the
// formatting settings. The zero value is an invalid flag too.
func (f *FormatFlag) IsValid() bool {
	return bits.OnesCount(uint(*f)) == 1 &&
		*f <= FormatFlag(maxFormatConfig+1)>>1
}

// FormatConfig type is designed to control the flags responsible
// for adding in the log message additional information as:
// file path, function name and line number.
type FormatConfig FormatFlag

// Has method returns true if value contains the specified flag.
// Returns false and an error if the value is invalid or an
// invalid flag is specified.
func (f *FormatConfig) Has(flag FormatFlag) (bool, error) {
	switch {
	case !flag.IsValid():
		return false, errors.New("incorrect flag value")
	case !f.IsValid():
		return false, errors.New("the object is damaged")
	}

	return *f&FormatConfig(flag) == FormatConfig(flag), nil
}

// IsValid returns true if value contains zero, one or an
// unique sum of valid FormatFlag flags. The zero value is a valid value.
func (f *FormatConfig) IsValid() bool {
	return *f <= maxFormatConfig
}

// FilePath returns true if value contains the FullPath or ShortPath flags.
// Returns false and an error if the value is invalid.
func (f *FormatConfig) FilePath() (bool, error) {
	v, err := f.Has(FullPathFormat)
	if err == nil && v {
		return true, nil
	}

	return f.Has(ShortPathFormat)
}

// FullPath returns true if value contains the FullPath flag.
// Returns false and an error if the value is invalid.
func (f *FormatConfig) FullPath() (bool, error) {
	return f.Has(FullPathFormat)
}

// ShortPath returns true if value contains the ShortPath flag.
// Returns false and an error if the value is invalid.
func (f *FormatConfig) ShortPath() (bool, error) {
	return f.Has(ShortPathFormat)
}

// FuncName returns true if value contains the FuncName flag.
// Returns false and an error if the value is invalid.
func (f *FormatConfig) FuncName() (bool, error) {
	return f.Has(FuncNameFormat)
}

// LineNumber returns true if value contains the LineNumber flag.
// Returns false and an error if the value is invalid.
func (f *FormatConfig) LineNumber() (bool, error) {
	return f.Has(LineNumberFormat)
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (f *FormatConfig) Set(flags ...FormatFlag) (FormatConfig, error) {
	var r FormatConfig

	for _, flag := range flags {
		if !flag.IsValid() {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); !ok {
			r += FormatConfig(flag)
		}
	}

	*f = r
	return *f, nil
}

// Add adds the specified flags ignores duplicates or flags that value
// already contains. Returns a new value if all is well or old value and
// an error if one or more invalid flags are specified.
func (f *FormatConfig) Add(flags ...FormatFlag) (FormatConfig, error) {
	r := *f

	for _, flag := range flags {
		if !flag.IsValid() {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); !ok {
			r += FormatConfig(flag)
		}
	}

	*f = r
	return *f, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (f *FormatConfig) Delete(flags ...FormatFlag) (FormatConfig, error) {
	r := *f

	for _, flag := range flags {
		if !flag.IsValid() {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Has(flag); ok {
			r -= FormatConfig(flag)
		}
	}

	*f = r
	return *f, nil
}

// All returns true if all of the specified flags are set.
// Returns false and an error if one or more of the specified flags is invalid.
func (f *FormatConfig) All(flags ...FormatFlag) (bool, error) {
	for _, flag := range flags {
		if ok, err := f.Has(flag); !ok || err != nil {
			return false, err
		}
	}

	return true, nil
}

// Any returns true if at least one of the specified flags is set.
// Returns false and an error if one or more of the specified flags is invalid.
func (f *FormatConfig) Any(flags ...FormatFlag) (bool, error) {
	for _, flag := range flags {
		if ok, err := f.Has(flag); ok || err != nil {
			return ok, err
		}
	}

	return false, nil
}
