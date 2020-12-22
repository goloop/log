package log

// LevelFormatConfig is a special type for control of
// the formats in different log level
type LevelFormatConfig map[LevelFlag]string

// Default sets default options for log level formats.
func (lfc *LevelFormatConfig) Default() {
	*lfc = make(map[LevelFlag]string)
}

// Set sets custom formats for log level substring..
func (lfc *LevelFormatConfig) Set(format string) {
	*lfc = make(map[LevelFlag]string)
	for l := range LevelNames {
		(*lfc)[l] = format
	}
}

// Color sets different colors for the substring of the log level.
// Use this method only for write log messages to the *NIX console.
func (lfc *LevelFormatConfig) Color() {
	lfc.Colorf(LevelFormat)
}

// Colorf sets different colors for the substring of the log level
// with support for setting custom formatting for the level substring.
// Use this method only for write log messages to the *NIX console.
func (lfc *LevelFormatConfig) Colorf(format string) {
	*lfc = make(map[LevelFlag]string)

	if len(format) == 0 {
		format = LevelFormat
	}

	// See more: https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
	(*lfc)[Panic] = "\x1b[5m\x1b[1m\x1b[31m" + format + "\x1b[0m"
	(*lfc)[Fatal] = "\x1b[1m\x1b[31m" + format + "\x1b[0m"
	(*lfc)[Error] = "\x1b[31m" + format + "\x1b[0m"
	(*lfc)[Warn] = "\x1b[2m\x1b[1m\x1b[31m" + format + "\x1b[0m"
	(*lfc)[Info] = "\x1b[2m\x1b[30m" + format + "\x1b[0m"
	(*lfc)[Debug] = "\x1b[1m\x1b[32m" + format + "\x1b[0m"
	(*lfc)[Trace] = "\x1b[1m\x1b[33m" + format + "\x1b[0m"
}

// PrefixConfig is config type for the log message prefix.
type PrefixConfig struct {
	// TimestampFormat defines the time and date format for the
	// timestamp in the log message.
	TimestampFormat string

	// SpaceBetweenCells is string that is set between elements
	// of information blocks in the logging prefix.
	SpaceBetweenCells string

	// LevelFormat is format string of the log level substring.
	// The formatting string is specified for each level separately.
	// If no format string is specified for special level, the default
	// format will be used as LevelFormat.
	//
	// Examples:
	//   - Shows the Debug logging level with square brackets:
	//     Log.Config.Prefix.LevelFormat[log.Debug] = "[%s]";
	//   - Shows the Error logging level as red color:
	//     Log.Config.Prefix.LevelFormat[log.Error] = "\033[1;31m%s\033[0m";
	LevelFormat LevelFormatConfig
}
