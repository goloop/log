package log

import (
	"runtime"
	"strings"
)

// SKIP default stack offset values.
const SKIP = 4

// Trace contains the top-level trace information where the
// logging method was called.
type Trace struct {
	FileLine int
	FuncName string
	FilePath string
}

// The getTrace returns the stack trace of the call to the current function.
// The skip argument is the number of stack frames to skip before taking a cut.
func getTrace(skip int) *Trace {
	var trace = &Trace{}

	// Return program counters of function invocations on
	// the calling goroutine's stack and skipping function
	// call frames inside *Log.
	pc := make([]uintptr, skip+1) // program counters
	runtime.Callers(skip, pc)

	// Get a function at an address on the stack.
	fn := runtime.FuncForPC(pc[0])

	// Get name, path and line of the file.
	trace.FuncName = fn.Name()
	trace.FilePath, trace.FileLine = fn.FileLine(pc[0])
	if r := strings.Split(trace.FuncName, "."); len(r) > 0 {
		trace.FuncName = r[len(r)-1]
	}

	return trace
}
