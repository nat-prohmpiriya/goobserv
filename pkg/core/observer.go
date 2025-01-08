package core

import (
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
}

// NewObserver creates a new observer instance
func NewObserver(cfg Config) *Observer {
	// Ensure minimum flush interval
	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = 100 * time.Millisecond
	}

	obs := &Observer{
		buffer: make(chan *Entry, cfg.BufferSize),
		done:   make(chan struct{}),
	}

	// Start worker
	obs.wg.Add(1)
	go obs.worker(cfg.FlushInterval)

	return obs
}

// AddOutput adds an output handler
func (o *Observer) AddOutput(output Output) {
	o.outputs = append(o.outputs, output)
}

// StartSpan starts a new span
func (o *Observer) StartSpan(ctx *Context, name string) *Context {
	return ctx.WithSpan(name)
}

// EndSpan ends the current span
func (o *Observer) EndSpan(ctx *Context) {
	ctx.EndSpan()
}

// Debug logs a debug message
func (o *Observer) Debug(ctx *Context, msg string, args ...interface{}) {
	o.log(ctx, LevelDebug, msg, args...)
}

// Info logs an info message
func (o *Observer) Info(ctx *Context, msg string, args ...interface{}) {
	o.log(ctx, LevelInfo, msg, args...)
}

// Warn logs a warning message
func (o *Observer) Warn(ctx *Context, msg string, args ...interface{}) {
	o.log(ctx, LevelWarn, msg, args...)
}

// Error logs an error message
func (o *Observer) Error(ctx *Context, msg string, args ...interface{}) {
	o.log(ctx, LevelError, msg, args...)
}

// Close closes the observer
func (o *Observer) Close() error {
	o.mu.Lock()
	if o.closed {
		o.mu.Unlock()
		return nil
	}
	o.closed = true
	close(o.done)
	o.mu.Unlock()

	o.wg.Wait()

	// Close outputs
	for _, output := range o.outputs {
		if err := output.Close(); err != nil {
			return err
		}
	}

	return nil
}

// log logs a message with the given level
func (o *Observer) log(ctx *Context, level Level, msg string, args ...interface{}) {
	// Create entry
	entry := &Entry{
		Time:      time.Now(),
		Level:     level,
		Message:   msg,
		TraceID:   ctx.TraceID(),
		SpanID:    ctx.SpanID(),
		RequestID: ctx.RequestID(),
		Data:      make(map[string]interface{}),
	}

	// Add data
	for i := 0; i < len(args)-1; i += 2 {
		if key, ok := args[i].(string); ok {
			entry.Data[key] = args[i+1]
		}
	}

	// Send entry to buffer
	select {
	case o.buffer <- entry:
		fmt.Printf("Observer: entry sent to buffer: %s\n", entry.Message)
	default:
		fmt.Printf("Observer: buffer full, dropping entry: %s\n", entry.Message)
	}
}

// worker processes entries from buffer
func (o *Observer) worker(interval time.Duration) {
	defer o.wg.Done()

	fmt.Printf("Observer: worker started with interval %v\n", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	entries := make([]*Entry, 0)

	for {
		select {
		case entry := <-o.buffer:
			fmt.Printf("Observer: received entry from buffer: %s\n", entry.Message)
			entries = append(entries, entry)
		case <-ticker.C:
			if len(entries) > 0 {
				fmt.Printf("Observer: flushing %d entries\n", len(entries))
				o.flush(entries)
				entries = make([]*Entry, 0)
			}
		case <-o.done:
			if len(entries) > 0 {
				fmt.Printf("Observer: final flush of %d entries\n", len(entries))
				o.flush(entries)
			}
			fmt.Println("Observer: worker stopped")
			return
		}
	}
}

// flush flushes entries to outputs
func (o *Observer) flush(entries []*Entry) {
	for _, output := range o.outputs {
		if err := output.Write(entries); err != nil {
			// Handle error
			continue
		}
	}

	for _, output := range o.outputs {
		if err := output.Flush(); err != nil {
			// Handle error
			continue
		}
	}
}
