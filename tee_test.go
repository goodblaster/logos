package logos

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogos_Tee_WithLoggers(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	log1 := NewLogger(LevelInfo, JSONFormatter(), buf1)
	log2 := NewLogger(LevelDebug, JSONFormatter(), buf2)

	// Tee log2 to log1
	teeLog := log1.Tee(log2)

	// Info message should go to both (log1 at Info, log2 at Debug accepts Info)
	teeLog.Info("test message")

	m1 := Map(buf1)
	m2 := Map(buf2)
	assert.Equal(t, "test message", m1["msg"])
	assert.Equal(t, "info", m1["level"])
	assert.Equal(t, "test message", m2["msg"])
	assert.Equal(t, "info", m2["level"])
}

func TestLogos_Tee_DifferentLevels(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	buf3 := &bytes.Buffer{}

	mainLog := NewLogger(LevelInfo, JSONFormatter(), buf1)
	infoLog := NewLogger(LevelInfo, JSONFormatter(), buf2)
	debugLog := NewLogger(LevelDebug, JSONFormatter(), buf3)

	teeLog := mainLog.Tee(infoLog, debugLog)

	// Info message: mainLog (Info), infoLog (Info), debugLog (Debug accepts Info)
	teeLog.Info("info message")

	m1 := Map(buf1)
	m2 := Map(buf2)
	m3 := Map(buf3)
	assert.Equal(t, "info message", m1["msg"])
	assert.Equal(t, "info message", m2["msg"])
	assert.Equal(t, "info message", m3["msg"])

	// Debug message: only debugLog (mainLog and infoLog are at Info, which filters Debug)
	teeLog.Debug("debug message")

	assert.Empty(t, buf1.String()) // mainLog filtered it
	assert.Empty(t, buf2.String()) // infoLog filtered it
	m3 = Map(buf3)
	assert.Equal(t, "debug message", m3["msg"]) // debugLog received it
}

func TestLogos_Tee_DifferentFormatters(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	jsonLog := NewLogger(LevelDebug, JSONFormatter(), buf1)
	textLog := NewLogger(LevelDebug, TextFormatter(), buf2)

	teeLog := jsonLog.Tee(textLog)
	teeLog.Info("test message")

	// buf1 should have JSON format
	jsonContent := buf1.String()
	assert.Contains(t, jsonContent, `"level":"info"`)
	assert.Contains(t, jsonContent, `"msg":"test message"`)

	// buf2 should have text format (not JSON)
	textContent := buf2.String()
	assert.NotContains(t, textContent, `"level"`)
	assert.Contains(t, textContent, "info")
	assert.Contains(t, textContent, "test message")
}

func TestLogos_Tee_WithFields(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	mainLog := NewLogger(LevelDebug, JSONFormatter(), buf1).
		With("main_field", "main_value")
	teeLog := NewLogger(LevelDebug, JSONFormatter(), buf2).
		With("tee_field", "tee_value")

	combinedLog := mainLog.Tee(teeLog)
	combinedLog.Info("test")

	m1 := Map(buf1)
	m2 := Map(buf2)

	// Main log should have main_field
	assert.Equal(t, "main_value", m1.Field("main_field"))

	// Tee log should have tee_field
	assert.Equal(t, "tee_value", m2.Field("tee_field"))
}

func TestLogos_Tee_PackageLevel(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	SetDefaultLogger(NewLogger(LevelDebug, JSONFormatter(), buf1))
	teeLogger := NewLogger(LevelDebug, JSONFormatter(), buf2)

	teeLog := Tee(teeLogger)
	teeLog.Info("package tee message")

	m1 := Map(buf1)
	m2 := Map(buf2)
	assert.Equal(t, "package tee message", m1["msg"])
	assert.Equal(t, "package tee message", m2["msg"])
}

func TestLogos_Tee_EmptyLoggers(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	// Tee with no loggers should return original logger
	teeLog := log.Tee()
	teeLog.Info("original message")

	m := Map(buf)
	assert.Equal(t, "original message", m["msg"])
}
