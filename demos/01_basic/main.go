package main

import (
	"github.com/goodblaster/logos"
)

// This demo shows basic logging operations using the default logger.
func main() {
	// Use package-level functions with the default logger
	logos.Debug("This is a debug message")
	logos.Info("This is an info message")
	logos.Warn("This is a warning message")
	logos.Error("This is an error message")

	// Formatted logging
	logos.Infof("User %s logged in from %s", "alice", "192.168.1.1")
	logos.Debugf("Processing request #%d", 12345)

	// Print always shows regardless of level
	logos.Print("This message always appears")
}
