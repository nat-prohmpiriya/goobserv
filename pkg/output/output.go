package output

import (
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

// Output represents an output handler interface
type Output interface {
	// Write writes entries to output
	Write(entries []*core.Entry) error

	// Flush flushes any buffered entries
	Flush() error

	// Close closes the output
	Close() error
}

// LogEntry represents a log entry for output
type LogEntry struct {
	TraceID      string                 `json:"trace_id,omitempty"`
	RequestID    string                 `json:"request_id,omitempty"`
	StartTime    string                 `json:"start_time"`
	EndTime      string                 `json:"end_time,omitempty"`
	Duration     float64                `json:"duration,omitempty"`
	State        string                 `json:"state"`
	Method       string                 `json:"method,omitempty"`
	OriginalPath string                 `json:"original_path,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Errors       []string               `json:"errors,omitempty"`
	Spans        []*SpanEntry           `json:"spans,omitempty"`
}

// SpanEntry represents a span entry for output
type SpanEntry struct {
	Function  string                 `json:"function"`
	StartTime string                 `json:"start_time"`
	EndTime   string                 `json:"end_time,omitempty"`
	Duration  float64                `json:"duration,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
	Output    map[string]interface{} `json:"output,omitempty"`
	Event     *EventEntry            `json:"event,omitempty"`
	SpanID    string                 `json:"span_id"`
}

// EventEntry represents an event entry for output
type EventEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}
