package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goodblaster/logos"
)

// This demo shows how to create custom formatters.
func main() {
	// Use a simple custom formatter
	simpleLog := logos.NewLogger(logos.LevelInfo, SimpleFormatter{}, os.Stdout)

	logos.Print("Simple custom formatter:")
	logos.Print("=========================")
	simpleLog.Info("Basic message")
	simpleLog.With("user", "alice").Info("Message with field")

	// Use a more complex key-value formatter
	logos.Print("\n\nKey-Value formatter:")
	logos.Print("=====================")
	kvLog := logos.NewLogger(logos.LevelInfo, KeyValueFormatter{}, os.Stdout)
	kvLog.With("request_id", "req-123").
		With("user_id", "user-456").
		With("duration_ms", 42).
		Info("Request completed")
}

// SimpleFormatter is a minimal formatter that outputs: [LEVEL] message
type SimpleFormatter struct{}

func (f SimpleFormatter) Format(level logos.Level, entry logos.Entry) string {
	levelStr := strings.ToUpper(level.String())
	if len(entry.Fields) > 0 {
		fieldsJSON, _ := json.Marshal(entry.Fields)
		return fmt.Sprintf("[%s] %s | %s", levelStr, entry.Msg, string(fieldsJSON))
	}
	return fmt.Sprintf("[%s] %s", levelStr, entry.Msg)
}

// KeyValueFormatter outputs logs in a key=value format
type KeyValueFormatter struct{}

func (f KeyValueFormatter) Format(level logos.Level, entry logos.Entry) string {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	parts := []string{
		fmt.Sprintf("time=%s", timestamp),
		fmt.Sprintf("level=%s", level.String()),
		fmt.Sprintf("msg=%q", entry.Msg),
	}

	// Add fields
	for key, value := range entry.Fields {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}

	// Add error if present
	if entry.Error != nil {
		parts = append(parts, fmt.Sprintf("error=%q", entry.Error.Error()))
	}

	return strings.Join(parts, " ")
}
