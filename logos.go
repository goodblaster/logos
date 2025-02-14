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
	return logger.With("error", err)
}

func (logger Logger) Log(level Level, format string, args ...any) {
	if *logger.level > level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	line := logger.formatter.Format(level, msg, logger.fields)
	_, _ = fmt.Fprintln(logger.writer, line)
}

// LogFunc - use for expensive operations where you don't want to calculate the message if the level is not enabled.
func (logger Logger) LogFunc(level Level, msg func() string) {
	if *logger.level > level {
		return
	}
	logger.Log(level, msg())
}

func (logger Logger) Print(format string, args ...any) {
	logger.Log(LevelPrint, format, args...)
}

func (logger Logger) Debug(format string, args ...any) {
	logger.Log(LevelDebug, format, args...)
}

func (logger Logger) Info(format string, args ...any) {
	logger.Log(LevelInfo, format, args...)
}

func (logger Logger) Warn(format string, args ...any) {
	logger.Log(LevelWarn, format, args...)
}

func (logger Logger) Error(format string, args ...any) {
	logger.Log(LevelError, format, args...)
}

func (logger Logger) Fatal(format string, args ...any) {
	logger.Log(LevelFatal, format, args...)
	panic(fmt.Sprintf(format, args...))
}
