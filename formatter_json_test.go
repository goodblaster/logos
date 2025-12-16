package logos

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/goodblaster/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewJsonFormatter(t *testing.T) {
	// Default config. Timestamp is close.
	cfg := DefaultConfig
	fmtr := NewJsonFormatter(cfg)
	line := fmtr.Format(LevelInfo, Entry{Msg: "Test"})
	var m map[string]any //
	assert.NoError(t, json.Unmarshal([]byte(line), &m))
	then, err := time.ParseInLocation(DefaultTimestampFormat, m["timestamp"].(string), time.Local)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().Local(), then.Local(), time.Second)

	// Custom config. Timestamp is static.
	static := "2020-01-01T00:00:00"
	cfg = Config{
		Timestamp: func() string {
			return static
		},
	}
	fmtr = NewJsonFormatter(cfg)
	line = fmtr.Format(LevelInfo, Entry{Msg: "Test"})
	assert.NoError(t, json.Unmarshal([]byte(line), &m))
	assert.Equal(t, static, m["timestamp"])

	// UTC
	cfg = Config{
		Timestamp: func() string {
			return time.Now().UTC().Format(DefaultTimestampFormat)
		},
	}
	fmtr = NewJsonFormatter(cfg)
	line = fmtr.Format(LevelInfo, Entry{Msg: "Test"})
	assert.NoError(t, json.Unmarshal([]byte(line), &m))
	then, err = time.ParseInLocation(DefaultTimestampFormat, m["timestamp"].(string), time.UTC)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().UTC(), then.UTC(), time.Second)
}

func Test_jsonFormatter_Format(t *testing.T) {
	type params struct {
		cfg Config
	}
	type args struct {
		level  Level
		msg    string
		fields map[string]any
	}
	tests := []struct {
		name     string
		params   params
		args     args
		contains map[string]any
	}{
		{
			name: "Msg only",
			params: params{
				cfg: DefaultConfig,
			},
			args: args{
				level:  LevelInfo,
				msg:    "Test",
				fields: nil,
			},
			contains: map[string]any{
				"msg": "Test",
			},
		},
		{
			name: "Msg with fields",
			params: params{
				cfg: DefaultConfig,
			},
			args: args{
				level: LevelInfo,
				msg:   "Test",
				fields: map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
			contains: map[string]any{
				"msg": "Test",
				"fields": map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := jsonFormatter{
				cfg: tt.params.cfg,
			}
			got := f.Format(tt.args.level, Entry{Msg: tt.args.msg, Fields: tt.args.fields})
			var m map[string]any
			assert.NoError(t, json.Unmarshal([]byte(got), &m))
			for key, value := range tt.contains {
				assert.EqualValues(t, m[key], value)
			}

		})
	}
}

func TestJsonFormatter_ErrorWrapping(t *testing.T) {
	buf := &bytes.Buffer{}
	fmtr := NewJsonFormatter(DefaultConfig)
	log := NewLogger(LevelDebug, fmtr, buf)

	errMsgs := []string{
		"high-level error",
		"wrapped error",
		"base error",
	}
	err := errors.New(errMsgs[2])
	err = errors.Wrap(err, errMsgs[1])
	err = errors.Wrap(err, errMsgs[0])

	log.WithError(err).Error("Test")
	m := Map(buf)

	errList := m.StringList("error")
	assert.EqualValues(t, errMsgs, errList)
}

func TestJsonFormatter_MarshalError(t *testing.T) {
	buf := &bytes.Buffer{}
	fmtr := NewJsonFormatter(DefaultConfig)
	log := NewLogger(LevelDebug, fmtr, buf)

	// Create a circular reference that will cause marshal errors
	type circular struct {
		Self *circular
	}
	circ := &circular{}
	circ.Self = circ

	// This should not panic, but should log an error message instead
	log.With("circular", circ).Info("test with circular reference")

	// Verify we got an error message instead of a panic
	m := Map(buf)
	assert.Equal(t, "[LOG ERROR: failed to marshal entry]", m["msg"])
	assert.Equal(t, "error", m["level"])
}

func TestJsonFormatter_CustomConfig(t *testing.T) {
	buf := &bytes.Buffer{}

	customLevelNames := map[Level]string{
		LevelDebug: "trace",
		LevelInfo:  "information",
		LevelError: "err",
	}

	cfg := Config{
		Timestamp:  func() string { return "2024-01-01" },
		LevelNames: customLevelNames,
	}

	fmtr := NewJsonFormatter(cfg)
	log := NewLogger(LevelDebug, fmtr, buf)

	log.Debug("debug message")
	m := Map(buf)
	assert.Equal(t, "trace", m["level"], "Should use custom level name")
	assert.Equal(t, "2024-01-01", m["timestamp"], "Should use custom timestamp")

	log.Info("info message")
	m = Map(buf)
	assert.Equal(t, "information", m["level"], "Should use custom level name")
}
