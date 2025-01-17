package core

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// SpanStatus represents the status of a span
type SpanStatus int

const (
	SpanStatusOK SpanStatus = iota
	SpanStatusError
)

// Span represents a single operation within a trace
type Span struct {
	TraceID    string
	SpanID     string
	ParentID   string
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	Status     SpanStatus
	Attributes map[string]interface{}
	Events     []SpanEvent
	mu         sync.RWMutex
}

// SpanEvent represents an event within a span
type SpanEvent struct {
	Time       time.Time
	Name       string
	Attributes map[string]interface{}
}

// TraceID generates a new trace ID
func TraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// SpanID generates a new span ID
func SpanID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// RequestID generates a new request ID
func RequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// NewSpan creates a new span
func NewSpan(name string) *Span {
	return &Span{
		TraceID:    TraceID(),
		SpanID:     SpanID(),
		ParentID:   "",
		Name:       name,
		StartTime:  time.Now(),
		Status:     SpanStatusOK,
		Attributes: make(map[string]interface{}),
		Events:     make([]SpanEvent, 0),
	}
}

// End ends the span
func (s *Span) End() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EndTime = time.Now()
}

// SetStatus sets the span status
func (s *Span) SetStatus(status SpanStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = status
}

// SetAttribute sets a span attribute
func (s *Span) SetAttribute(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Attributes[key] = value
}

// AddEvent adds an event to the span
func (s *Span) AddEvent(name string, attrs map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Events = append(s.Events, SpanEvent{
		Time:       time.Now(),
		Name:       name,
		Attributes: attrs,
	})
}

// Duration returns the span duration
func (s *Span) Duration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}
