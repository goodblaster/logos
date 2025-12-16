package logos

import (
	"io"
	"testing"
)

// BenchmarkBasicLogging tests basic logging performance
func BenchmarkBasicLogging(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("test message")
	}
}

// BenchmarkLoggingWithFields tests logging with fields
func BenchmarkLoggingWithFields(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard).
		WithFields(Fields{
			"user_id":    12345,
			"request_id": "abc-def-ghi",
			"endpoint":   "/api/users",
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("request processed")
	}
}

// BenchmarkLoggingFiltered tests performance when logs are filtered out
func BenchmarkLoggingFiltered(b *testing.B) {
	log := NewLogger(LevelError, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Debug("this should be filtered")
	}
}

// BenchmarkJSONFormatter tests JSON formatter performance
func BenchmarkJSONFormatter(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("test message")
	}
}

// BenchmarkTextFormatter tests Text formatter performance
func BenchmarkTextFormatter(b *testing.B) {
	log := NewLogger(LevelInfo, TextFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("test message")
	}
}

// BenchmarkConsoleFormatter tests Console formatter performance
func BenchmarkConsoleFormatter(b *testing.B) {
	log := NewLogger(LevelInfo, ConsoleFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("test message")
	}
}

// BenchmarkLogf tests formatted logging performance
func BenchmarkLogf(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Infof("test message %d %s", i, "formatted")
	}
}

// BenchmarkLogFunc tests lazy evaluation performance
func BenchmarkLogFunc(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.LogFunc(LevelInfo, func() string {
			return "expensive computation"
		})
	}
}

// BenchmarkLogFuncFiltered tests lazy evaluation when filtered
func BenchmarkLogFuncFiltered(b *testing.B) {
	log := NewLogger(LevelError, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.LogFunc(LevelDebug, func() string {
			return "this should never be called"
		})
	}
}

// BenchmarkWithChaining tests performance of chaining With() calls
func BenchmarkWithChaining(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.With("key1", "value1").
			With("key2", "value2").
			With("key3", "value3").
			Info("test message")
	}
}

// BenchmarkTeeLogging tests performance with tee loggers
func BenchmarkTeeLogging(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard).
		Tee(
			NewLogger(LevelDebug, TextFormatter(), io.Discard),
			NewLogger(LevelWarn, ConsoleFormatter(), io.Discard),
		)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("test message")
	}
}

// BenchmarkConcurrentLogging tests concurrent logging performance
func BenchmarkConcurrentLogging(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("concurrent message")
		}
	})
}

// BenchmarkConcurrentLoggingWithFields tests concurrent logging with fields
func BenchmarkConcurrentLoggingWithFields(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			log.With("goroutine_id", i).Info("concurrent message")
			i++
		}
	})
}

// BenchmarkIsLevelEnabled tests IsLevelEnabled performance
func BenchmarkIsLevelEnabled(b *testing.B) {
	log := NewLogger(LevelInfo, JSONFormatter(), io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if log.IsLevelEnabled(LevelDebug) {
			log.Debug("test")
		}
	}
}
