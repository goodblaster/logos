package logos

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogos_ConvenienceFunctions(t *testing.T) {
	buf := &bytes.Buffer{}

	// As debug, all logs should be printed
	log := NewLogger(LevelDebug, JSONFormatter(), buf)
	log.Debug("logos")
	m := Map(buf)
	assert.Equal(t, "logos", m["msg"])
	assert.Equal(t, "debug", m["level"])
	log.Info("logos")
	assert.Equal(t, "info", Map(buf)["level"])
	log.Warn("logos")
	assert.Equal(t, "warn", Map(buf)["level"])
	log.Error("logos")
	assert.Equal(t, "error", Map(buf)["level"])
	log.Print("logos")
	assert.Equal(t, "print", Map(buf)["level"])
}

func TestLogos_Levels(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	for _, level := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPrint} {
		log.Log(level, "logos")
		assert.Equal(t, level.String(), Map(buf)["level"])
	}

	// Change the level to error, only error and print logs should be printed.
	log.SetLevel(LevelError)
	log.Debug("logos")
	assert.Empty(t, buf.String())
	log.Info("logos")
	assert.Empty(t, buf.String())
	log.Warn("logos")
	assert.Empty(t, buf.String())
	log.Error("logos")
	assert.Equal(t, "error", Map(buf)["level"])
	log.Print("logos")
	assert.Equal(t, "print", Map(buf)["level"])
}

func TestLogos_With(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	log.With("key", "value").Log(LevelDebug, "logos")
	assert.Equal(t, "value", Map(buf).Field("key"))
}

func TestLogos_WithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	log.WithFields(map[string]any{"key": "value"}).Log(LevelDebug, "logos")
	assert.Equal(t, "value", Map(buf).Field("key"))
}

func TestLogos_LogFunc(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	log.LogFunc(LevelDebug, func() string {
		return "logos"
	})
	assert.Equal(t, "logos", Map(buf)["msg"])
}

func TestLogos_SetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	log.SetLevel(LevelError)
	log.Debug("logos")
	assert.Empty(t, buf.String())
	log.Error("logos")
	assert.Equal(t, "error", Map(buf)["level"])
}

func TestLogos_CustomLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	const (
		LevelApple Level = iota + 100
		LevelBanana
		LevelCherry
	)

	oldLevels := map[Level]string{}
	for level, name := range LevelNames {
		oldLevels[level] = name
	}

	defer func() {
		LevelNames = oldLevels
	}()

	LevelNames = map[Level]string{
		LevelApple:  "apple",
		LevelBanana: "banana",
		LevelCherry: "cherry",
	}

	log.Log(LevelApple, "apple")
	assert.Equal(t, "apple", Map(buf)["level"])
	log.Log(LevelBanana, "banana")
	assert.Equal(t, "banana", Map(buf)["level"])
	log.Log(LevelCherry, "cherry")
	assert.Equal(t, "cherry", Map(buf)["level"])

	log.SetLevel(LevelBanana)
	log.Log(LevelApple, "apple")
	assert.Empty(t, buf.String())
	log.Log(LevelBanana, "banana")
	assert.Equal(t, "banana", Map(buf)["level"])
	log.Log(LevelCherry, "cherry")
	assert.Equal(t, "cherry", Map(buf)["level"])
}

func TestLogos_DefaultLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)
	SetDefaultLogger(log)

	Log(LevelDebug, "logos")
	assert.Equal(t, "debug", Map(buf)["level"])
}

func TestLogos_SubLoggers(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	subLog := log.With("key", "value")
	sublog2 := subLog.With("key2", "value2")
	sublog3 := sublog2.With("key3", "value3")

	sublog3.Log(LevelDebug, "logos")
	m := Map(buf)
	assert.Equal(t, "value", m.Field("key"))
	assert.Equal(t, "value2", m.Field("key2"))
	assert.Equal(t, "value3", m.Field("key3"))

	sublog2.Log(LevelDebug, "logos")
	m = Map(buf)
	assert.Equal(t, "value", m.Field("key"))
	assert.Equal(t, "value2", m.Field("key2"))
	assert.Empty(t, m.Field("key3"))

	subLog.Log(LevelDebug, "logos")
	m = Map(buf)
	assert.Equal(t, "value", m.Field("key"))
	assert.Empty(t, m.Field("key2"))
	assert.Empty(t, m.Field("key3"))
}

func TestLogger_Getters(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf).
		With("key1", "value1").
		With("key2", "value2").
		WithError(assert.AnError)

	// Test GetLevel
	assert.Equal(t, LevelDebug, log.GetLevel())

	// Test GetFields
	fields := log.GetFields()
	assert.NotNil(t, fields)
	assert.Equal(t, "value1", fields["key1"])
	assert.Equal(t, "value2", fields["key2"])

	// Ensure returned fields are a copy
	fields["key1"] = "modified"
	originalFields := log.GetFields()
	assert.Equal(t, "value1", originalFields["key1"], "GetFields should return a copy")

	// Test GetError
	assert.Equal(t, assert.AnError, log.GetError())

	// Test GetTeeCount
	assert.Equal(t, 0, log.GetTeeCount())

	teeLog := log.Tee(NewLogger(LevelInfo, JSONFormatter(), &bytes.Buffer{}))
	assert.Equal(t, 1, teeLog.GetTeeCount())
}

func TestLogger_WithLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	log1 := NewLogger(LevelDebug, JSONFormatter(), buf)

	// Create a new logger with different level
	log2 := log1.WithLevel(LevelError)

	// Verify they have independent levels
	assert.Equal(t, LevelDebug, log1.GetLevel())
	assert.Equal(t, LevelError, log2.GetLevel())

	// Debug should log on log1 but not log2
	log1.Debug("test1")
	m := Map(buf)
	assert.Equal(t, "test1", m["msg"])

	log2.Debug("test2")
	// Buffer should not have new content since log2 filters debug
	assert.Equal(t, 0, buf.Len())

	// Error should log on log2
	log2.Error("test3")
	m = Map(buf)
	assert.Equal(t, "test3", m["msg"])
}

func TestLogger_Fatal_Panics(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	// Fatal should panic
	assert.Panics(t, func() {
		log.Fatal("fatal message")
	})

	// Fatalf should also panic
	assert.Panics(t, func() {
		log.Fatalf("fatal message %s", "formatted")
	})
}

func TestLogger_Concurrent_Logging(t *testing.T) {
	// The logger now synchronizes writes to prevent concurrent access issues
	// This test verifies that concurrent logging works safely with bytes.Buffer
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	// Run concurrent logging operations
	const numGoroutines = 10
	const numLogs = 100
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numLogs; j++ {
				log.With("goroutine", id).Info("concurrent log")
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify we got logs (should have content from all goroutines)
	assert.Greater(t, buf.Len(), 0, "Should have logged content")
}

func TestLogger_NilLevel(t *testing.T) {
	buf := &bytes.Buffer{}

	// Create a logger with nil level (simulating invalid state)
	log := Logger{
		level:     nil,
		formatter: JSONFormatter(),
		writer:    buf,
	}

	// Should not panic, should just not log
	log.Debug("test")
	assert.Equal(t, 0, buf.Len(), "Should not log with nil level")

	// GetLevel should return safe default
	assert.Equal(t, LevelInfo, log.GetLevel())
}

func TestLogger_NilWriter(t *testing.T) {
	// Create a logger with nil writer
	log := Logger{
		level:     ptrTo(LevelDebug),
		formatter: JSONFormatter(),
		writer:    nil,
	}

	// Should not panic, should just not log
	assert.NotPanics(t, func() {
		log.Debug("test")
	})
}

func TestLogger_NilFormatter(t *testing.T) {
	buf := &bytes.Buffer{}

	// Create a logger with nil formatter
	log := Logger{
		level:     ptrTo(LevelDebug),
		formatter: nil,
		writer:    buf,
	}

	// Should not panic, should just not log
	assert.NotPanics(t, func() {
		log.Debug("test")
	})
}

// Helper function to create a pointer to a Level
func ptrTo(level Level) *Level {
	return &level
}

type BMap map[string]any

func (m BMap) Field(key string) any {
	return m["fields"].(map[string]any)[key]
}

func (m BMap) StringList(key string) []string {
	field, ok := m[key].([]any)
	if !ok {
		return nil
	}

	var result []string
	for _, v := range field {
		if str, ok := v.(string); ok {
			result = append(result, str)
		}
	}

	return result
}

func Map(buf *bytes.Buffer) BMap {
	m := make(BMap)
	_ = json.Unmarshal(buf.Bytes(), &m)
	buf.Reset()
	return m
}
