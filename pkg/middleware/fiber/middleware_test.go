package fiber

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/vongga-platform/goobserv/pkg/core"
	"github.com/vongga-platform/goobserv/pkg/output"
)

func TestMiddleware(t *testing.T) {
	// Create observer
	obs := core.NewObserver(core.Config{
		BufferSize:    1000,
		FlushInterval: time.Second,
	})
	defer obs.Close()

	// Create test output
	testOutput := &output.TestOutput{}
	obs.AddOutput(testOutput)

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
		obs.Info(ctx, "Test endpoint")
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		ctx := GetContext(c)
		obs.Error(ctx, "Test error")
		return fiber.NewError(fiber.StatusInternalServerError, "test error")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

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
			// Reset test output
			testOutput.Reset()

			// Create request
			req := httptest.NewRequest("GET", tt.path, nil)
			req.Header.Set("X-Trace-ID", "test-trace")
			req.Header.Set("X-Request-ID", "test-request")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Read response body
			_, err = io.ReadAll(resp.Body)
			assert.NoError(t, err)

			// Check response
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Check logs
			if tt.shouldLog {
				assert.True(t, testOutput.HasEntries())
				entry := testOutput.LastEntry()
				assert.Equal(t, "HTTP Request", entry.Message)
				assert.Equal(t, "GET", entry.Data["method"])
				assert.Equal(t, tt.path, entry.Data["path"])
				assert.Equal(t, tt.expectedStatus, entry.Data["status"])
				assert.Contains(t, entry.Data, "duration_ms")
			} else {
				assert.False(t, testOutput.HasEntries())
			}
		})
	}
}

func TestGetContext(t *testing.T) {
	// Create fiber app and context
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Test without context
	obsCtx := GetContext(ctx)
	assert.NotNil(t, obsCtx)

	// Test with context
	testCtx := core.NewContext(ctx.Context())
	ctx.Locals("observContext", testCtx)
	obsCtx = GetContext(ctx)
	assert.Equal(t, testCtx, obsCtx)
}

func TestGetObserver(t *testing.T) {
	// Create fiber app and context
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Test without observer
	obs := GetObserver(ctx)
	assert.Nil(t, obs)

	// Test with observer
	testObs := core.NewObserver(core.Config{})
	ctx.Locals("observer", testObs)
	obs = GetObserver(ctx)
	assert.Equal(t, testObs, obs)
}
