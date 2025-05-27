# logos

---
## λόγος (logos)
### Meaning "word," "speech," "reason," or "account"

---

**Logos** is a lightweight, flexible logging library for Go. It focuses on simplicity, clarity, and customizability—especially when it comes to log levels, output formatting, and structured logging.

## Getting Started

```go
import (
    "os"
    "github.com/goodblaster/logos"
)

func main() {
    log := logos.NewLogger(logos.LevelDebug, logos.FormatConsole, os.Stdout)
    log.Debug("Starting app...")
}
```

## Global Logger
You can set a global logger to avoid passing it throughout your code:

```go
logos.SetDefaultLogger(log)
logos.Info("This uses the default logger")
```

## Features
- Easily adjustable log levels
- Structured field and error logging
- Multiple built-in formatters (Text, JSON, Console)
- Global and per-instance logging
- Simple customization of level names and colors

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

log := logos.NewLogger(LevelApple, logos.FormatConsole, os.Stdout)
log.Log(LevelApple, "apple log")
log.Log(LevelBanana, "banana log")
log.Log(LevelCherry, "cherry log")
```

## Adding Fields and Errors
```go
log.With("user_id", 42).Info("User logged in")
log.WithFields(map[string]any{"path": "/login", "method": "POST"}).Info("Handling request")
log.WithError(err).Error("Something went wrong")
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

## Default Log Levels
- `LevelDebug`
- `LevelInfo`
- `LevelWarn`
- `LevelError`
- `LevelFatal`
- `LevelPrint`

These can be extended or overridden to suit your needs.

---
