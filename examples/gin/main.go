package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	middleware "github.com/nat-prohmpiriya/goobserv/pkg/middleware/gin"
	"github.com/nat-prohmpiriya/goobserv/pkg/output"
)

func main() {
	// Create observer
	obs := core.NewObserver(core.Config{
		BufferSize:    1000,
		FlushInterval: 1 * time.Second,
	})
	defer obs.Close()

	// Add stdout output with color
	stdout := output.NewStdoutOutput(output.StdoutConfig{
		Colored: true,
	})
	obs.AddOutput(stdout)

	// Create gin engine
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Middleware(middleware.Config{
		Observer: obs,
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
	}))

	// Add routes
	r.GET("/hello", func(c *gin.Context) {
		// Get observer context
		ctx := middleware.GetContext(c)

		// Start span
		span, newCtx := obs.StartSpan(ctx, "process_hello")
		defer obs.EndSpan(span)

		// Log info with new context
		obs.Info(newCtx, "Processing hello request")

		// Simulate work
		time.Sleep(100 * time.Millisecond)

		// Return response
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/error", func(c *gin.Context) {
		// Get observer context
		ctx := middleware.GetContext(c)

		// Start span
		span, newCtx := obs.StartSpan(ctx, "process_error")
		defer obs.EndSpan(span)

		// Log error with new context
		obs.Error(newCtx, "Something went wrong").
			WithError(fmt.Errorf("test error"))

		// Return error response
		c.JSON(500, gin.H{
			"error": "Something went wrong",
		})
	})

	// Add health check
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	// Run server
	r.Run(":8080")
}
