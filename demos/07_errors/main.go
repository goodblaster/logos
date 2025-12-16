package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/goodblaster/logos"
)

// This demo shows error handling with WithError and error handlers.
func main() {
	log := logos.NewLogger(logos.LevelInfo, logos.JSONFormatter(), os.Stdout)

	logos.Print("Logging with errors:")
	logos.Print("=====================\n")

	// Attach an error to log entries
	err := errors.New("database connection failed")
	log.WithError(err).Error("Failed to connect to database")

	// Chain WithError with fields
	log.WithError(err).
		With("host", "db.example.com").
		With("port", 5432).
		Error("Connection attempt failed")

	// Error handler - called when write errors occur
	logos.Print("\n\nError handler example:")
	logos.Print("=======================")

	// Create a logger with an error handler
	errorHandlerLog := log.WithErrorHandler(func(writeErr error) {
		fmt.Fprintf(os.Stderr, "WRITE ERROR: %v\n", writeErr)
	})

	errorHandlerLog.Info("This should write successfully")

	// Simulate write error with a failing writer
	failingWriter := &failWriter{}
	failingLog := logos.NewLogger(logos.LevelInfo, logos.JSONFormatter(), failingWriter).
		WithErrorHandler(func(writeErr error) {
			fmt.Fprintf(os.Stderr, "ERROR HANDLER CALLED: %v\n", writeErr)
		})

	failingLog.Info("This will trigger the error handler")
}

// failWriter always returns an error on Write
type failWriter struct{}

func (w *failWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("simulated write failure")
}
