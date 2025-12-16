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
// If marshaling fails, it returns an error message in JSON format instead of panicking.
func (f jsonFormatter) Format(level Level, entry Entry) string {
	type JsonStruct struct {
		Level     string         `json:"level"`
		Timestamp string         `json:"timestamp"`
		Error     error          `json:"error,omitempty"`
		Fields    map[string]any `json:"fields,omitempty"`
		Msg       string         `json:"msg"`
	}

	jsonStruct := JsonStruct{
		Level:     GetLevelName(level, &f.cfg),
		Timestamp: f.cfg.Timestamp(),
		Error:     entry.Error,
		Msg:       entry.Msg,
		Fields:    entry.Fields,
	}

	b, err := json.Marshal(jsonStruct)
	if err != nil {
		// Instead of panicking, return an error message in JSON format
		errorMsg := JsonStruct{
			Level:     "error",
			Timestamp: f.cfg.Timestamp(),
			Error:     errors.Wrap(err, "failed to marshal log entry"),
			Msg:       "[LOG ERROR: failed to marshal entry]",
			Fields:    nil, // Don't include fields that might have caused the error
		}
		if errorBytes, innerErr := json.Marshal(errorMsg); innerErr == nil {
			return string(errorBytes)
		}
		// If even the error message can't be marshaled, return a simple string
		return `{"level":"error","msg":"[LOG ERROR: catastrophic marshal failure]"}`
	}

	return string(b)
}
