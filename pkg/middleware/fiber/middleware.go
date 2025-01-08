package fiber

import (
	"fmt"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vongga-platform/goobserv/pkg/core"
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
				return c.Next()
			}
		}

		// Extract trace ID and request ID
		traceID := c.Get("X-Trace-ID")
		requestID := c.Get("X-Request-ID")

		// Create context
		ctx := core.NewContext(c.Context()).
			WithTraceID(traceID).
			WithRequestID(requestID)

		// Set context and observer
		c.Locals("observContext", ctx)
		c.Locals("observer", cfg.Observer)

		// Start span
		spanCtx := cfg.Observer.StartSpan(ctx, fmt.Sprintf("%s %s", c.Method(), path))
		defer func(start time.Time) {
			// Log request
			cfg.Observer.Info(spanCtx, "HTTP Request",
				"method", c.Method(),
				"path", path,
				"status", c.Response().StatusCode(),
				"duration_ms", time.Since(start).Milliseconds(),
			)
		}(time.Now())

		// Process request
		return c.Next()
	}
}

// GetContext returns the observer context from fiber context
func GetContext(c *fiber.Ctx) *core.Context {
	if ctx := c.Locals("observContext"); ctx != nil {
		if obsCtx, ok := ctx.(*core.Context); ok {
			return obsCtx
		}
	}
	return core.NewContext(c.Context())
}

// GetObserver returns the observer from fiber context
func GetObserver(c *fiber.Ctx) *core.Observer {
	if obs := c.Locals("observer"); obs != nil {
		if observer, ok := obs.(*core.Observer); ok {
			return observer
		}
	}
	return nil
}
