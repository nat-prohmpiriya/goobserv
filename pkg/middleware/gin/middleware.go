package gin

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

// Config represents middleware configuration
type Config struct {
	Observer  *core.Observer
	SkipPaths []string
}

// Middleware returns a gin middleware handler
func Middleware(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip paths
		path := c.Request.URL.Path
		for _, p := range cfg.SkipPaths {
			if p == path {
				c.Next()
				return
			}
		}

		// Extract trace ID and request ID
		traceID := c.GetHeader("X-Trace-ID")
		requestID := c.GetHeader("X-Request-ID")

		// Create context
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		ctx = context.WithValue(ctx, "trace_id", traceID)

		// Set context and observer
		c.Set("observContext", ctx)
		c.Set("observer", cfg.Observer)

		// Start span
		span, ctx := cfg.Observer.StartSpan(ctx, fmt.Sprintf("%s %s", c.Request.Method, path))
		defer func(start time.Time) {
			// Log request
			status := c.Writer.Status()
			if status >= 500 {
				cfg.Observer.Error(ctx, "HTTP Request").
					WithField("method", c.Request.Method).
					WithField("path", path).
					WithField("status", status).
					WithField("duration_ms", time.Since(start).Milliseconds())
			} else if status >= 400 {
				cfg.Observer.Warn(ctx, "HTTP Request").
					WithField("method", c.Request.Method).
					WithField("path", path).
					WithField("status", status).
					WithField("duration_ms", time.Since(start).Milliseconds())
			} else {
				cfg.Observer.Info(ctx, "HTTP Request").
					WithField("method", c.Request.Method).
					WithField("path", path).
					WithField("status", status).
					WithField("duration_ms", time.Since(start).Milliseconds())
			}
			cfg.Observer.EndSpan(span)
		}(time.Now())

		// Process request
		c.Next()
	}
}

// GetContext returns the observer context from gin context
func GetContext(c *gin.Context) context.Context {
	if ctx, exists := c.Get("observContext"); exists {
		if obsCtx, ok := ctx.(context.Context); ok {
			return obsCtx
		}
	}
	return context.Background()
}

// GetObserver returns the observer from gin context
func GetObserver(c *gin.Context) *core.Observer {
	if obs, exists := c.Get("observer"); exists {
		if observer, ok := obs.(*core.Observer); ok {
			return observer
		}
	}
	return nil
}
