package logos

import (
	"io"
)

// WithTee returns a new Logger that writes to both the existing writer and the provided writers.
// All writers will receive the same log output.
func (logger Logger) WithTee(writers ...io.Writer) Logger {
	if len(writers) == 0 {
		return logger
	}

	// Combine existing writer with new writers
	allWriters := make([]io.Writer, 0, len(writers)+1)
	allWriters = append(allWriters, logger.writer)
	allWriters = append(allWriters, writers...)

	newLogger := logger.Copy()
	newLogger.writer = io.MultiWriter(allWriters...)
	return newLogger
}
