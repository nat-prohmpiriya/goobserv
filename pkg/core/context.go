package core

import (
	"context"
	"time"
)

// Context represents an observability context with tracing information
type Context struct {
	ctx       context.Context
	traceID   string
	spanID    string
	parentID  string
	startTime time.Time
	attrs     map[string]interface{}
	requestID string
}

// NewContext creates a new observability context
func NewContext(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	
	return &Context{
		ctx:       ctx,
		startTime: time.Now(),
		attrs:     make(map[string]interface{}),
	}
}

// WithTraceID sets the trace ID
func (c *Context) WithTraceID(traceID string) *Context {
	c.traceID = traceID
	return c
}

// WithSpanID sets the span ID
func (c *Context) WithSpanID(spanID string) *Context {
	c.spanID = spanID
	return c
}

// WithParentID sets the parent span ID
func (c *Context) WithParentID(parentID string) *Context {
	c.parentID = parentID
	return c
}

// WithAttribute adds an attribute to the context
func (c *Context) WithAttribute(key string, value interface{}) *Context {
	c.attrs[key] = value
	return c
}

// WithRequestID adds request ID to context
func (c *Context) WithRequestID(requestID string) *Context {
	c.requestID = requestID
	return c
}

// TraceID returns the trace ID
func (c *Context) TraceID() string {
	return c.traceID
}

// SpanID returns the span ID
func (c *Context) SpanID() string {
	return c.spanID
}

// ParentID returns the parent span ID
func (c *Context) ParentID() string {
	return c.parentID
}

// Attribute returns an attribute value
func (c *Context) Attribute(key string) interface{} {
	return c.attrs[key]
}

// Attributes returns all attributes
func (c *Context) Attributes() map[string]interface{} {
	return c.attrs
}

// RequestID returns request ID
func (c *Context) RequestID() string {
	return c.requestID
}

// StartTime returns the context creation time
func (c *Context) StartTime() time.Time {
	return c.startTime
}

// Duration returns the duration since context creation
func (c *Context) Duration() time.Duration {
	return time.Since(c.startTime)
}

// Deadline implements context.Context
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done implements context.Context
func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err implements context.Context
func (c *Context) Err() error {
	return c.ctx.Err()
}

// Value implements context.Context
func (c *Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}
