# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Logos** is a lightweight, flexible logging library for Go that emphasizes simplicity, customizability, and structured logging. The package name comes from the Greek λόγος (logos), meaning "word," "speech," "reason," or "account."

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v -run TestName

# Run tests for a specific file
go test -v formatter_json_test.go formatter_json.go formatter.go levels.go colors.go formats.go
```

### Building
```bash
# Build the package
go build ./...

# Run the demo
go run demo/main.go
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

## Architecture

### Core Components

**Logger (`logos.go`)**
- The main `Logger` struct is immutable through copy-on-write semantics
- Each method that modifies state (`With`, `WithFields`, `WithError`) returns a new logger
- Thread-safe through mutex-protected operations
- Supports both instance-based and package-level logging via `DefaultLogger`

**Level System (`levels.go`)**
- Levels are simple `int` constants, allowing custom levels beyond the defaults
- Default levels: Debug (-1), Info (0), Warn (1), Error (2), Fatal (3), Print (MaxInt)
- `LevelNames` and `LevelColors` are global maps that can be overridden for custom level schemes
- Level filtering: a logger with level X will log messages at level X and higher

**Formatters (`formatter.go`, `formatter_*.go`)**
- Three built-in formatters: JSON, Text, and Console (ANSI-colored)
- All implement the `Formatter` interface with a single `Format(Level, Entry) string` method
- Formatters are stateless except for the `Config` (which holds timestamp function)
- Each formatter handles the `Entry` struct containing Fields, Msg, and Error

**Tee Logging (`tee.go`)**
- Loggers can be composed using the `Tee()` method to write to multiple destinations
- Each tee logger maintains its own level, formatter, and fields
- Log messages are sent to all tee destinations; each destination does its own filtering
- Critical implementation detail: tee loggers are stored in the `teeLoggers` slice and each handles its own level checking

**Context Integration (`context.go`)**
- Loggers can be stored in and retrieved from `context.Context`
- Uses a private `contextKey` type to avoid collisions
- `FromContext()` returns `DefaultLogger` if no logger is found in context
- Useful for request-scoped logging in HTTP handlers and middleware

**Global Default Logger (`defaults.go`)**
- Package exports convenience functions (`Info()`, `Debug()`, etc.) that use an internal default logger
- Default logger is initialized from environment variables:
  - `LOG_LEVEL`: debug, info, warn, error, fatal
  - `LOG_FORMAT`: console, text, json
- Can be replaced with `SetDefaultLogger()` (thread-safe with RWMutex protection)
- All package-level functions are thread-safe through mutex-protected access

### Key Design Patterns

**Immutability Through Copying**
- `With()`, `WithFields()`, and `WithError()` create new loggers with copied state
- The `Copy()` method performs deep copies of fields and tee loggers
- This prevents accidental sharing of mutable state between loggers

**Lazy Evaluation**
- `LogFunc()` takes a `func() string` that's only called if the level is enabled
- `LogIf()` takes a `func()` that runs only if the level is enabled
- Both check the main logger AND all tee loggers to determine if evaluation is needed

**Entry-Based Formatting**
- All formatters receive an `Entry` struct containing Fields, Msg, and Error
- Separates logging logic from formatting logic cleanly
- Makes it easy to add new formatters without changing core logging code

### File Organization

- `logos.go` - Core Logger type and methods
- `levels.go` - Level definitions and customization
- `formatter.go` - Formatter interface and factory functions
- `formatter_*.go` - Specific formatter implementations (JSON, Text, Console)
- `tee.go` - Multi-destination logging
- `context.go` - Context integration
- `defaults.go` - Package-level convenience functions and DefaultLogger
- `colors.go` - ANSI color code constants
- `formats.go` - Format type enum

### Testing Patterns

Tests use `bytes.Buffer` as the writer to capture output for assertions:

```go
buf := &bytes.Buffer{}
log := NewLogger(LevelDebug, JSONFormatter(), buf)
log.Info("message")
// Parse buf.String() to verify output
```

## Important Implementation Notes

- **Fatal logs panic**: The `Fatal()` and `Fatalf()` methods log the message and then call `panic()`
- **Print level always prints**: `LevelPrint` is set to `math.MaxInt` so it bypasses all filtering
- **Error overwriting**: `WithError()` will warn if replacing an existing error
- **Tee propagation**: When calling `With()` on a logger with tees, fields are added to the main logger only, not to tee destinations
- **Custom levels**: Users can define their own levels and update `LevelNames` and `LevelColors` maps
- **Timestamp customization**: Formatters accept a `Config` with a custom `Timestamp func() string`
- **Thread-safety**:
  - Logger writes are protected by per-logger mutex (safe for concurrent logging to same writer)
  - Default logger access is protected by RWMutex (safe for concurrent reads, exclusive writes)
  - Each logger instance has its own mutex, preventing deadlocks

## Dependencies

- `github.com/goodblaster/errors` - Used in JSON formatter for error wrapping
- `github.com/stretchr/testify` - Used in tests for assertions
