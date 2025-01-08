package core

import (
	"time"
	"your-project/pkg/core/level"
)

// EntryType represents the type of entry
type EntryType string

const (
	LogEntry    EntryType = "log"
	MetricEntry EntryType = "metric"
	TraceEntry  EntryType = "trace"
)

// Entry represents a log entry
type Entry struct {
	Time      time.Time              `json:"time"`
	Level     level.Level            `json:"level"`
	Message   string                 `json:"message"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
	ParentID  string                 `json:"parent_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Type      EntryType             `json:"type,omitempty"`
}

// NewEntry creates a new entry
func NewEntry() *Entry {
	return &Entry{
		Time: time.Now(),
		Data: make(map[string]interface{}),
		Type: LogEntry,
	}
}

// WithField adds a field to the entry
func (e *Entry) WithField(key string, value interface{}) *Entry {
	e.Data[key] = value
	return e
}

// WithError adds an error to the entry
func (e *Entry) WithError(err error) *Entry {
	if err != nil {
		e.Data["error"] = err.Error()
	}
	return e
}
