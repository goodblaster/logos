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
	level      *Level
	formatter  Formatter
	writer     io.Writer
	sync       *sync.Mutex
	fields     Fields
	error      error
	teeLoggers []Logger
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
// Deprecated: Use WithLevel() instead for immutable logger creation.
func (logger Logger) SetLevel(level Level) {
	*logger.level = level
}

// GetLevel returns the current logging level.
func (logger Logger) GetLevel() Level {
	if logger.level == nil {
		return LevelInfo // Safe default
	}
	return *logger.level
}

// GetFields returns a copy of the logger's fields.
// The returned map is independent and modifications won't affect the logger.
func (logger Logger) GetFields() Fields {
	if logger.fields == nil {
		return nil
	}
	fieldsCopy := make(Fields, len(logger.fields))
	for k, v := range logger.fields {
		fieldsCopy[k] = v
	}
	return fieldsCopy
}

// GetError returns the logger's associated error.
func (logger Logger) GetError() error {
	return logger.error
}

// GetTeeCount returns the number of tee loggers attached to this logger.
func (logger Logger) GetTeeCount() int {
	return len(logger.teeLoggers)
}

// WithLevel returns a new Logger with the specified logging level.
// Unlike SetLevel(), this method returns a new logger instance with an independent level,
// maintaining immutability and preventing shared state issues.
func (logger Logger) WithLevel(level Level) Logger {
	newLogger := logger.Copy()
	newLevel := level
	newLogger.level = &newLevel // New pointer, not shared
	return newLogger
}

// Copy creates a deep copy of the logger, duplicating any fields and tee loggers.
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

	if len(logger.teeLoggers) > 0 {
		newLogger.teeLoggers = make([]Logger, len(logger.teeLoggers))
		for i, teeLogger := range logger.teeLoggers {
			newLogger.teeLoggers[i] = teeLogger.Copy()
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
	// Defensive nil checks
	if logger.level == nil || logger.formatter == nil || logger.writer == nil {
		return
	}

	msg := fmt.Sprint(a...)
	entry := Entry{
		Fields: logger.fields,
		Msg:    msg,
		Error:  logger.error,
	}

	// Write to main writer if level is enabled
	if *logger.level <= level {
		line := logger.formatter.Format(level, entry)
		// Lock to prevent concurrent writes to the same writer (e.g., bytes.Buffer)
		logger.sync.Lock()
		_, _ = fmt.Fprintln(logger.writer, line)
		logger.sync.Unlock()
	}

	// Always call Log on each tee logger (they handle their own level checking and formatting)
	for _, teeLogger := range logger.teeLoggers {
		teeLogger.Log(level, a...)
	}
}

// Logf logs a formatted message at the specified level.
func (logger Logger) Logf(level Level, format string, args ...any) {
	// Defensive nil checks
	if logger.level == nil || logger.formatter == nil || logger.writer == nil {
		return
	}

	msg := fmt.Sprintf(format, args...)
	entry := Entry{
		Fields: logger.fields,
		Msg:    msg,
		Error:  logger.error,
	}

	// Write to main writer if level is enabled
	if *logger.level <= level {
		line := logger.formatter.Format(level, entry)
		// Lock to prevent concurrent writes to the same writer (e.g., bytes.Buffer)
		logger.sync.Lock()
		_, _ = fmt.Fprintln(logger.writer, line)
		logger.sync.Unlock()
	}

	// Always call Logf on each tee logger (they handle their own level checking and formatting)
	for _, teeLogger := range logger.teeLoggers {
		teeLogger.Logf(level, format, args...)
	}
}

// LogFunc evaluates the message-producing function only if at least one logger (main or tee) has the level enabled.
func (logger Logger) LogFunc(level Level, msg func() string) {
	// Defensive nil check
	if logger.level == nil {
		return
	}

	// Check if main logger or any tee logger would accept this level
	shouldLog := *logger.level <= level
	if !shouldLog {
		for _, teeLogger := range logger.teeLoggers {
			// Nil check for tee logger level
			if teeLogger.level != nil && *teeLogger.level <= level {
				shouldLog = true
				break
			}
		}
	}

	if !shouldLog {
		return
	}

	logger.Log(level, msg())
}

// LogIf calls the provided function if at least one logger (main or tee) has the level enabled.
func (logger Logger) LogIf(level Level, log func()) {
	// Defensive nil check
	if logger.level == nil {
		return
	}

	// Check if main logger or any tee logger would accept this level
	shouldLog := *logger.level <= level
	if !shouldLog {
		for _, teeLogger := range logger.teeLoggers {
			// Nil check for tee logger level
			if teeLogger.level != nil && *teeLogger.level <= level {
				shouldLog = true
				break
			}
		}
	}

	if !shouldLog {
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
