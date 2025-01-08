package gin

import (
	"fmt"
	"path"
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
		path := path.Clean(c.Request.URL.Path)
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
		ctx := core.NewContext(c.Request.Context()).
			WithTraceID(traceID).
			WithRequestID(requestID)

		// Set context and observer
		c.Set("observContext", ctx)
		c.Set("observer", cfg.Observer)

		// Start span
		spanCtx := cfg.Observer.StartSpan(ctx, fmt.Sprintf("%s %s", c.Request.Method, path))
		defer func(start time.Time) {
			// Log request
			cfg.Observer.Info(spanCtx, "HTTP Request",
				"method", c.Request.Method,
				"path", path,
				"status", c.Writer.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
			)
			cfg.Observer.EndSpan(spanCtx)
		}(time.Now())

		// Process request
		c.Next()
	}
}

// GetContext returns the observer context from gin context
func GetContext(c *gin.Context) *core.Context {
	if ctx, exists := c.Get("observContext"); exists {
		if obsCtx, ok := ctx.(*core.Context); ok {
			return obsCtx
		}
	}
	return core.NewContext(c.Request.Context())
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
