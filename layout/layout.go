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
	overflowLayoutValue Layout = (1 << iota)

	// Default is the default format for the log message.
	Default = ShortFilePath | FuncName | LineNumber
)

// Layout is the type of single flags of the the Layout.
type Layout uint8

// IsSingle returns true if value contains single of the available flag.
// The custom flags cannot be valid since they should not affect the
// formatting settings. The zero value is an invalid flag too.
func (l *Layout) IsSingle() bool {
	return bits.OnesCount(uint(*l)) == 1 &&
		*l <= Layout(overflowLayoutValue+1)>>1
}

// Contains method returns true if value contains the specified flag.
// Returns false and an error if the value is invalid or an
// invalid flag is specified.
func (l *Layout) Contains(flag Layout) (bool, error) {
	switch {
	case !flag.IsValid():
		return false, errors.New("incorrect flag value")
	case !l.IsValid():
		return false, errors.New("the object is damaged")
	}

	return *l&Layout(flag) == Layout(flag), nil
}

// IsValid returns true if value contains zero, one or an
// unique sum of valid FormatFlag flags. The zero value is a valid value.
func (l *Layout) IsValid() bool {
	// Check if object is zero, which is a valid value.
	if *l == 0 {
		return true
	}

	copy := *l
	// Iterate over all possible values of the constants and
	// check whether they are part of object.
	for layout := Layout(1); layout < overflowLayoutValue; layout <<= 1 {
		// If layout is part of the object, remove it from object.
		if copy&layout == layout {
			copy ^= layout
		}
	}

	// Check whether all bits of t were "turned off".
	// If t is zero, it means that all bits were matched values
	// of constants, and therefore t is valid.
	return copy == 0
}

// FilePath returns true if value contains the FullPath or ShortPath flags.
// Returns false and an error if the value is invalid.
func (l *Layout) FilePath() bool {
	ffp, err := l.Contains(FullFilePath)
	if err == nil && ffp {
		return true
	}

	sfp, err := l.Contains(ShortFilePath)
	if err == nil && sfp {
		return true
	}

	return false
}

// FullFilePath returns true if value contains the FullPath flag.
func (l *Layout) FullFilePath() bool {
	v, _ := l.Contains(FullFilePath)
	return v
}

// ShortFilePath returns true if value contains the ShortPath flag.
func (l *Layout) ShortFilePath() bool {
	v, _ := l.Contains(ShortFilePath)
	return v
}

// FuncName returns true if value contains the FuncName flag.
func (l *Layout) FuncName() bool {
	v, _ := l.Contains(FuncName)
	return v
}

// FuncAddress returns true if value contains the FuncAddress flag.
func (l *Layout) FuncAddress() bool {
	v, _ := l.Contains(FuncAddress)
	return v
}

// LineNumber returns true if value contains the LineNumber flag.
func (l *Layout) LineNumber() bool {
	v, _ := l.Contains(LineNumber)
	return v
}

// Set sets the specified flags ignores duplicates.
// The flags that were set previously will be discarded.
// Returns a new value if all is well or old value and an
// error if one or more invalid flags are specified.
func (l *Layout) Set(flags ...Layout) (Layout, error) {
	var r Layout

	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is an invalid flag value", flag)
		}
		r |= flag
	}

	*l = r
	return *l, nil
}

// Add adds the specified flags ignores duplicates or flags that value
// already contains. Returns a new value if all is well or old value and
// an error if one or more invalid flags are specified.
func (l *Layout) Add(flags ...Layout) (Layout, error) {
	r := *l
	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is an invalid flag value", flag)
		}
		r |= flag
	}

	*l = r
	return *l, nil
}

// Delete deletes the specified flags ignores duplicates or
// flags that were not set. Returns a new value if all is well or
// old value and an error if one or more invalid flags are specified.
func (l *Layout) Delete(flags ...Layout) (Layout, error) {
	r := *l
	for _, flag := range flags {
		if !flag.IsValid() {
			return *l, fmt.Errorf("the %d is an invalid flag value", flag)
		}
		r &^= flag
	}

	*l = r
	return *l, nil
}

// All returns true if all of the specified flags are set.
func (l *Layout) All(flags ...Layout) bool {
	for _, flag := range flags {
		if ok, _ := l.Contains(flag); !ok {
			return false
		}
	}

	return true
}

// Any returns true if at least one of the specified flags is set.
func (l *Layout) Any(flags ...Layout) bool {
	for _, flag := range flags {
		if ok, _ := l.Contains(flag); ok {
			return true
		}
	}

	return false
}
