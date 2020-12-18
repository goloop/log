package log

// TimestampFormt is default date and time format for a timestamp.
const (
	// FatalStatusCode is default value of the an exit code when
	// calling the Fatal method.
	FatalStatusCode = 1

	// TimestampFormat default value of the time and date format
	// for the timestamp in the log message.
	TimestampFormat = "01.02.2006 15:04:05"

	// SpaceBetweenCells is default value of the string that is set
	// between elements of information blocks in the logging prefix.
	SpaceBetweenCells = " "

	// LevelFormat is default value of the format string
	// of the log level substring.
	LevelFormat = "" // "[%s]"
)

// PrefixConfig is config type for the log message prefix.
type PrefixConfig struct {
	// TimestampFormat defines the time and date format for the
	// timestamp in the log message.
	TimestampFormat string

	// SpaceBetweenCells is string that is set between elements
	// of information blocks in the logging prefix.
	SpaceBetweenCells string

	// LevelFormat is format string of the log level substring.
	LevelFormat string
}

// Config is the type of logging configurations: message display
// parameters, log levels, etc.
type Config struct {
	// Formats is the flag-holder where flags responsible for
	// formatting the log message prefix.
	Formats FormatConfig

	// Levels is the flag-holder where flags responsible for
	// levels of the logging: Panic, Fatal, Error, Warn, Info etc.
	Levels LevelConfig

	// FatalStatusCode is an exit code when calling the Fatal method.
	// Default - 1. If the code is <= 0, the forced exit will not occur.
	FatalStatusCode int

	// Prefix is prefix config of the log-message.
	Prefix *PrefixConfig
}

// FatalAllowed reutrns ture if the exit code for Fatal methot not equal zero.
func (c Config) FatalAllowed() bool {
	return c.FatalStatusCode > 0
}
