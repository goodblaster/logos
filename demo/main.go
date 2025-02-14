package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/goodblaster/logos"
)

func main() {
	// Just use it, as console logger, with no setup.
	logos.Debug("This is a debug message")

	// Include some key/values pairs.
	logos.With("key2", "value2").Info("This is an %s message", "info")

	// Break off into sub-loggers.
	sublog1 := logos.With("sublog", 1)
	sublog1.Info("This is an %s message with %s", "info", "sublog1")

	sublog2 := sublog1.With("sublog", 2)
	sublog2.Info("This is an %s message with %s", "info", "sublog2")

	sublog1.Warn("This is a %s message with %s", "warn", "sublog1")

	// Messages for all levels.
	for defLevel := logos.LevelDebug; defLevel < logos.LevelFatal; defLevel++ {
		logos.SetLevel(defLevel)
		logos.Print("")
		logos.Print("DEFAULT LEVEL = %s -----", defLevel)
		for level := logos.LevelDebug; level < logos.LevelFatal; level++ {
			logos.Log(level, "This is a %s message", level)
		}
	}

	// Custom levels.
	const (
		LevelApple logos.Level = iota
		LevelBanana
		LevelCherry
	)

	// With custom names.
	logos.LevelNames = map[logos.Level]string{
		LevelApple:  "apple",
		LevelBanana: "banana",
		LevelCherry: "cherry",
	}

	// And custom console colors.
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
		log.Print("DEFAULT LEVEL = %s -----", defLevel)
		for level := LevelApple; level <= LevelCherry; level++ {
			log.Log(level, "This is a %s message", level)
		}
	}

	// Customer formatter
	custom := logos.NewLogger(LevelApple, customFormatter{}, os.Stdout)
	custom.With("key1", "value1").With("key2", "value2").Log(LevelCherry, "This is a custom message.")
	custom.With("key1", "value1").With("key2", "value2").Log(LevelCherry, "This is a custom message.")
}

type customFormatter struct{}

func (f customFormatter) Format(level logos.Level, msg string, fields map[string]any) string {
	var tuples []string
	for key, value := range fields {
		b, _ := json.Marshal(value)
		tuples = append(tuples, fmt.Sprintf("%s=%v", key, string(b)))
	}
	slices.Sort(tuples)

	// ANSI color codes
	textColor := logos.LevelColors[level]

	line := fmt.Sprintf("%s\t%s",
		strings.ToUpper(level.String()),
		time.Now().UTC().Format(time.ANSIC))

	for _, tuple := range tuples {
		line += fmt.Sprintf("\n\t%s", tuple)
	}

	lineMsg := fmt.Sprintf("%s%s%s%s", logos.ColorBgBlack, textColor, msg, logos.ColorReset)
	line += fmt.Sprintf("\n\t%s\n", lineMsg)
	return line
}
