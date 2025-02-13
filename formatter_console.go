package logos

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// Just like the text formatter, but with colors
type consoleFormatter struct {
	cfg Config
}

func NewConsoleFormatter(cfg Config) Formatter {
	return &consoleFormatter{cfg: cfg}
}

func (f consoleFormatter) Format(level Level, msg string, fields map[string]any) string {
	var tuples []string
	for key, value := range fields {
		b, _ := json.Marshal(value)
		tuples = append(tuples, fmt.Sprintf("%s=%v", key, string(b)))
	}
	slices.Sort(tuples)

	// ANSI color codes
	textColor := LevelColors[level]

	return fmt.Sprintf("%s\t%s%s%s\t%s\t%s", f.cfg.Timestamp(), textColor, level.String(), ColorReset, strings.Join(tuples, " "), msg)
}
