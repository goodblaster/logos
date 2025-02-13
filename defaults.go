package logos

import "os"

func init() {
	DefaultLogger = NewLogger(LevelDebug, ConsoleFormatter(), os.Stdout)
}

var DefaultLogger Logger

func SetDefaultLogger(logger Logger) {
	DefaultLogger = logger
}

func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

func With(key string, value any) Logger {
	return DefaultLogger.With(key, value)
}

func WithError(err error) Logger {
	return DefaultLogger.WithError(err)
}

func WithFields(fields map[string]any) Logger {
	return DefaultLogger.WithFields(fields)
}

func Log(level Level, format string, args ...any) {
	DefaultLogger.Log(level, format, args...)
}

func LogFunc(level Level, msg func() string) {
	DefaultLogger.LogFunc(level, msg)
}

func Print(format string, args ...any) {
	DefaultLogger.Print(format, args...)
}

func Debug(format string, args ...any) {
	DefaultLogger.Debug(format, args...)
}

func Info(format string, args ...any) {
	DefaultLogger.Info(format, args...)
}

func Warn(format string, args ...any) {
	DefaultLogger.Warn(format, args...)
}

func Error(format string, args ...any) {
	DefaultLogger.Error(format, args...)
}

func Fatal(format string, args ...any) {
	DefaultLogger.Fatal(format, args...)
}
