package log

// Allowed log level constants.
const (
	FATAL Level = "FATAL"
	ERROR Level = "ERROR"
	WARN  Level = "WARNING"
	INFO  Level = "INFO"
	DEBUG Level = "DEBUG"
	TRACE Level = "TRACE"
)

// Level identifies the logging level.
type Level string

// Levels contains active log levels.
type Levels map[Level]bool

// Set sets active log levels.
func (l *Levels) Set(levels ...Level) []Level {
	var result = make(map[Level]bool, len(levels))
	for _, level := range levels {
		result[level] = true
	}

	*l = result
	return activeLevels(*l)
}

// Add adds new levels to the list of active logging levels.
func (l *Levels) Add(levels ...Level) []Level {
	for _, level := range levels {
		(*l)[level] = true
	}

	return activeLevels(*l)
}

// Delete removes the specified logging levels from
// the list of active logging levels.
func (l *Levels) Delete(levels ...Level) []Level {
	for _, level := range levels {
		(*l)[level] = false
	}

	return activeLevels(*l)
}

// All returns true if all logging levels are supported.
func (l *Levels) All(levels ...Level) bool {
	for _, level := range levels {
		if v := (*l)[level]; !v {
			return false
		}
	}

	return true
}

// Any returns true if any logging level is supported.
func (l *Levels) Any(levels ...Level) bool {
	for _, level := range levels {
		if v := (*l)[level]; v {
			return true
		}
	}

	return false
}
