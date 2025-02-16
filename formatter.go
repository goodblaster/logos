package logos

import (
	"time"
)

type Formatter interface {
	Format(level Level, entry Entry) string
}

type Config struct {
	Timestamp func() string
}

var DefaultConfig = Config{
	Timestamp: DefaultTimestamp,
}

func NewFormatter(format Format) Formatter {
	return NewFormatterWithConfig(format, DefaultConfig)
}

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

const DefaultTimestampFormat = "2006-01-02T15:04:05"

func DefaultTimestamp() string {
	return time.Now().Local().Format(DefaultTimestampFormat)
}

func JSONFormatter() Formatter {
	return NewJsonFormatter(DefaultConfig)
}

func TextFormatter() Formatter {
	return NewTextFormatter(DefaultConfig)
}

func ConsoleFormatter() Formatter {
	return NewConsoleFormatter(DefaultConfig)
}
