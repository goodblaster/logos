package logos

import (
	"time"
)

// Formatter is the interface implemented by all formatters in the package.
// It defines how a log entry is rendered as a string.
type Formatter interface {
	Format(level Level, entry Entry) string
}

// Config defines the configuration for a Formatter, such as a custom timestamp function.
type Config struct {
	Timestamp func() string
}

// DefaultConfig is the fallback configuration using the DefaultTimestamp function.
var DefaultConfig = Config{
	Timestamp: DefaultTimestamp,
}

// NewFormatter returns a new Formatter based on the provided Format and the default configuration.
func NewFormatter(format Format) Formatter {
	return NewFormatterWithConfig(format, DefaultConfig)
}

// NewFormatterWithConfig returns a new Formatter based on the provided Format and Config.
func NewFormatterWithConfig(format Format, cfg Config) Formatter {
	switch format {
	case FormatJSON:
		return NewJsonFormatter(cfg)
	case FormatText:
		return NewTextFormatter(cfg)
	case FormatConsole:
		return NewConsoleFormatter(cfg)
	}
	panic("unknown format")
}

// DefaultTimestampFormat defines the layout used for default timestamps.
const DefaultTimestampFormat = "2006-01-02T15:04:05"

// DefaultTimestamp returns the current local time formatted using DefaultTimestampFormat.
func DefaultTimestamp() string {
	return time.Now().Local().Format(DefaultTimestampFormat)
}

// JSONFormatter returns a new JSON formatter with the default configuration.
func JSONFormatter() Formatter {
	return NewJsonFormatter(DefaultConfig)
}

// TextFormatter returns a new plain text formatter with the default configuration.
func TextFormatter() Formatter {
	return NewTextFormatter(DefaultConfig)
}

// ConsoleFormatter returns a new colorized console formatter with the default configuration.
func ConsoleFormatter() Formatter {
	return NewConsoleFormatter(DefaultConfig)
}
