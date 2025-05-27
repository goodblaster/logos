package logos

import (
	"math"
)

// Level represents the severity of a log message.
type Level int

// String returns the string representation of a logging level.
func (level Level) String() string {
	if name, ok := LevelNames[level]; ok {
		return name
	}

	// Keep "print" as a default option.
	if level == LevelPrint {
		return "print"
	}

	return "unknown"
}

const (
	// LevelDebug represents fine-grained debug information.
	LevelDebug Level = iota - 1
	// LevelInfo represents general operational entries about what's going on inside the application.
	LevelInfo
	// LevelWarn represents potentially harmful situations.
	LevelWarn
	// LevelError represents error events that might still allow the application to continue running.
	LevelError
	// LevelFatal represents very severe error events that will presumably lead the application to abort.
	LevelFatal
	// LevelPrint is used for messages that should always be printed regardless of level filtering.
	LevelPrint = math.MaxInt
)

// LevelNames maps Level values to human-readable strings. This can be overridden by the user.
var LevelNames map[Level]string

// LevelColors maps Level values to ANSI color codes for console output. This can be customized.
var LevelColors map[Level]Color

// init initializes default level names and associated colors.
func init() {
	LevelNames = map[Level]string{
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
		LevelFatal: "fatal",
		LevelPrint: "print",
	}

	LevelColors = map[Level]Color{
		LevelDebug: ColorTextBlue,
		LevelInfo:  ColorTextGreen,
		LevelWarn:  ColorTextYellow,
		LevelError: ColorTextRed,
		LevelFatal: ColorTextPurple,
		LevelPrint: ColorReset,
	}
}
