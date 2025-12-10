package logos

import (
	"io"
	"os"
	"strings"
)

// DefaultLogger is the global logger used by package-level log functions.
var DefaultLogger Logger

// init sets the default logger to output debug-level logs to the console.
func init() {
	level := LevelDebug
	if logLevel := strings.ToLower(os.Getenv("LOG_LEVEL")); logLevel != "" {
		if lvl, ok := DefaultLevels[logLevel]; ok {
			level = lvl
		}
	}

	formatter := ConsoleFormatter()
	switch strings.ToLower(os.Getenv("LOG_FORMAT")) {
	case "json":
		formatter = JSONFormatter()
	case "text":
		formatter = TextFormatter()
	case "console":
		formatter = ConsoleFormatter()
	}

	DefaultLogger = NewLogger(level, formatter, os.Stdout)
}

// SetDefaultLogger overrides the global DefaultLogger with a new one.
func SetDefaultLogger(logger Logger) {
	DefaultLogger = logger
}

// SetLevel sets the logging level of the DefaultLogger.
func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

// With returns a copy of the DefaultLogger with an additional field.
func With(key string, value any) Logger {
	return DefaultLogger.With(key, value)
}

// WithError returns a copy of the DefaultLogger with an associated error.
func WithError(err error) Logger {
	return DefaultLogger.WithError(err)
}

// WithFields returns a copy of the DefaultLogger with additional fields.
func WithFields(fields map[string]any) Logger {
	return DefaultLogger.WithFields(fields)
}

// WithTee returns a copy of the DefaultLogger that writes to both the existing writer and the provided writers.
func WithTee(writers ...io.Writer) Logger {
	return DefaultLogger.WithTee(writers...)
}

// Log logs a message at the specified level using the DefaultLogger.
func Log(level Level, a ...any) {
	DefaultLogger.Log(level, a...)
}

// Logf logs a formatted message at the specified level using the DefaultLogger.
func Logf(level Level, format string, args ...any) {
	DefaultLogger.Logf(level, format, args...)
}

// LogFunc logs a lazily-evaluated message using the DefaultLogger if the level is enabled.
func LogFunc(level Level, msg func() string) {
	DefaultLogger.LogFunc(level, msg)
}

// LogIf executes a function if the level is enabled using the DefaultLogger.
func LogIf(level Level, log func()) {
	DefaultLogger.LogIf(level, log)
}

// Print logs a message at the print level using the DefaultLogger.
func Print(a ...any) {
	DefaultLogger.Print(a...)
}

// Printf logs a formatted message at the print level using the DefaultLogger.
func Printf(format string, args ...any) {
	DefaultLogger.Printf(format, args...)
}

// Debug logs a message at the debug level using the DefaultLogger.
func Debug(s ...any) {
	DefaultLogger.Debug(s...)
}

// Debugf logs a formatted message at the debug level using the DefaultLogger.
func Debugf(format string, args ...any) {
	DefaultLogger.Debugf(format, args...)
}

// Info logs a message at the info level using the DefaultLogger.
func Info(a ...any) {
	DefaultLogger.Info(a...)
}

// Infof logs a formatted message at the info level using the DefaultLogger.
func Infof(format string, args ...any) {
	DefaultLogger.Infof(format, args...)
}

// Warn logs a message at the warn level using the DefaultLogger.
func Warn(a ...any) {
	DefaultLogger.Warn(a...)
}

// Warnf logs a formatted message at the warn level using the DefaultLogger.
func Warnf(format string, args ...any) {
	DefaultLogger.Warnf(format, args...)
}

// Error logs a message at the error level using the DefaultLogger.
func Error(a ...any) {
	DefaultLogger.Error(a...)
}

// Errorf logs a formatted message at the error level using the DefaultLogger.
func Errorf(format string, args ...any) {
	DefaultLogger.Errorf(format, args...)
}

// Fatal logs a message at the fatal level using the DefaultLogger and panics.
func Fatal(a ...any) {
	DefaultLogger.Fatal(a...)
}

// Fatalf logs a formatted message at the fatal level using the DefaultLogger and panics.
func Fatalf(format string, args ...any) {
	DefaultLogger.Fatalf(format, args...)
}
