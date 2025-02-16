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
func (f textFormatter) Format(level Level, entry Entry) string {
	var tuples []string
	if entry.Error != nil {
		errMsg := entry.Error.Error()
		tuples = append(tuples, fmt.Sprintf("error=%q", string(errMsg)))
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

	return fmt.Sprintf("%s\t%s\t%s%s", f.cfg.Timestamp(), level.String(), tupleString, entry.Msg)
}
