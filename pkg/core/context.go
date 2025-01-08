package core

import (
	"context"
	"sync"
)

// Context represents observer context
type Context struct {
	ctx       context.Context
	traceID   string
	spanID    string
	requestID string
	mu        sync.RWMutex
}

// NewContext creates a new context
func NewContext(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Context{
		ctx: ctx,
	}
}

// WithTraceID sets trace ID
func (c *Context) WithTraceID(traceID string) *Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.traceID = traceID
	return c
}

// WithSpan sets span ID
func (c *Context) WithSpan(name string) *Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spanID = name
	return c
}

// WithRequestID sets request ID
func (c *Context) WithRequestID(requestID string) *Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requestID = requestID
	return c
}

// EndSpan ends the current span
func (c *Context) EndSpan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spanID = ""
}

// TraceID returns trace ID
func (c *Context) TraceID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.traceID
}

// SpanID returns span ID
func (c *Context) SpanID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.spanID
}

// RequestID returns request ID
func (c *Context) RequestID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.requestID
}

// Context returns underlying context
func (c *Context) Context() context.Context {
	return c.ctx
}
