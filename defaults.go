package logos

import (
	"os"
	"strings"
	"sync"
)

// defaultLogger is the global logger used by package-level log functions.
var defaultLogger Logger
var defaultLoggerMu sync.RWMutex

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

	defaultLogger = NewLogger(level, formatter, os.Stdout)
}

// SetDefaultLogger overrides the global default logger with a new one.
// This function is thread-safe.
func SetDefaultLogger(logger Logger) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger = logger
}

// getDefaultLogger returns a copy of the default logger.
// This function is thread-safe.
func getDefaultLogger() Logger {
	defaultLoggerMu.RLock()
	defer defaultLoggerMu.RUnlock()
	return defaultLogger
}

// SetLevel sets the logging level of the default logger.
func SetLevel(level Level) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.SetLevel(level)
}

// IsLevelEnabled returns true if the default logger would log at the given level.
func IsLevelEnabled(level Level) bool {
	return getDefaultLogger().IsLevelEnabled(level)
}

// With returns a copy of the default logger with an additional field.
func With(key string, value any) Logger {
	return getDefaultLogger().With(key, value)
}

// WithError returns a copy of the default logger with an associated error.
func WithError(err error) Logger {
	return getDefaultLogger().WithError(err)
}

// WithFields returns a copy of the default logger with additional fields.
func WithFields(fields map[string]any) Logger {
	return getDefaultLogger().WithFields(fields)
}

// Tee adds one or more loggers as tee destinations to the default logger.
func Tee(loggers ...Logger) Logger {
	return getDefaultLogger().Tee(loggers...)
}

// Log logs a message at the specified level using the default logger.
func Log(level Level, a ...any) {
	getDefaultLogger().Log(level, a...)
}

// Logf logs a formatted message at the specified level using the default logger.
func Logf(level Level, format string, args ...any) {
	getDefaultLogger().Logf(level, format, args...)
}

// LogFunc logs a lazily-evaluated message using the default logger if the level is enabled.
func LogFunc(level Level, msg func() string) {
	getDefaultLogger().LogFunc(level, msg)
}

// LogIf executes a function if the level is enabled using the default logger.
func LogIf(level Level, log func()) {
	getDefaultLogger().LogIf(level, log)
}

// Print logs a message at the print level using the default logger.
func Print(a ...any) {
	getDefaultLogger().Print(a...)
}

// Printf logs a formatted message at the print level using the default logger.
func Printf(format string, args ...any) {
	getDefaultLogger().Printf(format, args...)
}

// Debug logs a message at the debug level using the default logger.
func Debug(s ...any) {
	getDefaultLogger().Debug(s...)
}

// Debugf logs a formatted message at the debug level using the default logger.
func Debugf(format string, args ...any) {
	getDefaultLogger().Debugf(format, args...)
}

// Info logs a message at the info level using the default logger.
func Info(a ...any) {
	getDefaultLogger().Info(a...)
}

// Infof logs a formatted message at the info level using the default logger.
func Infof(format string, args ...any) {
	getDefaultLogger().Infof(format, args...)
}

// Warn logs a message at the warn level using the default logger.
func Warn(a ...any) {
	getDefaultLogger().Warn(a...)
}

// Warnf logs a formatted message at the warn level using the default logger.
func Warnf(format string, args ...any) {
	getDefaultLogger().Warnf(format, args...)
}

// Error logs a message at the error level using the default logger.
func Error(a ...any) {
	getDefaultLogger().Error(a...)
}

// Errorf logs a formatted message at the error level using the default logger.
func Errorf(format string, args ...any) {
	getDefaultLogger().Errorf(format, args...)
}

// Fatal logs a message at the fatal level using the default logger and panics.
func Fatal(a ...any) {
	getDefaultLogger().Fatal(a...)
}

// Fatalf logs a formatted message at the fatal level using the default logger and panics.
func Fatalf(format string, args ...any) {
	getDefaultLogger().Fatalf(format, args...)
}
