package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/goodblaster/errors"
	"github.com/goodblaster/logos"
)

func main() {
	// Basic usage with default console logger.
	logos.Debug("This is a debug message")

	// Include structured key-value fields.
	logos.With("key2", "value2").Infof("This is an %s message", "info")

	// Attach an error.
	logos.WithError(errors.New("this is an error")).Error("this is a message")

	// Create sub-loggers for scoped context.
	sublog1 := logos.With("sublog", 1)
	sublog1.Infof("This is an %s message with %s", "info", "sublog1")

	sublog2 := sublog1.With("sublog", 2)
	sublog2.Infof("This is an %s message with %s", "info", "sublog2")
	sublog1.Warnf("This is a %s message with %s", "warn", "sublog1")

	// Demonstrate all default levels.
	for defLevel := logos.LevelDebug; defLevel < logos.LevelFatal; defLevel++ {
		logos.SetLevel(defLevel)
		logos.Print("")
		logos.Printf("DEFAULT LEVEL = %s -----", defLevel)
		for level := logos.LevelDebug; level < logos.LevelFatal; level++ {
			logos.Logf(level, "This is a %s message", level)
		}
	}

	// LogFunc (lazy evaluation)
	logos.LogFunc(logos.LevelDebug, func() string {
		return expensiveOperation()
	})

	// LogIf (conditional execution)
	logos.LogIf(logos.LevelInfo, func() {
		fmt.Println("This log is only printed if info level is enabled")
	})

	// Custom level demonstration.
	const (
		LevelApple logos.Level = iota
		LevelBanana
		LevelCherry
	)

	logos.LevelNames = map[logos.Level]string{
		LevelApple:  "apple",
		LevelBanana: "banana",
		LevelCherry: "cherry",
	}

	logos.LevelColors = map[logos.Level]logos.Color{
		LevelApple:  logos.ColorBgGreen + logos.ColorTextBlack,
		LevelBanana: logos.ColorBgYellow + logos.ColorTextBlack,
		LevelCherry: logos.ColorBgRed + logos.ColorTextBlack,
	}

	logos.Print("")
	logos.Print("CUSTOM LEVELS ----------")
	log := logos.NewLogger(LevelApple, logos.ConsoleFormatter(), os.Stdout)
	for defLevel := LevelApple; defLevel <= LevelCherry; defLevel++ {
		log.SetLevel(defLevel)
		log.Print("")
		log.Printf("DEFAULT LEVEL = %s -----", defLevel)
		for level := LevelApple; level <= LevelCherry; level++ {
			log.Logf(level, "This is a %s message", level)
		}
	}

	// Custom formatter demo.
	custom := logos.NewLogger(LevelApple, customFormatter{}, os.Stdout)
	custom.With("key1", "value1").With("key2", "value2").Log(LevelCherry, "This is a custom message.")

	// Context logging demo.
	logos.Print("")
	logos.Print("CONTEXT LOGGING ----------")

	// Create a logger with request-specific fields
	requestLogger := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout).
		With("request_id", "req-123").
		With("user_id", "user-456").
		With("ip", "192.168.1.1")

	// Store logger in context for request lifecycle
	ctx := logos.WithLogger(context.Background(), requestLogger)

	// Simulate request handling - logger is retrieved from context
	handleRequest(ctx)
	processPayment(ctx)

	// If no logger in context, returns DefaultLogger
	emptyCtx := context.Background()
	logos.FromContext(emptyCtx).Info("Using default logger (no context logger)")

	// Tee logging demo - write to multiple destinations simultaneously.
	logos.Print("")
	logos.Print("TEE LOGGING ----------")

	// Create separate loggers for different destinations with different levels
	consoleBuf := &bytes.Buffer{}
	fileBuf := &bytes.Buffer{}

	consoleLogger := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), consoleBuf)
	fileLogger := logos.NewLogger(logos.LevelDebug, logos.JSONFormatter(), fileBuf)

	// Create a main logger and tee the other loggers to it
	mainLogger := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout)
	teeLogger := mainLogger.Tee(consoleLogger, fileLogger)

	teeLogger.Info("This message goes to stdout, consoleBuf (Info), and fileBuf (Debug)")
	teeLogger.With("component", "auth").Warn("Warning also goes to all destinations")

	// Debug message: only goes to fileBuf (Debug level), not stdout or consoleBuf (Info level)
	teeLogger.Debug("This debug message only goes to fileBuf")

	// Verify tee logging worked
	fmt.Printf("\nConsole buffer received: %s\n", consoleBuf.String())
	fmt.Printf("File buffer received: %s\n", fileBuf.String())

	// Package-level tee logging
	logos.Print("")
	logos.Print("PACKAGE-LEVEL TEE LOGGING ----------")
	teeBuf := &bytes.Buffer{}
	teeBufLogger := logos.NewLogger(logos.LevelDebug, logos.JSONFormatter(), teeBuf)
	teeLog := logos.Tee(teeBufLogger)
	teeLog.Info("This goes to both DefaultLogger's writer and teeBuf")
}

// handleRequest simulates a request handler that uses logger from context
func handleRequest(ctx context.Context) {
	logger := logos.FromContext(ctx)
	logger.Info("Processing request")
	logger.With("action", "validate").Info("Validating input")
	logger.With("action", "authorize").Info("Checking permissions")
}

// processPayment simulates another function that uses logger from context
func processPayment(ctx context.Context) {
	logger := logos.FromContext(ctx)
	logger.With("action", "payment").Info("Processing payment")
	logger.With("amount", 99.99).Info("Payment successful")
}

func expensiveOperation() string {
	time.Sleep(50 * time.Millisecond)
	return "result of expensive operation"
}

type customFormatter struct{}

func (f customFormatter) Format(level logos.Level, entry logos.Entry) string {
	var tuples []string
	for key, value := range entry.Fields {
		b, _ := json.Marshal(value)
		tuples = append(tuples, fmt.Sprintf("%s=%v", key, string(b)))
	}
	slices.Sort(tuples)

	textColor := logos.LevelColors[level]

	line := fmt.Sprintf("%s\t%s",
		strings.ToUpper(level.String()),
		time.Now().UTC().Format(time.ANSIC))

	for _, tuple := range tuples {
		line += fmt.Sprintf("\n\t%s", tuple)
	}

	lineMsg := fmt.Sprintf("%s%s%s%s", logos.ColorBgBlack, textColor, entry.Msg, logos.ColorReset)
	line += fmt.Sprintf("\n\t%s\n", lineMsg)
	return line
}
