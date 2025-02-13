# logos

---
## λόγος (logos)
### Meaning "word," "speech," "reason," or "account"

---

Logos is a very lightweight, simple logger with all the functionality I ever need.
It is simple to use and modify.
Getting started:
```go
    log := logos.NewLogger(LevelDebug, FormatConsole, os.Stdout)
    log.Debug("logos ...")
```
You can set it as a default logger so you don't have to pass the logger around.
```go
    logos.SetDefaultLogger(log)
    logos.Debug("default ...")
```

If you want to change log levels, and names, do whatever you like.
Note that this is a global change. 
If you are running multiple loggers, this could cause a problem.
Example:

```go
func main() {
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

    logos.LevelColors = map[logos.Level]logos.TextColor{
        LevelApple:  colors.TextGreen,
        LevelBanana: colors.TextYellow,
        LevelCherry: colors.TextRed,
    }
    
    log := logos.NewLogger(LevelApple, FormatConsole, os.Stdout)
    log.Log(LevelApple, "apple ...")
    log.Log(LevelBanana, "banana ...")
    log.Log(LevelCherry, "cherry ...")
}
```
