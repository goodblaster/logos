package logos

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogos_WithTee(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	log := NewLogger(LevelDebug, JSONFormatter(), buf1)
	teeLog := log.WithTee(buf2)

	// Log should go to both buffers
	teeLog.Info("test message")

	m1 := Map(buf1)
	m2 := Map(buf2)
	assert.Equal(t, "test message", m1["msg"])
	assert.Equal(t, "info", m1["level"])
	assert.Equal(t, "test message", m2["msg"])
	assert.Equal(t, "info", m2["level"])
}

func TestLogos_WithTee_MultipleWriters(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	buf3 := &bytes.Buffer{}

	log := NewLogger(LevelDebug, JSONFormatter(), buf1)
	teeLog := log.WithTee(buf2, buf3)

	// Log should go to all three buffers
	teeLog.Warn("tee message")

	m1 := Map(buf1)
	m2 := Map(buf2)
	m3 := Map(buf3)
	assert.Equal(t, "tee message", m1["msg"])
	assert.Equal(t, "warn", m1["level"])
	assert.Equal(t, "tee message", m2["msg"])
	assert.Equal(t, "warn", m2["level"])
	assert.Equal(t, "tee message", m3["msg"])
	assert.Equal(t, "warn", m3["level"])
}

func TestLogos_WithTee_PackageLevel(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	SetDefaultLogger(NewLogger(LevelDebug, JSONFormatter(), buf1))
	teeLog := WithTee(buf2)

	// Log should go to both buffers
	teeLog.Info("package tee message")

	m1 := Map(buf1)
	m2 := Map(buf2)
	assert.Equal(t, "package tee message", m1["msg"])
	assert.Equal(t, "package tee message", m2["msg"])
}

func TestLogos_WithTee_EmptyWriters(t *testing.T) {
	buf := &bytes.Buffer{}
	log := NewLogger(LevelDebug, JSONFormatter(), buf)

	// WithTee with no writers should return original logger
	teeLog := log.WithTee()
	teeLog.Info("original message")

	m := Map(buf)
	assert.Equal(t, "original message", m["msg"])
}
