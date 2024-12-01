// Package log provides a flexible, feature-rich logging system for Go applications.
// It offers a comprehensive logging solution with multiple severity levels,
// customizable outputs, formatting options, and concurrent-safe operations.
//
// Key Features:
//
// Log Levels: Supports hierarchical logging levels:
//   - Panic: Logs message then panics
//   - Fatal: Logs message then calls os.Exit(1)
//   - Error: For error conditions
//   - Warn: For warning conditions
//   - Info: For informational messages
//   - Debug: For debugging information
//   - Trace: For very detailed debugging
//
// Multiple Output Formats:
//   - Text format with customizable layouts
//   - JSON format for structured logging
//   - Support for ANSI colors in terminal output
//   - Customizable timestamp formats
//
// Flexible Configuration:
//   - Multiple simultaneous outputs (file, stdout, custom writers)
//   - Per-output configuration (levels, format, prefix visibility)
//   - Custom prefixes for application identification
//   - Adjustable stack frame skipping for wrapper libraries
//
// Output Formatting Options:
//   - File paths (full or short)
//   - Function names and addresses
//   - Line numbers
//   - Custom prefix handling
//   - Timestamp formatting
//   - Level label formatting
//
// Thread Safety:
//   - Concurrent-safe logging operations
//   - Safe for multi-goroutine environments
//
// Basic Usage:
//
//	logger := log.New("APP")                // create new logger with prefix
//	logger.Info("Starting application...")  // simple logging
//	logger.Errorf("Failed: %v", err)        // formatted logging
//	logger.Debugln("Debug information")     // line logging
//
// Output Configuration:
//
//	logger.SetOutputs(
//	    log.Output{
//	        Name:      "console",
//	        Writer:    os.Stdout,
//	        Levels:    log.level.Info | log.level.Debug,
//	        WithColor: true,
//	    },
//	    log.Output{
//	        Name:      "file",
//	        Writer:    fileWriter,
//	        Levels:    log.level.Error | log.level.Fatal,
//	        TextStyle: false,  // JSON output
//	    },
//	)
//
// Default Output Format:
//
//	MY-APP: 2023/06/26 11:42:08 ERROR .../ground/log/main.go main:13 some text
//	======  =================== ===== ============================== =========
//	   |            |             |                 |                  |__ message
//	   |            |             |                 |_____________________ location
//	   |            |             |_______________________________________ level
//	   |            |_____________________________________________________ timestamp
//	   |__________________________________________________________________ prefix
//
// Layout Options:
//   - FullFilePath: Complete path to source file
//   - ShortFilePath: Abbreviated path with configurable sections
//   - FuncName: Name of calling function
//   - FuncAddress: Memory address of calling function
//   - LineNumber: Source code line number
//
// Level-specific Methods:
//
//	Each level (Panic, Fatal, Error, etc.) provides three method variants:
//	- Level(args...): Basic logging with space-separated arguments
//	- Levelf(format, args...): Printf-style formatted logging
//	- Levelln(args...): Line logging with space-separated arguments and newline
//
// Thread Safety:
//
//	All logging operations are protected by mutex locks, making the logger
//	safe for concurrent use across multiple goroutines.
//
// Global Logger:
//
//	  The package provides a default global logger instance accessible through
//	  package-level functions. This instance is initialized during package
//	  initialization and can be used directly:
//
//		log.Info("Using global logger")
//		log.SetPrefix("GLOBAL")
//
// Performance Considerations:
//   - Use appropriate log levels to minimize runtime overhead
//   - Consider using JSON format for structured logging needs
//   - Disable unused log levels in production
//
// Note on Fatal and Panic:
//   - Fatal functions call os.Exit(1) after logging
//   - Panic functions call panic() after logging
//   - Use these levels appropriately based on recovery needs
package log
