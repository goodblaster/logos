package logos

// Format represents the output format used by the logger.
type Format int

// String returns the string representation of a Format.
func (f Format) String() string {
	if name, ok := FormatNames[f]; ok {
		return name
	}
	return "unknown"
}

const (
	// FormatJSON outputs logs in structured JSON format.
	FormatJSON Format = iota
	// FormatText outputs logs as plain, uncolored text.
	FormatText
	// FormatConsole outputs logs as colored text suitable for terminals.
	FormatConsole
)

// Formats is the list of all supported output formats.
var Formats = []Format{
	FormatJSON,
	FormatText,
	FormatConsole,
}

// FormatNames maps Format values to their string identifiers.
var FormatNames = map[Format]string{
	FormatJSON:    "JSON",
	FormatText:    "TEXT",
	FormatConsole: "CONSOLE",
}
