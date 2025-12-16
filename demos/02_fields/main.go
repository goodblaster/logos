package main

import (
	"github.com/goodblaster/logos"
)

// This demo shows structured logging with fields.
func main() {
	// Add a single field
	logos.With("user_id", "user-123").Info("User logged in")

	// Add multiple fields
	logos.WithFields(logos.Fields{
		"user_id":    "user-456",
		"request_id": "req-789",
		"ip":         "10.0.1.5",
	}).Info("Request processed")

	// Chain With() calls to build up context
	log := logos.With("component", "auth").
		With("version", "1.2.3")

	log.Info("Authentication started")
	log.Warn("Invalid credentials")

	// Create sub-loggers with additional context
	userLog := log.With("user_id", "user-999")
	userLog.Info("Processing user request")
	userLog.With("action", "update_profile").Info("Profile updated")
}
