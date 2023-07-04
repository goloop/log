// Package log is a custom logging library. It defines a type, Logger,
// with methods for formatting output.  Here's an overview of its
// main features:
//
//   - Log Levels: This package supports various log levels, including Panic,
//     Fatal, Error, Warn, Info, Debug, Trace. Each level has three methods:
//     plain (like Info), formatted (like Infof), and line-based (like Infoln).
//     This distinction allows for flexible log message creation.
//
//   - Flexible Outputs: The logger is designed to output logs to multiple
//     destinations (referred to as "Outputs" in the code).
//
//   - Prefixing: Each logger can have a prefix, which can be set during
//     creation with the New function or later with SetPrefix.
//
//   - Skipping Stack Frames: The package also provides a mechanism to skip
//     a certain number of stack frames when logging, which can be useful when
//     wrapping this logger inside other libraries or utilities.
//
//   - Singleton Pattern: The package uses a singleton pattern to instantiate
//     a single logger instance that can be used throughout an application.
//     It provides global functions that operate on this instance.
//
//   - Panic and Fatal: The Fatal functions call os.Exit(1) after writing
//     the log message. The Panic functions call panic after writing the
//     log message.
package log

/* Default Output Format
MY-APP: 2023/06/26 11:42:08 ERROR .../ground/log/main.go main:13 some text
======  =================== ===== ============================== =========
   |            |             |                 |                  |__ message
   |            |             |                 |_____________________ data
   |            |             |_______________________________________ level
   |            |_____________________________________________________ date and time
   |__________________________________________________________________ prefix (optional) from labels

*/
