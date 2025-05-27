package logos

import (
	"encoding/json"

	"github.com/goodblaster/errors"
)

// jsonFormatter is a log formatter that outputs logs in JSON format.
type jsonFormatter struct {
	cfg Config
}

// NewJsonFormatter creates a new jsonFormatter using the provided configuration.
func NewJsonFormatter(cfg Config) Formatter {
	return &jsonFormatter{cfg: cfg}
}

// Format renders the log entry as a JSON string.
// It includes the log level, timestamp, error (if any), fields, and message.
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
