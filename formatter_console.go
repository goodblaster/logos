package logos

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// consoleFormatter is a log formatter that adds ANSI color codes for terminal output.
type consoleFormatter struct {
	cfg Config
}

// NewConsoleFormatter creates a new consoleFormatter using the provided configuration.
func NewConsoleFormatter(cfg Config) Formatter {
	return &consoleFormatter{cfg: cfg}
}

// Format renders the log entry as a colored string using ANSI escape codes for terminal output.
func (f consoleFormatter) Format(level Level, entry Entry) string {
	// ANSI color codes
	textColor := LevelColors[level]

	var tuples []string
	if entry.Error != nil {
		errMsg := entry.Error.Error()
		tuples = append(tuples, fmt.Sprintf("error=%s%q%s", textColor, errMsg, ColorReset))
	}

	for key, value := range entry.Fields {
		b, _ := json.Marshal(value)
		tuples = append(tuples, fmt.Sprintf("%s=%v", key, string(b)))
	}
	slices.Sort(tuples)

	// If there are tuples, add a tab to separate them from the message
	var tupleString string
	if len(tuples) > 0 {
		tupleString = strings.Join(tuples, " ") + "\t"
	}

	return fmt.Sprintf("%s\t%s%s%s\t%s%s",
		f.cfg.Timestamp(),
		textColor, level.String(), ColorReset,
		tupleString,
		entry.Msg)
}
