package logos

import (
	"encoding/json"

	"github.com/goodblaster/errors"
)

type jsonFormatter struct {
	cfg Config
}

func NewJsonFormatter(cfg Config) Formatter {
	return &jsonFormatter{cfg: cfg}
}

func (f jsonFormatter) Format(level Level, msg string, fields map[string]any) string {
	type Entry struct {
		Level     string         `json:"level"`
		Timestamp string         `json:"timestamp"`
		Fields    map[string]any `json:"fields,omitempty"`
		Msg       string         `json:"msg"`
	}

	entry := Entry{
		Level:     level.String(),
		Timestamp: f.cfg.Timestamp(),
		Msg:       msg,
		Fields:    fields,
	}

	b, err := json.Marshal(entry)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal log entry"))
	}

	return string(b)
}
