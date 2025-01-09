package output

import (
	"encoding/json"
	"os"
	"time"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

// StdoutConfig represents stdout output configuration
type StdoutConfig struct {
	Pretty bool
}

// StdoutOutput represents stdout output handler
type StdoutOutput struct {
	stdout *os.File
	config StdoutConfig
}

// NewStdoutOutput creates a new stdout output handler
func NewStdoutOutput(config StdoutConfig) *StdoutOutput {
	return &StdoutOutput{
		stdout: os.Stdout,
		config: config,
	}
}

// Write writes entries to stdout
func (o *StdoutOutput) Write(entries []*core.Entry) error {
	for _, entry := range entries {
		// Create log entry
		logEntry := LogEntry{
			TraceID:      entry.TraceID,
			RequestID:    entry.RequestID,
			StartTime:    entry.StartTime.Format(time.RFC3339),
			EndTime:      entry.EndTime.Format(time.RFC3339),
			Duration:     entry.Duration,
			State:       entry.State,
			Method:      entry.Method,
			OriginalPath: entry.OriginalPath,
			Metadata:    make(map[string]interface{}),
			Errors:      make([]string, 0),
		}

		// Convert spans
		for _, span := range entry.Spans {
			spanEntry := &SpanEntry{
				Function:  span.Function,
				StartTime: span.StartTime.Format(time.RFC3339),
				EndTime:   span.EndTime.Format(time.RFC3339),
				Duration:  span.Duration,
				Input:     span.Input,
				Output:    span.Output,
				SpanID:    span.SpanID,
			}
			
			if span.Event != nil {
				spanEntry.Event = &EventEntry{
					Level:   span.Event.Level,
					Message: span.Event.Message,
				}
			}
			
			logEntry.Spans = append(logEntry.Spans, spanEntry)
		}

		// Add error if present
		if entry.Error != nil {
			logEntry.Errors = append(logEntry.Errors, entry.Error.Message)
			logEntry.State = "error"
		}

		// Encode and write
		var err error
		if o.config.Pretty {
			encoder := json.NewEncoder(o.stdout)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(logEntry)
		} else {
			err = json.NewEncoder(o.stdout).Encode(logEntry)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// Flush flushes stdout
func (o *StdoutOutput) Flush() error {
	return nil
}

// Close closes stdout
func (o *StdoutOutput) Close() error {
	return nil
}
