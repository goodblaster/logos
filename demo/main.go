package main

import (
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
