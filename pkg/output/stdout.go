package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/vongga-platform/goobserv/pkg/core"
)

// StdoutOutput represents console output handler
type StdoutOutput struct {
	colored bool
	mu      sync.Mutex
}

// StdoutConfig represents stdout configuration
type StdoutConfig struct {
	Colored bool
}

// NewStdoutOutput creates a new stdout output
func NewStdoutOutput(config StdoutConfig) *StdoutOutput {
	return &StdoutOutput{
		colored: config.Colored,
	}
}

// Write writes an entry to stdout
func (o *StdoutOutput) Write(entry *core.Entry) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Format entry
	data := map[string]interface{}{
		"timestamp": entry.Timestamp,
		"level":     entry.Level.String(),
		"message":   entry.Message,
	}

	// Add context if available
	if entry.Context != nil {
		if entry.Context.TraceID() != "" {
			data["trace_id"] = entry.Context.TraceID()
		}
		if entry.Context.SpanID() != "" {
			data["span_id"] = entry.Context.SpanID()
		}
		if attrs := entry.Context.Attributes(); len(attrs) > 0 {
			data["attributes"] = attrs
		}
	}

	// Add extra fields
	if len(entry.Data) > 0 {
		data["data"] = entry.Data
	}

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write to stdout
	if o.colored {
		color := o.levelColor(entry.Level)
		fmt.Fprintf(os.Stdout, "%s%s%s\n", color, string(jsonData), "\033[0m")
	} else {
		fmt.Fprintln(os.Stdout, string(jsonData))
	}

	return nil
}

// Flush implements Output interface
func (o *StdoutOutput) Flush() error {
	return nil
}

// Close implements Output interface
func (o *StdoutOutput) Close() error {
	return nil
}

// levelColor returns ANSI color code for log level
func (o *StdoutOutput) levelColor(level core.Level) string {
	switch level {
	case core.DebugLevel:
		return "\033[37m" // white
	case core.InfoLevel:
		return "\033[32m" // green
	case core.WarnLevel:
		return "\033[33m" // yellow
	case core.ErrorLevel:
		return "\033[31m" // red
	default:
		return "\033[0m" // default
	}
}
