package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	"github.com/nat-prohmpiriya/goobserv/pkg/output"
)

func main() {
	// Create observer
	obs := core.NewObserver(core.Config{
		BufferSize:    1000,
		FlushInterval: time.Second,
	})
	defer obs.Close()

	// Add stdout output
	stdout := output.NewStdoutOutput(output.StdoutConfig{
		Colored: true,
	})
	obs.AddOutput(stdout)

	// Create gin engine
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(gin.Middleware(gin.Config{
		Observer: obs,
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
	}))

	// Add routes
	r.GET("/hello", func(c *gin.Context) {
		// Get observer context
		ctx := gin.GetContext(c)

		// Log info
		obs.Info(ctx, "Processing hello request")

		// Start span
		span, ctx := obs.StartSpan(ctx, "process_hello")
		defer obs.EndSpan(span)

		// Simulate work
		time.Sleep(100 * time.Millisecond)

		// Return response
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/error", func(c *gin.Context) {
		// Get observer context
		ctx := gin.GetContext(c)

		// Log error
		obs.Error(ctx, "Something went wrong").
			WithError(fmt.Errorf("test error"))

		// Return error
		c.JSON(500, gin.H{
			"error": "Internal Server Error",
		})
	})

	// Add health check
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	// Run server
	r.Run(":8080")
}
