# logos

---
## λόγος (logos)
### Meaning "word," "speech," "reason," or "account"

---

**Logos** is a lightweight, flexible logging library for Go. It focuses on simplicity, clarity, and customizability—especially when it comes to log levels, output formatting, and structured logging.

## Getting Started

### Basic Usage
```go
import (
    "os"
    "github.com/goodblaster/logos"
)

func main() {
    // Create a logger
    log := logos.NewLogger(logos.LevelDebug, logos.ConsoleFormatter(), os.Stdout)
    log.Debug("Starting app...")
    log.Info("Application initialized")
    log.With("port", 8080).Info("Server starting")
}
```

### Using the Global Logger
The package provides a default global logger that can be used without creating an instance:
```go
// Use the default logger directly
logos.Info("This uses the default logger")
logos.With("key", "value").Debug("Debug message with field")

// Or customize the default logger
logos.SetLevel(logos.LevelInfo)
logos.SetDefaultLogger(customLogger)
```

## Environment Variables
Logos respects environment variables for easy configuration:

- `LOG_LEVEL`: Set the default log level (debug, info, warn, error, fatal)
- `LOG_FORMAT`: Set the default format (console, text, json)

```bash
LOG_LEVEL=info LOG_FORMAT=json ./myapp
```

## Features
- Easily adjustable log levels
- Structured field and error logging
- Multiple built-in formatters (Text, JSON, Console)
- Global and per-instance logging
- Context-aware logging
- Tee logging (write to multiple destinations)
- Simple customization of level names and colors
- Lazy evaluation and conditional logging

## Custom Log Levels
Unlike many logging packages, Logos allows you to rename or define your own log levels easily:

```go
const (
    LevelApple logos.Level = iota
    LevelBanana
    LevelCherry
)

logos.LevelNames = map[logos.Level]string{
    LevelApple:  "apple",
    LevelBanana: "banana",
    LevelCherry: "cherry",
}

logos.LevelColors = map[logos.Level]logos.Color{
    LevelApple:  logos.ColorTextGreen,
    LevelBanana: logos.ColorTextYellow,
    LevelCherry: logos.ColorTextRed,
}

log := logos.NewLogger(LevelApple, logos.ConsoleFormatter(), os.Stdout)
log.Log(LevelApple, "apple log")
log.Log(LevelBanana, "banana log")
log.Log(LevelCherry, "cherry log")
```

## Adding Fields and Errors
```go
log.With("user_id", 42).Info("User logged in")
log.WithFields(map[string]any{"path": "/login", "method": "POST"}).Info("Handling request")
log.WithError(err).Error("Something went wrong")

// Fields can be chained
log.With("user_id", 42).With("session_id", "abc123").Info("User action")
```

## Formatters
You can choose how logs are rendered:
- `FormatConsole` — colorized terminal output
- `FormatText` — plain, human-readable text
- `FormatJSON` — structured JSON for machines

## Conditional and Lazy Logging
```go
log.LogFunc(logos.LevelDebug, func() string {
    return expensiveComputation()
})

log.LogIf(logos.LevelInfo, func() {
    fmt.Println("This block runs only if info level is enabled")
})
```

## Context Logging
Logos supports storing and retrieving loggers from Go's `context.Context`, making it easy to pass request-scoped loggers through your application:

```go
// Create a logger with request-specific fields
requestLogger := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout).
    With("request_id", "req-123").
    With("user_id", "user-456")

// Store logger in context
ctx := logos.WithLogger(context.Background(), requestLogger)

// In any function, retrieve the logger from context
func handleRequest(ctx context.Context) {
    logger := logos.FromContext(ctx)
    logger.Info("Processing request")
    logger.With("action", "validate").Info("Validating input")
}

// If no logger is in context, FromContext returns the DefaultLogger
logger := logos.FromContext(context.Background())
logger.Info("Using default logger")
```

## Tee Logging
Write logs to multiple destinations simultaneously. Useful for logging to both console and file, or multiple files with different formats:

```go
// Create a logger that writes to both stdout and a file
file, _ := os.Create("app.log")
teeLogger := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout).
    WithTee(file)

teeLogger.Info("This message goes to both stdout and the file")

// Tee to multiple destinations
consoleBuf := &bytes.Buffer{}
fileBuf := &bytes.Buffer{}
teeLogger := logos.NewLogger(logos.LevelInfo, logos.ConsoleFormatter(), os.Stdout).
    WithTee(consoleBuf, fileBuf)

// Package-level tee logging
teeLogger := logos.WithTee(file1, file2)
teeLogger.Info("Message goes to DefaultLogger's writer plus file1 and file2")
```

## Default Log Levels
- `LevelDebug`
- `LevelInfo`
- `LevelWarn`
- `LevelError`
- `LevelFatal`
- `LevelPrint`

These can be extended or overridden to suit your needs.

---
