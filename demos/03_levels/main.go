package main

import (
	"os"

	"github.com/goodblaster/logos"
)

// This demo shows log levels and filtering.
func main() {
	log := logos.NewLogger(logos.LevelDebug, logos.ConsoleFormatter(), os.Stdout)

	logos.Print("Level: Debug - All messages shown")
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")

	logos.Print("\nLevel: Info - Debug filtered out")
	log = log.WithLevel(logos.LevelInfo)
	log.Debug("This won't be shown")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")

	logos.Print("\nLevel: Warn - Info and Debug filtered out")
	log = log.WithLevel(logos.LevelWarn)
	log.Debug("This won't be shown")
	log.Info("This won't be shown")
	log.Warn("This is a warning message")
	log.Error("This is an error message")

	logos.Print("\nLevel: Error - Only errors shown")
	log = log.WithLevel(logos.LevelError)
	log.Debug("This won't be shown")
	log.Info("This won't be shown")
	log.Warn("This won't be shown")
	log.Error("This is an error message")

	// Check if a level is enabled before expensive operations
	logos.Print("\nUsing IsLevelEnabled:")
	if log.IsLevelEnabled(logos.LevelDebug) {
		logos.Print("Debug is enabled")
	} else {
		logos.Print("Debug is NOT enabled")
	}

	if log.IsLevelEnabled(logos.LevelError) {
		logos.Print("Error is enabled")
	}
}
