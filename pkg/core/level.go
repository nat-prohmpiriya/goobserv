package core

// Level represents log level
type Level string

const (
	// LevelDebug represents debug level
	LevelDebug Level = "debug"
	// LevelInfo represents info level
	LevelInfo Level = "info"
	// LevelWarn represents warn level
	LevelWarn Level = "warn"
	// LevelError represents error level
	LevelError Level = "error"
)

// String returns string representation of level
func (l Level) String() string {
	return string(l)
}
