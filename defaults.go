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

func Log(level Level, a ...any) {
	DefaultLogger.Log(level, a...)
}

func Logf(level Level, format string, args ...any) {
	DefaultLogger.Logf(level, format, args...)
}

func LogFunc(level Level, msg func() string) {
	DefaultLogger.LogFunc(level, msg)
}

func LogIf(level Level, log func()) {
	DefaultLogger.LogIf(level, log)
}

func Print(a ...any) {
	DefaultLogger.Print(a...)
}

func Printf(format string, args ...any) {
	DefaultLogger.Printf(format, args...)
}

func Debug(s ...any) {
	DefaultLogger.Debug(s...)
}

func Debugf(format string, args ...any) {
	DefaultLogger.Debugf(format, args...)
}

func Info(a ...any) {
	DefaultLogger.Info(a...)
}

func Infof(format string, args ...any) {
	DefaultLogger.Infof(format, args...)
}

func Warn(a ...any) {
	DefaultLogger.Warn(a...)
}

func Warnf(format string, args ...any) {
	DefaultLogger.Warnf(format, args...)
}

func Error(a ...any) {
	DefaultLogger.Error(a...)
}

func Errorf(format string, args ...any) {
	DefaultLogger.Errorf(format, args...)
}

func Fatal(a ...any) {
	DefaultLogger.Fatal(a...)
}

func Fatalf(format string, args ...any) {
	DefaultLogger.Fatalf(format, args...)
}
