package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goodblaster/logos"
)

// This demo shows lazy evaluation with LogFunc and LogIf.
func main() {
	log := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout)

	logos.Print("Lazy evaluation examples:")
	logos.Print("==========================\n")

	// LogFunc - evaluates function only if level is enabled
	logos.Print("LogFunc with Info level (function WILL be called):")
	log.LogFunc(logos.LevelInfo, func() string {
		logos.Print("  -> Expensive function executing...")
		time.Sleep(100 * time.Millisecond)
		return "Result of expensive operation"
	})

	logos.Print("\nLogFunc with Debug level (function will NOT be called):")
	log.LogFunc(logos.LevelDebug, func() string {
		logos.Print("  -> This should never print")
		return "This won't be evaluated"
	})

	// LogIf - executes function only if level is enabled
	logos.Print("\nLogIf with Info level (function WILL be called):")
	log.LogIf(logos.LevelInfo, func() {
		fmt.Println("  -> Custom logging logic executing...")
		log.Info("Message from LogIf")
	})

	logos.Print("\nLogIf with Debug level (function will NOT be called):")
	log.LogIf(logos.LevelDebug, func() {
		fmt.Println("  -> This should never print")
	})

	// Use case: avoid expensive formatting when debug is disabled
	logos.Print("\n\nPractical use case:")
	logos.Print("====================")

	data := map[string]interface{}{
		"users": []string{"alice", "bob", "charlie"},
		"count": 42,
	}

	// Bad: format always happens, even if debug is disabled
	// log.Debugf("Data: %+v", expensiveFormat(data))

	// Good: format only happens if debug is enabled
	log.LogFunc(logos.LevelDebug, func() string {
		return fmt.Sprintf("Data: %+v", expensiveFormat(data))
	})

	logos.Print("Debug was disabled, so expensiveFormat() was never called")
}

func expensiveFormat(data interface{}) interface{} {
	logos.Print("  -> expensiveFormat() called")
	time.Sleep(50 * time.Millisecond)
	return data
}
