package logos

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConsoleFormatter(t *testing.T) {
	cfg := DefaultConfig
	fmtr := NewConsoleFormatter(cfg)
	assert.NotNil(t, fmtr)

	// Default config. Timestamp is close.
	cfg = DefaultConfig
	fmtr = NewConsoleFormatter(cfg)
	line := fmtr.Format(LevelDebug, Entry{Msg: "Test"})
	assert.Equal(t, "\x1b[34mdebug\x1b[0m", strings.Fields(line)[1])
	assert.Equal(t, "Test", strings.Fields(line)[2])
	then, err := time.ParseInLocation(DefaultTimestampFormat, strings.Fields(line)[0], time.Local)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().Local(), then.Local(), time.Second)

	// static timestamp
	static := "2020-01-01T00:00:00"
	cfg = Config{
		Timestamp: func() string {
			return static
		},
	}
	fmtr = NewConsoleFormatter(cfg)
	line = fmtr.Format(LevelPrint, Entry{Msg: "Test"})
	assert.Equal(t, static, strings.Fields(line)[0])
	assert.Equal(t, "\x1b[0mprint\x1b[0m", strings.Fields(line)[1])
	assert.Equal(t, "Test", strings.Fields(line)[2])

	// UTC
	cfg = Config{
		Timestamp: func() string {
			return time.Now().UTC().Format(DefaultTimestampFormat)
		},
	}
	fmtr = NewConsoleFormatter(cfg)
	line = fmtr.Format(LevelInfo, Entry{Msg: "Test"})
	assert.Equal(t, "\x1b[32minfo\x1b[0m", strings.Fields(line)[1])
	assert.Equal(t, "Test", strings.Fields(line)[2])
	then, err = time.ParseInLocation(DefaultTimestampFormat, strings.Fields(line)[0], time.UTC)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().UTC(), then.UTC(), time.Second)

	// With some fields.
	line = fmtr.Format(LevelError, Entry{Msg: "Test", Fields: map[string]any{"key": "value"}})
	assert.Equal(t, "\x1b[31merror\x1b[0m", strings.Fields(line)[1])
	assert.Equal(t, "key=\"value\"", strings.Fields(line)[2])
	assert.Equal(t, "Test", strings.Fields(line)[3])
}

func Test_consoleFormatter_Format(t *testing.T) {
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
		contains []string
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
			contains: []string{
				"",
				"info",
				"Test",
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
				},
			},
			contains: []string{
				"",
				"info",
				"key1=\"value1\"",
				"Test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := consoleFormatter{
				cfg: tt.params.cfg,
			}
			got := f.Format(tt.args.level, Entry{Msg: tt.args.msg, Fields: tt.args.fields})
			fields := strings.Fields(got)
			for i, c := range tt.contains {
				if i < len(fields) {
					assert.Contains(t, fields[i], c)
				} else {
					assert.Fail(t, "Field not found")
				}
			}
		})
	}
}
