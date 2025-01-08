package fiber

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	"github.com/nat-prohmpiriya/goobserv/pkg/output"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestMiddleware(t *testing.T) {
	// Create test output
	testOutput := &output.TestOutput{}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		shouldLog      bool
	}{
		{
			name:           "Normal request",
			path:           "/test",
			expectedStatus: 200,
			shouldLog:      true,
		},
		{
			name:           "Error request",
			path:           "/error",
			expectedStatus: 500,
			shouldLog:      true,
		},
		{
			name:           "Skip path",
			path:           "/health",
			expectedStatus: 200,
			shouldLog:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create observer with shorter flush interval for tests
			obs := core.NewObserver(core.Config{
				BufferSize:    10000, // เพิ่ม buffer size
				FlushInterval: 100 * time.Millisecond,
			})
			defer obs.Close()

			// Add test output
			obs.AddOutput(testOutput)

			// Reset test output
			testOutput.Reset()

			// Create app with middleware
			app := fiber.New()
			app.Use(Middleware(Config{
				Observer: obs,
				SkipPaths: []string{
					"/health",
				},
			}))

			// Add test routes
			app.Get("/test", func(c *fiber.Ctx) error {
				ctx := GetContext(c)
				obs.Info(ctx, "Test endpoint").
					WithField("path", "/test")
				return c.JSON(fiber.Map{"status": "ok"})
			})

			app.Get("/error", func(c *fiber.Ctx) error {
				ctx := GetContext(c)
				obs.Error(ctx, "Test error").
					WithField("path", "/error")
				return c.Status(500).JSON(fiber.Map{"error": "test error"})
			})

			app.Get("/health", func(c *fiber.Ctx) error {
				return c.SendStatus(200)
			})

			// Create request
			req := httptest.NewRequest("GET", tt.path, nil)
			req.Header.Set("X-Trace-ID", "test-trace")
			req.Header.Set("X-Request-ID", "test-request")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check response
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Wait for logs to be processed (3 flush intervals)
			time.Sleep(300 * time.Millisecond)

			// Force flush before checking logs
			if err := obs.Close(); err != nil {
				t.Fatalf("Failed to close observer: %v", err)
			}

			// Create new observer for next test
			obs = core.NewObserver(core.Config{
				BufferSize:    10000,
				FlushInterval: 100 * time.Millisecond,
			})

			// Debug log
			t.Logf("Test case: %s", tt.name)
			t.Logf("Should log: %v", tt.shouldLog)
			t.Logf("Has entries: %v", testOutput.HasEntries())
			if entries := testOutput.Entries(); len(entries) > 0 {
				for i, e := range entries {
					t.Logf("Entry %d: Message=%s Level=%v", i, e.Message, e.Level)
				}
			} else {
				t.Log("No entries found")
			}

			// Check logs
			if tt.shouldLog {
				assert.True(t, testOutput.HasEntries(), "Should have log entries")
				entries := testOutput.Entries()
				assert.NotEmpty(t, entries, "Should have log entries")

				// Find HTTP Request entry
				var httpEntry *core.Entry
				for _, e := range entries {
					if e.Message == "HTTP Request" {
						httpEntry = e
						break
					}
				}

				assert.NotNil(t, httpEntry, "Should have HTTP Request entry")
				if httpEntry != nil {
					assert.Equal(t, "GET", httpEntry.Data["method"], "Wrong method")
					assert.Equal(t, tt.path, httpEntry.Data["path"], "Wrong path")
					assert.Equal(t, tt.expectedStatus, httpEntry.Data["status"], "Wrong status")
					assert.NotNil(t, httpEntry.Data["duration_ms"], "Missing duration")
				}

				// For error requests, check error message
				if tt.expectedStatus >= 500 {
					var errorEntry *core.Entry
					for _, e := range entries {
						if e.Message == "Test error" {
							errorEntry = e
							break
						}
					}
					assert.NotNil(t, errorEntry, "Should have error entry")
				}
			} else {
				assert.False(t, testOutput.HasEntries(), "Should not have log entries")
			}
		})
	}
}

func TestGetContext(t *testing.T) {
	// Create fiber app
	app := fiber.New()

	// Create context
	ctx := &fasthttp.RequestCtx{}
	c := app.AcquireCtx(ctx)
	defer app.ReleaseCtx(c)

	// Test without context
	obsCtx := GetContext(c)
	assert.NotNil(t, obsCtx)

	// Test with context
	testCtx := context.WithValue(context.Background(), "test", "value")
	c.Locals("observContext", testCtx)
	obsCtx = GetContext(c)
	assert.Equal(t, testCtx, obsCtx)
}

func TestGetObserver(t *testing.T) {
	// Create fiber app
	app := fiber.New()

	// Create context
	ctx := &fasthttp.RequestCtx{}
	c := app.AcquireCtx(ctx)
	defer app.ReleaseCtx(c)

	// Test without observer
	obs := GetObserver(c)
	assert.Nil(t, obs)

	// Test with observer
	testObs := core.NewObserver(core.Config{})
	c.Locals("observer", testObs)
	obs = GetObserver(c)
	assert.Equal(t, testObs, obs)
}
