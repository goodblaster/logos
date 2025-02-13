package logos

import (
	"math"
)

type Level int

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
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPrint = math.MaxInt
)

// LevelNames - change if you like.
var LevelNames map[Level]string

var LevelColors map[Level]Color

// init - set defaults
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
