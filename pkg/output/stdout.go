package output

import (
	"encoding/json"
	"os"

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
		// Prepare log data
		data := make(map[string]interface{})
		if entry.Data != nil {
			data = entry.Data
		}
		if entry.TraceID != "" {
			data["trace_id"] = entry.TraceID
		}
		if entry.SpanID != "" {
			data["span_id"] = entry.SpanID
		}
		if entry.RequestID != "" {
			data["request_id"] = entry.RequestID
		}

		// Create log entry
		logEntry := map[string]interface{}{
			"timestamp": entry.Time.Format("2006-01-02T15:04:05.000Z07:00"),
			"level":     entry.Level.String(),
			"message":   entry.Message,
			"data":      data,
		}

		// Add color if enabled
		if o.config.Colored {
			switch entry.Level {
			case core.LevelDebug:
				logEntry["color"] = "\033[36m" // Cyan
			case core.LevelInfo:
				logEntry["color"] = "\033[32m" // Green
			case core.LevelWarn:
				logEntry["color"] = "\033[33m" // Yellow
			case core.LevelError:
				logEntry["color"] = "\033[31m" // Red
			}
			if _, ok := logEntry["color"]; ok {
				logEntry["reset"] = "\033[0m"
			}
		}

		// Encode and write
		if err := json.NewEncoder(os.Stdout).Encode(logEntry); err != nil {
			return err
		}
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
