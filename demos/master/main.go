package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/goodblaster/logos"
)

// This master demo shows the most commonly used features of the logos library.
// It covers: basic logging, structured fields, log levels, formatters,
// context logging, error handling, and tee logging.
func main() {
	printSection("1. Basic Logging")
	basicLogging()

	printSection("2. Structured Logging with Fields")
	structuredLogging()

	printSection("3. Log Levels and Filtering")
	logLevels()

	printSection("4. Different Formatters")
	formatters()

	printSection("5. Context-Based Logging")
	contextLogging()

	printSection("6. Error Handling")
	errorHandling()

	printSection("7. Tee Logging (Multiple Destinations)")
	teeLogging()
}

func basicLogging() {
	// Use the default logger with package-level functions
	logos.Debug("Debug message - detailed information")
	logos.Info("Info message - general information")
	logos.Warn("Warning message - something unexpected")
	logos.Error("Error message - something went wrong")

	// Formatted messages
	logos.Infof("User %s performed action: %s", "alice", "login")
}

func structuredLogging() {
	// Add context with fields
	logos.With("user_id", "user-123").
		With("ip", "192.168.1.1").
		Info("User logged in")

	// Add multiple fields at once
	logos.WithFields(logos.Fields{
		"request_id": "req-456",
		"method":     "POST",
		"endpoint":   "/api/users",
		"status":     200,
	}).Info("Request processed")

	// Create a logger with persistent fields
	log := logos.With("component", "auth").
		With("version", "1.0.0")

	log.Info("Authentication module initialized")
	log.With("user", "bob").Info("User authenticated")
}

func logLevels() {
	// Create a logger at Info level
	log := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout)

	log.Debug("This won't show (below Info level)")
	log.Info("This will show")
	log.Warn("This will show")
	log.Error("This will show")

	// Change level dynamically
	log = log.WithLevel(logos.LevelWarn)
	fmt.Println("\n  (Changed level to Warn)")

	log.Info("This won't show (below Warn level)")
	log.Warn("This will show")
	log.Error("This will show")

	// Check if level is enabled before expensive operations
	if log.IsLevelEnabled(logos.LevelDebug) {
		// Do expensive debug formatting
		log.Debug("Expensive debug info")
	}
}

func formatters() {
	fmt.Println("\n  JSON Formatter:")
	jsonLog := logos.NewLogger(logos.LevelInfo, logos.JSONFormatter(), os.Stdout)
	jsonLog.With("format", "json").Info("Structured JSON output")

	fmt.Println("\n  Text Formatter:")
	textLog := logos.NewLogger(logos.LevelInfo, logos.TextFormatter(), os.Stdout)
	textLog.With("format", "text").Info("Simple text output")

	fmt.Println("\n  Console Formatter (colored):")
	consoleLog := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout)
	consoleLog.Info("Colored console output")
	consoleLog.Warn("Warning in yellow")
	consoleLog.Error("Error in red")
}

func contextLogging() {
	// Create a request-specific logger
	requestLog := logos.NewLogger(logos.LevelInfo, logos.JSONFormatter(), os.Stdout).
		With("request_id", "req-789").
		With("user_id", "user-999")

	// Store in context
	ctx := logos.WithLogger(context.Background(), requestLog)

	// Pass context through the call stack
	handleAPIRequest(ctx)
}

func handleAPIRequest(ctx context.Context) {
	// Retrieve logger from context
	log := logos.FromContext(ctx)

	log.Info("Processing API request")
	log.With("endpoint", "/api/orders").Info("Validating request")
	log.With("order_id", "order-555").Info("Order created")
}

func errorHandling() {
	// Attach errors to log entries
	err := errors.New("database connection timeout")
	logos.WithError(err).Error("Failed to connect to database")

	// Combine errors with fields
	logos.WithError(err).
		With("host", "db.example.com").
		With("port", 5432).
		Error("Connection failed")

	// Error handler for write failures
	log := logos.NewLogger(logos.LevelInfo, logos.JSONFormatter(), os.Stdout).
		WithErrorHandler(func(writeErr error) {
			fmt.Fprintf(os.Stderr, "[WRITE ERROR] %v\n", writeErr)
		})

	log.Info("Logs with error handler attached")
}

func teeLogging() {
	// Create loggers for different purposes
	errorLog := logos.NewLogger(logos.LevelError, logos.JSONFormatter(), os.Stdout)
	auditLog := logos.NewLogger(logos.LevelInfo, logos.TextFormatter(), os.Stdout)

	// Main logger with tee destinations
	log := logos.NewLogger(logos.LevelDebug, logos.ConsoleFormatter(), os.Stdout).
		Tee(errorLog, auditLog)

	log.Debug("Debug - only in main log")
	log.Info("Info - in main and audit logs")
	log.Error("Error - in all three logs")

	fmt.Println("\n  Each destination can have its own level and formatter!")
}

func printSection(title string) {
	fmt.Println()
	fmt.Println("=" + strings(len(title)+2, "="))
	fmt.Printf(" %s \n", title)
	fmt.Println("=" + strings(len(title)+2, "="))
	fmt.Println()
}

func strings(n int, s string) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
