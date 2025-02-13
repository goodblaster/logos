package logos

type Format int

func (f Format) String() string {
	if name, ok := FormatNames[f]; ok {
		return name
	}
	return "unknown"
}

const (
	FormatJSON Format = iota
	FormatText
	FormatConsole
)

var Formats = []Format{
	FormatJSON,
	FormatText,
	FormatConsole,
}

var FormatNames = map[Format]string{
	FormatJSON:    "JSON",
	FormatText:    "TEXT",
	FormatConsole: "CONSOLE",
}
