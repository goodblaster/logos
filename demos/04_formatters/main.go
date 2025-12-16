package main

import (
	"os"

	"github.com/goodblaster/logos"
)

// This demo shows different formatters: JSON, Text, and Console.
func main() {
	logos.Print("JSON Formatter:")
	logos.Print("================")
	jsonLog := logos.NewLogger(logos.LevelInfo, logos.JSONFormatter(), os.Stdout)
	jsonLog.With("user", "alice").With("action", "login").Info("User logged in")

	logos.Print("\nText Formatter:")
	logos.Print("================")
	textLog := logos.NewLogger(logos.LevelInfo, logos.TextFormatter(), os.Stdout)
	textLog.With("user", "bob").With("action", "logout").Info("User logged out")

	logos.Print("\nConsole Formatter (with colors):")
	logos.Print("==================================")
	consoleLog := logos.NewLogger(logos.LevelDebug, logos.ConsoleFormatter(), os.Stdout)
	consoleLog.Debug("This is a debug message")
	consoleLog.Info("This is an info message")
	consoleLog.Warn("This is a warning message")
	consoleLog.Error("This is an error message")
}
