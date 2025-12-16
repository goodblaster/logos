package main

import (
	"os"

	"github.com/goodblaster/logos"
)

// This demo shows how to create custom log levels.
func main() {
	// Define custom levels
	const (
		LevelTrace   logos.Level = iota - 2 // Lower number = more verbose
		LevelVerbose
		LevelNotice  logos.Level = iota + 2 // Between Warn and Error
	)

	// Register custom level names using thread-safe setter
	logos.SetLevelName(LevelTrace, "trace")
	logos.SetLevelName(LevelVerbose, "verbose")
	logos.SetLevelName(LevelNotice, "notice")

	// Register custom colors (optional, for console formatter)
	logos.SetLevelColor(LevelTrace, logos.ColorTextCyan)
	logos.SetLevelColor(LevelVerbose, logos.ColorTextBlue)
	logos.SetLevelColor(LevelNotice, logos.ColorTextMagenta)

	// Cleanup after demo
	defer func() {
		logos.SetLevelName(LevelTrace, "")
		logos.SetLevelName(LevelVerbose, "")
		logos.SetLevelName(LevelNotice, "")
	}()

	log := logos.NewLogger(LevelTrace, logos.ConsoleFormatter(), os.Stdout)

	logos.Print("Custom log levels:")
	logos.Print("===================\n")

	// Use custom levels
	log.Log(LevelTrace, "This is a trace message (most verbose)")
	log.Log(LevelVerbose, "This is a verbose message")
	log.Log(logos.LevelDebug, "This is a debug message")
	log.Log(logos.LevelInfo, "This is an info message")
	log.Log(logos.LevelWarn, "This is a warning message")
	log.Log(LevelNotice, "This is a notice message (between warn and error)")
	log.Log(logos.LevelError, "This is an error message")

	// Change level to filter out verbose messages
	logos.Print("\n\nWith level set to Debug (Trace and Verbose filtered):")
	logos.Print("========================================================")
	log = log.WithLevel(logos.LevelDebug)

	log.Log(LevelTrace, "This won't show")
	log.Log(LevelVerbose, "This won't show")
	log.Log(logos.LevelDebug, "This is a debug message")
	log.Log(logos.LevelInfo, "This is an info message")
	log.Log(LevelNotice, "This is a notice message")
}
