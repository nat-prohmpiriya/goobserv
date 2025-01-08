package fiber

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

// Config represents middleware configuration
type Config struct {
	Observer  *core.Observer
	SkipPaths []string
}

// Middleware returns a fiber middleware handler
func Middleware(cfg Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip paths
		path := path.Clean(c.Path())
		for _, p := range cfg.SkipPaths {
			if p == path {
				c.Next()
				return nil
			}
		}

		// Extract trace ID and request ID
		traceID := c.Get("X-Trace-ID")
		requestID := c.Get("X-Request-ID")

		// Create context
		ctx := context.WithValue(c.Context(), "request_id", requestID)
		ctx = context.WithValue(ctx, "trace_id", traceID)

		// Set context and observer
		c.Locals("observContext", ctx)
		c.Locals("observer", cfg.Observer)

		// Start span
		span, ctx := cfg.Observer.StartSpan(ctx, fmt.Sprintf("%s %s", c.Method(), path))
		defer func(start time.Time) {
			status := c.Response().StatusCode()
			fmt.Printf("Middleware: status=%d path=%s\n", status, path)

			// Log request with appropriate level based on status code
			if status >= 500 {
				cfg.Observer.Error(ctx, "HTTP Request").
					WithField("method", c.Method()).
					WithField("path", path).
					WithField("status", status).
					WithField("duration_ms", time.Since(start).Milliseconds())
			} else if status >= 400 {
				cfg.Observer.Warn(ctx, "HTTP Request").
					WithField("method", c.Method()).
					WithField("path", path).
					WithField("status", status).
					WithField("duration_ms", time.Since(start).Milliseconds())
			} else {
				cfg.Observer.Info(ctx, "HTTP Request").
					WithField("method", c.Method()).
					WithField("path", path).
					WithField("status", status).
					WithField("duration_ms", time.Since(start).Milliseconds())
			}
			cfg.Observer.EndSpan(span)
		}(time.Now())

		// Process request
		err := c.Next()
		if err != nil {
			fmt.Printf("Middleware: error from Next: %v\n", err)
			return err
		}
		return nil
	}
}

// GetContext returns the observer context from fiber context
func GetContext(c *fiber.Ctx) context.Context {
	if ctx := c.Locals("observContext"); ctx != nil {
		return ctx.(context.Context)
	}
	return context.Background()
}

// GetObserver returns the observer from fiber context
func GetObserver(c *fiber.Ctx) *core.Observer {
	if obs := c.Locals("observer"); obs != nil {
		return obs.(*core.Observer)
	}
	return nil
}
