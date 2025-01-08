package output

import (
	"fmt"
	"os"
	"time"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

// StdoutConfig represents stdout output configuration
type StdoutConfig struct {
	Colored bool // Enable colored output
}

// StdoutOutput represents stdout output handler
type StdoutOutput struct {
	config StdoutConfig
}

// NewStdoutOutput creates a new stdout output handler
func NewStdoutOutput(config StdoutConfig) *StdoutOutput {
	return &StdoutOutput{
		config: config,
	}
}

// Write writes entries to stdout
func (o *StdoutOutput) Write(entries []*core.Entry) error {
	for _, entry := range entries {
		// Prepare color codes
		var color, reset string
		if o.config.Colored {
			switch entry.Level {
			case core.LevelDebug:
				color = "\033[36m" // Cyan
			case core.LevelInfo:
				color = "\033[32m" // Green
			case core.LevelWarn:
				color = "\033[33m" // Yellow
			case core.LevelError:
				color = "\033[31m" // Red
			}
			if color != "" {
				reset = "\033[0m"
			}
		}

		// Format log message
		msg := fmt.Sprintf("%s[%s] %s | trace=%s span=%s",
			color,
			entry.Level.String(),
			entry.Message,
			entry.TraceID,
			entry.SpanID,
		)

		// Add error if present
		if entry.Data != nil {
			if err, ok := entry.Data["error"]; ok {
				msg += fmt.Sprintf(" | error=%v", err)
			}
		}

		// Add reset code
		msg += reset

		// Write to stdout with timestamp
		fmt.Fprintf(os.Stdout, "%s %s\n",
			entry.Time.Format(time.RFC3339),
			msg,
		)
	}

	return nil
}

// Flush flushes the output
func (o *StdoutOutput) Flush() error {
	return nil
}

// Close closes the output
func (o *StdoutOutput) Close() error {
	return nil
}
