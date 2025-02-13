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

type BMap map[string]any

func (m BMap) Field(key string) any {
	return m["fields"].(map[string]any)[key]
}

func Map(buf *bytes.Buffer) BMap {
	m := make(BMap)
	_ = json.Unmarshal(buf.Bytes(), &m)
	buf.Reset()
	return m
}
