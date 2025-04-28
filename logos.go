package logos

import (
	"fmt"
	"io"
	"sync"
)

type Fields = map[string]any

type Logger struct {
	level     *Level
	formatter Formatter
	writer    io.Writer
	sync      *sync.Mutex
	fields    Fields
	error     error
}

func NewLogger(level Level, formatter Formatter, writer io.Writer) Logger {
	return Logger{
		level:     &level,
		formatter: formatter,
		writer:    writer,
		sync:      &sync.Mutex{},
		fields:    nil,
	}
}

func (logger Logger) SetLevel(level Level) {
	*logger.level = level
}

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

func (logger Logger) With(key string, value any) Logger {
	return logger.WithFields(Fields{key: value})
}

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

// LogFunc - use for expensive operations where you don't want to calculate the message if the level is not enabled.
func (logger Logger) LogFunc(level Level, msg func() string) {
	if *logger.level > level {
		return
	}
	logger.Log(level, msg())
}

func (logger Logger) LogIf(level Level, log func()) {
	if *logger.level > level {
		return
	}
	log()
}

func (logger Logger) Print(a ...any) {
	logger.Log(LevelPrint, a...)
}

func (logger Logger) Printf(format string, args ...any) {
	logger.Logf(LevelPrint, format, args...)
}

func (logger Logger) Debug(a ...any) {
	logger.Log(LevelDebug, a...)
}

func (logger Logger) Debugf(format string, args ...any) {
	logger.Logf(LevelDebug, format, args...)
}

func (logger Logger) Info(a ...any) {
	logger.Log(LevelInfo, a...)
}

func (logger Logger) Infof(format string, args ...any) {
	logger.Logf(LevelInfo, format, args...)
}

func (logger Logger) Warn(a ...any) {
	logger.Log(LevelWarn, a...)
}

func (logger Logger) Warnf(format string, args ...any) {
	logger.Logf(LevelWarn, format, args...)
}

func (logger Logger) Error(a ...any) {
	logger.Log(LevelError, a...)
}

func (logger Logger) Errorf(format string, args ...any) {
	logger.Logf(LevelError, format, args...)
}

func (logger Logger) Fatal(a ...any) {
	logger.Log(LevelFatal, a...)
	panic(fmt.Sprint(a...))
}

func (logger Logger) Fatalf(format string, args ...any) {
	logger.Logf(LevelFatal, format, args...)
	panic(fmt.Sprintf(format, args...))
}

type Entry struct {
	Fields Fields
	Msg    string
	Error  error
}
