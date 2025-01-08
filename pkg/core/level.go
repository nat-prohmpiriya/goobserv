package core

// Level represents a log level
type Level string

const (
	LevelDebug = Level("debug")
	LevelInfo  = Level("info")
	LevelWarn  = Level("warn")
	LevelError = Level("error")
)
