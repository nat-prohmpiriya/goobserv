package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	fibermw "github.com/nat-prohmpiriya/goobserv/pkg/middleware/fiber"
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

	// Create fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Get observer context
			ctx := fibermw.GetContext(c)

			// Start span
			span, newCtx := obs.StartSpan(ctx, "handle_error")
			defer obs.EndSpan(span)

			// Log error
			obs.Error(newCtx, "Request error").
				WithError(err)

			// Return error response
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(fibermw.Middleware(fibermw.Config{
		Observer: obs,
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
	}))

	// Add routes
	app.Get("/hello", func(c *fiber.Ctx) error {
		// Get observer context
		ctx := fibermw.GetContext(c)

		// Start span
		span, newCtx := obs.StartSpan(ctx, "process_hello")
		defer obs.EndSpan(span)

		// Log info
		obs.Info(newCtx, "Processing hello request")

		// Simulate work
		time.Sleep(100 * time.Millisecond)

		// Return response
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		// Get observer context
		ctx := fibermw.GetContext(c)

		// Start span
		span, newCtx := obs.StartSpan(ctx, "process_error")
		defer obs.EndSpan(span)

		// Log error
		obs.Error(newCtx, "Something went wrong").
			WithError(fmt.Errorf("test error"))

		// Return error
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	})

	// Add health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Start server
	app.Listen(":8080")
}
