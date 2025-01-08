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

// Observer represents an observability handler
type Observer struct {
	buffer   *RingBuffer
	outputs  []Output
	wg       sync.WaitGroup
	stopChan chan struct{}
}

// NewObserver creates a new observer
func NewObserver(cfg Config) *Observer {
	if cfg.BufferSize <= 0 {
		cfg.BufferSize = 1000
	}
	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = time.Second
	}

	obs := &Observer{
		buffer:   NewRingBuffer(cfg.BufferSize),
		outputs:  make([]Output, 0),
		stopChan: make(chan struct{}),
	}

	obs.wg.Add(1)
	go obs.flushLoop(cfg.FlushInterval)

	return obs
}

// AddOutput adds an output handler
func (o *Observer) AddOutput(output Output) {
	o.outputs = append(o.outputs, output)
}

// StartSpan starts a new span
func (o *Observer) StartSpan(ctx context.Context, name string) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return NewContext(ctx)
}

// Info logs an info message
func (o *Observer) Info(ctx context.Context, msg string, args ...interface{}) {
	o.log(ctx, LevelInfo, msg, args...)
}

// Error logs an error message
func (o *Observer) Error(ctx context.Context, msg string, args ...interface{}) {
	o.log(ctx, LevelError, msg, args...)
}

// Debug logs a debug message
func (o *Observer) Debug(ctx context.Context, msg string, args ...interface{}) {
	o.log(ctx, LevelDebug, msg, args...)
}

// Warn logs a warning message
func (o *Observer) Warn(ctx context.Context, msg string, args ...interface{}) {
	o.log(ctx, LevelWarn, msg, args...)
}

// Close closes the observer
func (o *Observer) Close() error {
	close(o.stopChan)
	o.wg.Wait()

	// Flush remaining entries
	o.flush()

	// Close outputs
	for _, output := range o.outputs {
		if err := output.Close(); err != nil {
			return fmt.Errorf("failed to close output: %v", err)
		}
	}

	return nil
}

func (o *Observer) log(ctx context.Context, level Level, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	var obsCtx *Context
	if c, ok := ctx.(*Context); ok {
		obsCtx = c
	} else {
		obsCtx = NewContext(ctx)
	}

	entry := &Entry{
		Time:      time.Now(),
		Level:     level,
		Message:   msg,
		TraceID:   obsCtx.TraceID(),
		SpanID:    obsCtx.SpanID(),
		ParentID:  obsCtx.ParentID(),
		Data:      obsCtx.Attributes(),
		RequestID: obsCtx.RequestID(),
	}

	o.buffer.Write(entry)
}

func (o *Observer) flushLoop(interval time.Duration) {
	defer o.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.flush()
		case <-o.stopChan:
			return
		}
	}
}

func (o *Observer) flush() {
	entries := o.buffer.ReadAll()
	if len(entries) == 0 {
		return
	}

	for _, output := range o.outputs {
		for _, entry := range entries {
			if err := output.Write(entry); err != nil {
				fmt.Printf("Failed to write entry: %v\n", err)
			}
		}
		if err := output.Flush(); err != nil {
			fmt.Printf("Failed to flush output: %v\n", err)
		}
	}
}
