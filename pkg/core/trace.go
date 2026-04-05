package core

import (
	"time"
)

// // Event represents a log event
// type Event struct {
// 	Level   string `json:"level"` // debug, info, warn, error
// 	Message string `json:"message"`
// }

// // Span represents a function execution or manual log
// type Span struct {
// 	Function  string                 `json:"function"` // package.function
// 	StartTime time.Time              `json:"start_time"`
// 	EndTime   time.Time              `json:"end_time"`
// 	Duration  float64                `json:"duration"`
// 	Input     map[string]interface{} `json:"input,omitempty"`
// 	Output    map[string]interface{} `json:"output,omitempty"`
// 	Event     *Event                 `json:"event,omitempty"` // for manual logging
// 	SpanID    string                 `json:"span_id"`
// }

// // Error represents error details
// type Error struct {
// 	Code       string                 `json:"code"`
// 	Message    string                 `json:"message"`
// 	StackTrace string                 `json:"stack_trace,omitempty"`
// 	Details    map[string]interface{} `json:"details,omitempty"`
// }

// type entryKey struct{}

// // Log Entry represents a complete request log
// type Entry struct {
// 	RequestID    string    `json:"request_id"`
// 	TraceID      string    `json:"trace_id"`
// 	UserID       string    `json:"user_id,omitempty"`
// 	StartTime    time.Time `json:"start_time"`
// 	EndTime      time.Time `json:"end_time"`
// 	Duration     float64   `json:"duration"`
// 	State        string    `json:"state"` // processing, success, error
// 	Method       string    `json:"method"`
// 	OriginalPath string    `json:"original_path"`
// 	Spans        []Span    `json:"spans"`
// 	Error        *Error    `json:"error,omitempty"`
// }

// // NewEntry creates a new entry
// func NewEntry() *Entry {
// 	return &Entry{
// 		StartTime: time.Now(),
// 		State:     "processing",
// 		Spans:     make([]Span, 0),
// 	}
// }

// // WithEntry adds an entry to context
// func WithEntry(ctx context.Context, entry *Entry) context.Context {
// 	return context.WithValue(ctx, entryKey{}, entry)
// }

// // GetEntry gets the current entry from context
// func GetEntry(ctx context.Context) *Entry {
// 	if entry, ok := ctx.Value(entryKey{}).(*Entry); ok {
// 		return entry
// 	}
// 	return nil
// }

// // AddSpan adds a span to the entry
// func (e *Entry) AddSpan(span Span) {
// 	span.SpanID = fmt.Sprintf("%d", len(e.Spans)+1)
// 	e.Spans = append(e.Spans, span)
// }

// // End marks the entry as completed
// func (e *Entry) End() {
// 	e.EndTime = time.Now()
// 	e.Duration = e.EndTime.Sub(e.StartTime).Seconds()
// 	if e.Error == nil {
// 		e.State = "success"
// 	} else {
// 		e.State = "error"
// 	}
// }

// // WithError adds error details to the entry
// func (e *Entry) WithError(err error, code string, details map[string]interface{}) *Entry {
// 	e.Error = &Error{
// 		Code:    code,
// 		Message: err.Error(),
// 		Details: details,
// 	}
// 	if obs := GetObserver(context.Background()); obs != nil && obs.config.Development {
// 		e.Error.StackTrace = string(debug.Stack())
// 	}
// 	return e
// }

type Span struct {
}

func NewSpan() *Span {
	return &Span{}
}

type Trace struct {
	Spans     []*Span
	TraceID   string
	RequestID string
	UserID    string
	StartTime time.Time
	EndTime   time.Time
    Duration  float64
    State     string
    Method    string
    OriginalPath string
    
}

func NewTrace() *Trace {
	return &Trace{}
}
