package logos

// Tee adds one or more loggers as tee destinations.
// Each tee logger will receive all log messages and handle its own level checking and formatting.
// This allows different destinations to have different levels, formatters, and fields.
func (logger Logger) Tee(loggers ...Logger) Logger {
	if len(loggers) == 0 {
		return logger
	}

	newLogger := logger.Copy()
	if newLogger.teeLoggers == nil {
		newLogger.teeLoggers = make([]Logger, 0, len(loggers))
	}

	// Copy each tee logger to avoid shared state
	for _, teeLogger := range loggers {
		newLogger.teeLoggers = append(newLogger.teeLoggers, teeLogger.Copy())
	}

	return newLogger
}
