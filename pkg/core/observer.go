package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Config represents observer configuration
type Config struct {
	BufferSize    int
	FlushInterval time.Duration
}

// Observer represents an observer instance
type Observer struct {
	outputs []Output
	buffer  chan *Entry
	done   chan struct{}
	closed bool
	mu     sync.Mutex
	wg     sync.WaitGroup

	// Metrics
	metrics map[string]interface{}
}

// NewObserver creates a new observer instance
func NewObserver(cfg Config) *Observer {
	// Ensure minimum flush interval
	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = 100 * time.Millisecond
	}

	// Set default buffer size
	// Calculate based on expected log rate and flush interval
	// Assuming max 1000 logs/second with 2x safety factor
	if cfg.BufferSize <= 0 {
		expectedLogsPerSecond := 1000
		safetyFactor := 2
		cfg.BufferSize = int(float64(expectedLogsPerSecond) * cfg.FlushInterval.Seconds() * float64(safetyFactor))
		if cfg.BufferSize < 100 {
			cfg.BufferSize = 100  // Minimum buffer size
		}
	}

	obs := &Observer{
		buffer: make(chan *Entry, cfg.BufferSize),
		done:   make(chan struct{}),
		metrics: make(map[string]interface{}),
	}

	// Start worker
	obs.wg.Add(1)
	go obs.worker(cfg.FlushInterval)

	return obs
}

// AddOutput adds an output handler
func (o *Observer) AddOutput(output Output) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.outputs = append(o.outputs, output)
}

// Counter creates a new counter metric
func (o *Observer) Counter(name string) *Counter {
	o.mu.Lock()
	defer o.mu.Unlock()

	if m, ok := o.metrics[name]; ok {
		if counter, ok := m.(*Counter); ok {
			return counter
		}
	}

	counter := NewCounter(name)
	o.metrics[name] = counter
	return counter
}

// Gauge creates a new gauge metric
func (o *Observer) Gauge(name string) *Gauge {
	o.mu.Lock()
	defer o.mu.Unlock()

	if m, ok := o.metrics[name]; ok {
		if gauge, ok := m.(*Gauge); ok {
			return gauge
		}
	}

	gauge := NewGauge(name)
	o.metrics[name] = gauge
	return gauge
}

// Histogram creates a new histogram metric
func (o *Observer) Histogram(name string, buckets []float64) *Histogram {
	o.mu.Lock()
	defer o.mu.Unlock()

	if m, ok := o.metrics[name]; ok {
		if histogram, ok := m.(*Histogram); ok {
			return histogram
		}
	}

	histogram := NewHistogram(name, buckets)
	o.metrics[name] = histogram
	return histogram
}

// StartSpan starts a new span
func (o *Observer) StartSpan(ctx context.Context, name string) (*Span, context.Context) {
	span := NewSpan(name)
	span.TraceID = TraceID()
	span.SpanID = SpanID()

	if parentSpan, ok := ctx.Value(spanKey).(*Span); ok {
		span.ParentID = parentSpan.SpanID
		span.TraceID = parentSpan.TraceID
	}

	return span, context.WithValue(ctx, spanKey, span)
}

// EndSpan ends the current span
func (o *Observer) EndSpan(span *Span) {
	span.End()
}

// Debug logs a debug message
func (o *Observer) Debug(ctx context.Context, msg string) *Entry {
	return o.log(ctx, LevelDebug, msg)
}

// Info logs an info message
func (o *Observer) Info(ctx context.Context, msg string) *Entry {
	return o.log(ctx, LevelInfo, msg)
}

// Warn logs a warning message
func (o *Observer) Warn(ctx context.Context, msg string) *Entry {
	return o.log(ctx, LevelWarn, msg)
}

// Error logs an error message
func (o *Observer) Error(ctx context.Context, msg string) *Entry {
	return o.log(ctx, LevelError, msg)
}

// Close closes the observer
func (o *Observer) Close() error {
	o.mu.Lock()
	if o.closed {
		o.mu.Unlock()
		return nil
	}
	o.closed = true
	o.mu.Unlock()

	// Signal worker to stop
	close(o.done)

	// Wait for worker to finish
	o.wg.Wait()

	// Close outputs
	for _, output := range o.outputs {
		if err := output.Close(); err != nil {
			return fmt.Errorf("failed to close output: %v", err)
		}
	}

	return nil
}

// log logs a message with the given level
func (o *Observer) log(ctx context.Context, level Level, msg string) *Entry {
	entry := &Entry{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
		Data:    make(map[string]interface{}),
	}

	// Add span information
	if span, ok := ctx.Value(spanKey).(*Span); ok {
		entry.TraceID = span.TraceID
		entry.SpanID = span.SpanID
	}

	// Add request ID
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		entry.RequestID = reqID
	}

	// Send to buffer
	select {
	case o.buffer <- entry:
	default:
		// Buffer is full, drop entry
	}

	return entry
}

// worker processes entries from buffer
func (o *Observer) worker(interval time.Duration) {
	defer o.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var entries []*Entry

	for {
		select {
		case <-o.done:
			if len(entries) > 0 {
				o.flush(entries)
			}
			return
		case entry := <-o.buffer:
			entries = append(entries, entry)
			if len(entries) >= cap(o.buffer) {
				o.flush(entries)
				entries = entries[:0]
			}
		case <-ticker.C:
			if len(entries) > 0 {
				o.flush(entries)
				entries = entries[:0]
			}
		}
	}
}

// flush flushes entries to outputs
func (o *Observer) flush(entries []*Entry) {
	o.mu.Lock()
	outputs := make([]Output, len(o.outputs))
	copy(outputs, o.outputs)
	o.mu.Unlock()

	for _, output := range outputs {
		if err := output.Write(entries); err != nil {
			fmt.Printf("failed to write entries: %v\n", err)
		}
	}
}

// Flush flushes any pending entries
func (o *Observer) Flush() error {
	o.mu.Lock()
	outputs := make([]Output, len(o.outputs))
	copy(outputs, o.outputs)
	o.mu.Unlock()

	for _, output := range outputs {
		if err := output.Flush(); err != nil {
			return fmt.Errorf("failed to flush output: %v", err)
		}
	}

	return nil
}
