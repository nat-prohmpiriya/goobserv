package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

// StdoutOutput represents stdout output handler
type StdoutOutput struct{}

// NewStdoutOutput creates a new stdout output handler
func NewStdoutOutput() *StdoutOutput {
	return &StdoutOutput{}
}

// Write writes entries to stdout
func (o *StdoutOutput) Write(entries []*core.Entry) error {
	for _, entry := range entries {
		// Format timestamp
		timestamp := entry.Time.Format(time.RFC3339)

		// Format data
		data := make(map[string]interface{})
		for k, v := range entry.Data {
			data[k] = v
		}

		// Add trace and request IDs
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
			"timestamp": timestamp,
			"level":     entry.Level.String(),
			"message":   entry.Message,
			"data":      data,
		}

		// Marshal to JSON
		jsonBytes, err := json.Marshal(logEntry)
		if err != nil {
			return fmt.Errorf("failed to marshal log entry: %v", err)
		}

		// Write to stdout
		fmt.Fprintln(os.Stdout, string(jsonBytes))
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
