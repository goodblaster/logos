package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goodblaster/logos"
)

// This demo shows tee logging - writing to multiple destinations simultaneously.
func main() {
	// Create buffers to simulate different destinations
	errorBuf := &bytes.Buffer{}
	auditBuf := &bytes.Buffer{}

	// Create loggers for each destination with different levels and formatters
	errorLogger := logos.NewLogger(logos.LevelError, logos.JSONFormatter(), errorBuf)
	auditLogger := logos.NewLogger(logos.LevelInfo, logos.TextFormatter(), auditBuf)

	// Create main logger that outputs to stdout
	mainLogger := logos.NewLogger(logos.LevelDebug, logos.ConsoleFormatter(), os.Stdout)

	// Tee the other loggers to the main logger
	log := mainLogger.Tee(errorLogger, auditLogger)

	logos.Print("Logging to multiple destinations:")
	logos.Print("===================================\n")

	// Debug message - only goes to mainLogger (stdout)
	log.Debug("Debug message - only in stdout")

	// Info message - goes to mainLogger (stdout) and auditLogger
	log.Info("Info message - in stdout and audit log")

	// Warning message - goes to mainLogger (stdout) and auditLogger
	log.Warn("Warning message - in stdout and audit log")

	// Error message - goes to all three loggers
	log.Error("Error message - in stdout, audit, and error logs")

	// Show what was captured
	logos.Print("\n\nError log buffer (JSON, Error level only):")
	logos.Print("============================================")
	fmt.Print(errorBuf.String())

	logos.Print("\nAudit log buffer (Text, Info+ levels):")
	logos.Print("========================================")
	fmt.Print(auditBuf.String())
}
