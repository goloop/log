package log

import (
	"runtime"
	"strings"
)

// StackSlice contains the top-level trace information
// where the logging method was called.
type StackSlice struct {
	FileLine int
	FuncName string
	FilePath string
}

// The getStackSlice returns the stack slice. The skip argument
// is the number of stack frames to skip before taking a slice.
func getStackSlice(skip int) *StackSlice {
	var ss = &StackSlice{}

	// Return program counters of function invocations on
	// the calling goroutine's stack and skipping function
	// call frames inside *Log.
	pc := make([]uintptr, skip+1) // program counters
	runtime.Callers(skip, pc)

	// Get a function at an address on the stack.
	fn := runtime.FuncForPC(pc[0])

	// Get name, path and line of the file.
	ss.FuncName = fn.Name()
	ss.FilePath, ss.FileLine = fn.FileLine(pc[0])
	if r := strings.Split(ss.FuncName, "."); len(r) > 0 {
		ss.FuncName = r[len(r)-1]
	}

	return ss
}
