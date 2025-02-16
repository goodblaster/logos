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

func (f jsonFormatter) Format(level Level, entry Entry) string {
	type JsonStruct struct {
		Level     string         `json:"level"`
		Timestamp string         `json:"timestamp"`
		Error     error          `json:"error,omitempty"`
		Fields    map[string]any `json:"fields,omitempty"`
		Msg       string         `json:"msg"`
	}

	jsonStruct := JsonStruct{
		Level:     level.String(),
		Timestamp: f.cfg.Timestamp(),
		Error:     entry.Error,
		Msg:       entry.Msg,
		Fields:    entry.Fields,
	}

	b, err := json.Marshal(jsonStruct)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal log entry"))
	}

	return string(b)
}
