package core

import (
	"context"
)

type observerKey struct{}

// Config represents observer configuration
type Config struct {
	Development bool
	BufferSize  int
}

// Observer handles logging and tracing
type Observer struct {
	buffer chan *Entry
	config *Config
}

// NewObserver creates a new observer
func NewObserver(config *Config) *Observer {
	if config == nil {
		config = &Config{
			Development: false,
			BufferSize:  1000,
		}
	}
	return &Observer{
		buffer: make(chan *Entry, config.BufferSize),
		config: config,
	}
}

// Buffer returns the entry buffer channel
func (o *Observer) Buffer() chan *Entry {
	return o.buffer
}

// WithObserver adds observer to context
func WithObserver(ctx context.Context, obs *Observer) context.Context {
	return context.WithValue(ctx, observerKey{}, obs)
}

// GetObserver gets observer from context
func GetObserver(ctx context.Context) *Observer {
	if obs, ok := ctx.Value(observerKey{}).(*Observer); ok {
		return obs
	}
	return nil
}
