package logos

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type textFormatter struct {
	cfg Config
}

func NewTextFormatter(cfg Config) Formatter {
	return &textFormatter{cfg: cfg}
}
func (f textFormatter) Format(level Level, msg string, fields map[string]any) string {
	var tuples []string
	for key, value := range fields {
		b, _ := json.Marshal(value)
		tuples = append(tuples, fmt.Sprintf("%s=%v", key, string(b)))
	}
	slices.Sort(tuples)

	return fmt.Sprintf("%s\t%s\t%s\t%s", f.cfg.Timestamp(), level.String(), strings.Join(tuples, " "), msg)
}
