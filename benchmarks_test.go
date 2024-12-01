package log

import (
	"bytes"
	"testing"

	"github.com/goloop/log/level"
)

// Implement test writer that does nothing to avoid I/O overhead in benchmarks
type nopWriter struct{}

func (w *nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Setup functions
func setupLogger() *Logger {
	logger := New()
	logger.SetOutputs(Output{
		Name:   "benchmark",
		Writer: &nopWriter{},
		Levels: level.Default,
	})
	return logger
}

func setupBufferedLogger() (*Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	logger := New()
	logger.SetOutputs(Output{
		Name:   "benchmark",
		Writer: buf,
		Levels: level.Default,
	})
	return logger, buf
}

// Basic logging benchmarks
func BenchmarkLoggerInfo(b *testing.B) {
	logger := setupLogger()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("test message")
	}
}

func BenchmarkLoggerInfof(b *testing.B) {
	logger := setupLogger()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Infof("test message %d", i)
	}
}

func BenchmarkLoggerInfoln(b *testing.B) {
	logger := setupLogger()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Infoln("test message")
	}
}

// Benchmark different log levels
func BenchmarkLogLevels(b *testing.B) {
	logger := setupLogger()

	b.Run("Trace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Trace("test message")
		}
	})

	b.Run("Debug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Debug("test message")
		}
	})

	b.Run("Info", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})

	b.Run("Warn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Warn("test message")
		}
	})

	b.Run("Error", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Error("test message")
		}
	})
}

// Benchmark with different output configurations
func BenchmarkOutputConfigurations(b *testing.B) {
	b.Run("WithPrefix", func(b *testing.B) {
		logger := New("TEST")
		logger.SetOutputs(Output{
			Name:   "benchmark",
			Writer: &nopWriter{},
			Levels: level.Default,
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})

	b.Run("WithColor", func(b *testing.B) {
		logger := New()
		logger.SetOutputs(Output{
			Name:      "benchmark",
			Writer:    &nopWriter{},
			Levels:    level.Default,
			WithColor: 1,
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})

	b.Run("JSONStyle", func(b *testing.B) {
		logger := New()
		logger.SetOutputs(Output{
			Name:      "benchmark",
			Writer:    &nopWriter{},
			Levels:    level.Default,
			TextStyle: -1, // JSON style
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})
}

// Benchmark concurrent logging
func BenchmarkConcurrentLogging(b *testing.B) {
	logger := setupLogger()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("test message")
		}
	})
}

// Benchmark with different message sizes
func BenchmarkMessageSizes(b *testing.B) {
	logger := setupLogger()

	tiny := "t"
	small := "test message"
	medium := "test message with some additional context about what's happening"
	large := string(make([]byte, 1024))

	b.Run("TinyMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info(tiny)
		}
	})

	b.Run("SmallMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info(small)
		}
	})

	b.Run("MediumMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info(medium)
		}
	})

	b.Run("LargeMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info(large)
		}
	})
}

// Benchmark different writers
func BenchmarkWriters(b *testing.B) {
	b.Run("NopWriter", func(b *testing.B) {
		logger := New()
		logger.SetOutputs(Output{
			Name:   "benchmark",
			Writer: &nopWriter{},
			Levels: level.Default,
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})

	b.Run("BufferWriter", func(b *testing.B) {
		logger := New()
		logger.SetOutputs(Output{
			Name:   "benchmark",
			Writer: &bytes.Buffer{},
			Levels: level.Default,
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})

	b.Run("MultipleWriters", func(b *testing.B) {
		logger := New()
		logger.SetOutputs(
			Output{
				Name:   "benchmark1",
				Writer: &nopWriter{},
				Levels: level.Default,
			},
			Output{
				Name:   "benchmark2",
				Writer: &nopWriter{},
				Levels: level.Default,
			},
		)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})
}

// Benchmark logger creation and setup
func BenchmarkLoggerCreation(b *testing.B) {
	b.Run("NewLogger", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New()
		}
	})

	b.Run("NewLoggerWithPrefix", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New("TEST")
		}
	})

	b.Run("NewLoggerWithSetup", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger := New()
			logger.SetOutputs(Output{
				Name:   "benchmark",
				Writer: &nopWriter{},
				Levels: level.Default,
			})
		}
	})
}

// Benchmark layout configurations
func BenchmarkLayoutConfigurations(b *testing.B) {
	msg := "test message"

	b.Run("DefaultLayout", func(b *testing.B) {
		logger := setupLogger()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info(msg)
		}
	})

	b.Run("FullFilePath", func(b *testing.B) {
		logger := New()
		logger.SetOutputs(Output{
			Name:    "benchmark",
			Writer:  &nopWriter{},
			Levels:  level.Default,
			Layouts: 1, // FullFilePath
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info(msg)
		}
	})

	b.Run("AllLayouts", func(b *testing.B) {
		logger := New()
		logger.SetOutputs(Output{
			Name:    "benchmark",
			Writer:  &nopWriter{},
			Levels:  level.Default,
			Layouts: 31, // All layout flags
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info(msg)
		}
	})
}
