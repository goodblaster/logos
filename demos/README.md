# Logos Demos

This directory contains demonstrations of various logos logging library features.

## Quick Start

Run any demo with:
```bash
cd demos/<demo-name>
go run main.go
```

## Demo Index

### Core Functionality

- **[master](./master/)** - **START HERE!** Comprehensive demo showing the most commonly used features that most projects will need.

### Individual Feature Demos

1. **[01_basic](./01_basic/)** - Basic logging operations (Debug, Info, Warn, Error)
2. **[02_fields](./02_fields/)** - Structured logging with fields using With() and WithFields()
3. **[03_levels](./03_levels/)** - Log levels, filtering, and IsLevelEnabled()
4. **[04_formatters](./04_formatters/)** - Different output formats: JSON, Text, Console
5. **[05_context](./05_context/)** - Context-based logging for request-scoped loggers
6. **[06_tee](./06_tee/)** - Tee logging to multiple destinations simultaneously
7. **[07_errors](./07_errors/)** - Error handling with WithError() and error handlers
8. **[08_lazy_evaluation](./08_lazy_evaluation/)** - LogFunc and LogIf for performance optimization
9. **[09_custom_levels](./09_custom_levels/)** - Creating and using custom log levels
10. **[10_custom_formatter](./10_custom_formatter/)** - Implementing custom formatters

### Advanced

- **[comprehensive_demo](./comprehensive_demo/)** - Exhaustive demo showing all library features including custom levels, custom formatters, and advanced use cases

## Recommended Learning Path

1. Start with **master/** to understand core concepts
2. Explore individual demos based on your needs
3. Check **comprehensive_demo/** for advanced use cases

## Running All Demos

To run all demos in sequence:

```bash
for dir in demos/*/; do
  echo "Running $(basename "$dir")..."
  (cd "$dir" && go run main.go)
  echo ""
done
```
