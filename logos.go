package logos

import (
	"fmt"
	"io"
	"sync"
)

// Fields represents a key-value pair used to annotate log entries.
type Fields = map[string]any

// Logger is the primary struct for logging messages with optional fields and errors.
type Logger struct {
	level     *Level
	formatter Formatter
	writer    io.Writer
	sync      *sync.Mutex
	fields    Fields
	error     error
}

// NewLogger creates a new Logger instance with the given level, formatter, and output writer.
func NewLogger(level Level, formatter Formatter, writer io.Writer) Logger {
	return Logger{
		level:     &level,
		formatter: formatter,
		writer:    writer,
		sync:      &sync.Mutex{},
		fields:    nil,
	}
}

// SetLevel sets the logging level of the logger.
func (logger Logger) SetLevel(level Level) {
	*logger.level = level
}

// Copy creates a deep copy of the logger, duplicating any fields.
func (logger Logger) Copy() Logger {
	logger.sync.Lock()
	defer logger.sync.Unlock()

	newLogger := logger

	if newLogger.fields != nil {
		newLogger.fields = make(Fields)
		for key, value := range logger.fields {
			newLogger.fields[key] = value
		}
	}

	return newLogger
}

// With returns a new Logger with an added single key-value field.
func (logger Logger) With(key string, value any) Logger {
	return logger.WithFields(Fields{key: value})
}

// WithFields returns a new Logger with additional key-value pairs.
func (logger Logger) WithFields(fields Fields) Logger {
	newLogger := logger.Copy()

	if newLogger.fields == nil {
		newLogger.fields = make(Fields)
	}

	for key, value := range fields {
		newLogger.fields[key] = value
	}

	return newLogger
}

// WithError returns a new Logger with an associated error.
func (logger Logger) WithError(err error) Logger {
	newLogger := logger.Copy()
	if newLogger.error != nil {
		logger.With("old_error", newLogger.error).
			With("new_error", err).
			Error("overwriting old error with new error")
	}

	newLogger.error = err
	return newLogger
}

// Log logs a message at the specified level.
func (logger Logger) Log(level Level, a ...any) {
	if *logger.level > level {
		return
	}
	msg := fmt.Sprint(a...)
	line := logger.formatter.Format(level, Entry{
		Fields: logger.fields,
		Msg:    msg,
		Error:  logger.error,
	})
	_, _ = fmt.Fprintln(logger.writer, line)
}

// Logf logs a formatted message at the specified level.
func (logger Logger) Logf(level Level, format string, args ...any) {
	if *logger.level > level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	line := logger.formatter.Format(level, Entry{
		Fields: logger.fields,
		Msg:    msg,
		Error:  logger.error,
	})
	_, _ = fmt.Fprintln(logger.writer, line)
}

// LogFunc evaluates the message-producing function only if the log level is enabled.
func (logger Logger) LogFunc(level Level, msg func() string) {
	if *logger.level > level {
		return
	}
	logger.Log(level, msg())
}

// LogIf calls the provided function if the log level is enabled.
func (logger Logger) LogIf(level Level, log func()) {
	if *logger.level > level {
		return
	}
	log()
}

// Print logs a message at the print level.
func (logger Logger) Print(a ...any) {
	logger.Log(LevelPrint, a...)
}

// Printf logs a formatted message at the print level.
func (logger Logger) Printf(format string, args ...any) {
	logger.Logf(LevelPrint, format, args...)
}

// Debug logs a message at the debug level.
func (logger Logger) Debug(a ...any) {
	logger.Log(LevelDebug, a...)
}

// Debugf logs a formatted message at the debug level.
func (logger Logger) Debugf(format string, args ...any) {
	logger.Logf(LevelDebug, format, args...)
}

// Info logs a message at the info level.
func (logger Logger) Info(a ...any) {
	logger.Log(LevelInfo, a...)
}

// Infof logs a formatted message at the info level.
func (logger Logger) Infof(format string, args ...any) {
	logger.Logf(LevelInfo, format, args...)
}

// Warn logs a message at the warn level.
func (logger Logger) Warn(a ...any) {
	logger.Log(LevelWarn, a...)
}

// Warnf logs a formatted message at the warn level.
func (logger Logger) Warnf(format string, args ...any) {
	logger.Logf(LevelWarn, format, args...)
}

// Error logs a message at the error level.
func (logger Logger) Error(a ...any) {
	logger.Log(LevelError, a...)
}

// Errorf logs a formatted message at the error level.
func (logger Logger) Errorf(format string, args ...any) {
	logger.Logf(LevelError, format, args...)
}

// Fatal logs a message at the fatal level and then panics.
func (logger Logger) Fatal(a ...any) {
	logger.Log(LevelFatal, a...)
	panic(fmt.Sprint(a...))
}

// Fatalf logs a formatted message at the fatal level and then panics.
func (logger Logger) Fatalf(format string, args ...any) {
	logger.Logf(LevelFatal, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Entry holds the log data including fields, message, and error.
type Entry struct {
	Fields Fields
	Msg    string
	Error  error
}
