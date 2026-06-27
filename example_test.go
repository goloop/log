package log_test

import (
	"log/slog"
	"os"

	"github.com/goloop/log/v2"
	"github.com/goloop/log/v2/level"
)

func ExampleNew() {
	logger := log.New("APP")
	logger.Info("application started")
	logger.Debugf("loaded %d modules", 12)
}

func ExampleLogger_Enabled() {
	logger := log.New("APP")
	if logger.Enabled(level.Debug) {
		logger.Debug("debug is enabled for at least one output")
	}
}

func ExampleLogger_SetOutputs() {
	logger := log.New("APP")
	logger.SetOutputs(log.Output{
		Name:   "console",
		Writer: os.Stdout,
		Levels: level.Info | level.Warn | level.Error,
	})
	logger.Info("configured a single console output")
}

func ExampleNewSlog() {
	slogger := log.NewSlog("APP")
	slogger.Info("user logged in", "user", "bob", "id", 42)
}

func ExampleLogger_Handler() {
	logger := log.New("APP")
	slogger := slog.New(logger.Handler())
	slogger.Warn("disk almost full", "free_mb", 128)
}
