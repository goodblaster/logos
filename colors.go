package logos

type Color string

const (
	ColorReset Color = "\033[0m"

	ColorTextRed     Color = "\033[31m"
	ColorTextYellow  Color = "\033[33m"
	ColorTextGreen   Color = "\033[32m"
	ColorTextBlue    Color = "\033[34m"
	ColorTextMagenta Color = "\033[35m"
	ColorTextCyan    Color = "\033[36m"
	ColorTextWhite   Color = "\033[37m"
	ColorTextPurple  Color = "\033[35m"
	ColorTextBlack   Color = "\033[30m"

	ColorBgRed     Color = "\033[41m"
	ColorBgYellow  Color = "\033[43m"
	ColorBgGreen   Color = "\033[42m"
	ColorBgBlue    Color = "\033[44m"
	ColorBgMagenta Color = "\033[45m"
	ColorBgCyan    Color = "\033[46m"
	ColorBgWhite   Color = "\033[47m"
	ColorBgBlack   Color = "\033[40m"
)
