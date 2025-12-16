package logos

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// textFormatter formats log entries as plain text without color codes.
type textFormatter struct {
	cfg Config
}

// NewTextFormatter creates a new textFormatter using the provided configuration.
func NewTextFormatter(cfg Config) Formatter {
	return &textFormatter{cfg: cfg}
}

// Format renders the log entry as a plain text string.
// It includes the timestamp, log level, optional fields, and message.
func (f textFormatter) Format(level Level, entry Entry) string {
	var tuples []string
	if entry.Error != nil {
		errMsg := entry.Error.Error()
		tuples = append(tuples, fmt.Sprintf("error=%q", string(errMsg)))
	}

	for key, value := range entry.Fields {
		b, err := json.Marshal(value)
		if err != nil {
			// If marshal fails, include an error indicator instead of silently failing
			tuples = append(tuples, fmt.Sprintf("%s=<marshal_error>", key))
		} else {
			tuples = append(tuples, fmt.Sprintf("%s=%v", key, string(b)))
		}
	}
	slices.Sort(tuples)

	// If there are tuples, add a tab to separate them from the message
	var tupleString string
	if len(tuples) > 0 {
		tupleString = strings.Join(tuples, " ") + "\t"
	}

	return fmt.Sprintf("%s\t%s\t%s%s", f.cfg.Timestamp(), GetLevelName(level, &f.cfg), tupleString, entry.Msg)
}
