package layout

import (
	"errors"
	"fmt"
	"math/bits"
)

const (
	// FullFilePath flag adding in the log message the path to
	// the go-file where the logging method was called.
	FullFilePath Layout = 1 << iota

	// ShortFilePath flag adding in the log message the short path
	// to the go-file where the logging method was called.
	ShortFilePath

	// FuncName flag adding in the log message the function's name
	// where the logging method was called.
	FuncName

	// FuncAddress flag adding in the log message the function's address
	// where the logging method was called.
	FuncAddress

	// LineNumber flag adding in the log message the line number
	// of the go-file where the logging method was called.
	LineNumber

	// The overflowLayoutValue is a exceeding the limit of permissible
	// values for the Layout.
	overflowLayoutValue Layout = (1 << iota) - 1

	// Default is the default format for the log message.
	Default = ShortFilePath | FuncName | LineNumber
)

// Layout is the type of single flags of the the Layout.
type Layout uint8

// IsSingle returns true if value contains single of the available flag.
// The custom flags cannot be valid since they should not affect the
// formatting settings. The zero value is an invalid flag too.
func (f *Layout) IsSingle() bool {
	return bits.OnesCount(uint(*f)) == 1 &&
		*f <= Layout(overflowLayoutValue+1)>>1
}

// Contains method returns true if value contains the specified flag.
// Returns false and an error if the value is invalid or an
// invalid flag is specified.
func (f *Layout) Contains(flag Layout) (bool, error) {
	switch {
	case !flag.IsValid():
		return false, errors.New("incorrect flag value")
	case !f.IsValid():
		return false, errors.New("the object is damaged")
	}

	return *f&Layout(flag) == Layout(flag), nil
}

// IsValid returns true if value contains zero, one or an
// unique sum of valid FormatFlag flags. The zero value is a valid value.
func (f *Layout) IsValid() bool {
	return *f <= overflowLayoutValue
}

// FilePath returns true if value contains the FullPath or ShortPath flags.
// Returns false and an error if the value is invalid.
func (f *Layout) FilePath() bool {
	ffp, err := f.Contains(FullFilePath)
	if err == nil && ffp {
		return true
	}

	sfp, err := f.Contains(ShortFilePath)
	if err == nil && sfp {
		return true
	}

	return false
}

// FullFilePath returns true if value contains the FullPath flag.
func (f *Layout) FullFilePath() bool {
	v, _ := f.Contains(FullFilePath)
	return v
}

// ShortFilePath returns true if value contains the ShortPath flag.
func (f *Layout) ShortFilePath() bool {
	v, _ := f.Contains(ShortFilePath)
	return v
}

// FuncName returns true if value contains the FuncName flag.
func (f *Layout) FuncName() bool {
	v, _ := f.Contains(FuncName)
	return v
}

// FuncAddress returns true if value contains the FuncAddress flag.
func (f *Layout) FuncAddress() bool {
	v, _ := f.Contains(FuncAddress)
	return v
}

// LineNumber returns true if value contains the LineNumber flag.
func (f *Layout) LineNumber() bool {
	v, _ := f.Contains(LineNumber)
	return v
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (f *Layout) Set(flags ...Layout) (Layout, error) {
	var r Layout

	for _, flag := range flags {
		if !flag.IsValid() {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Contains(flag); !ok {
			r += Layout(flag)
		}
	}

	*f = r
	return *f, nil
}

// Add adds the specified flags ignores duplicates or flags that value
// already contains. Returns a new value if all is well or old value and
// an error if one or more invalid flags are specified.
func (f *Layout) Add(flags ...Layout) (Layout, error) {
	r := *f

	for _, flag := range flags {
		if !flag.IsValid() {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Contains(flag); !ok {
			r += Layout(flag)
		}
	}

	*f = r
	return *f, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (f *Layout) Delete(flags ...Layout) (Layout, error) {
	r := *f

	for _, flag := range flags {
		if !flag.IsValid() {
			return *f, fmt.Errorf("the %d is invalid flag value", flag)
		}

		if ok, _ := r.Contains(flag); ok {
			r -= Layout(flag)
		}
	}

	*f = r
	return *f, nil
}

// All returns true if all of the specified flags are set.
func (f *Layout) All(flags ...Layout) bool {
	for _, flag := range flags {
		if ok, _ := f.Contains(flag); !ok {
			return false
		}
	}

	return true
}

// Any returns true if at least one of the specified flags is set.
func (f *Layout) Any(flags ...Layout) bool {
	for _, flag := range flags {
		if ok, _ := f.Contains(flag); ok {
			return true
		}
	}

	return false
}
