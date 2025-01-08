package gin

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	"github.com/nat-prohmpiriya/goobserv/pkg/output"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	// Switch to test mode
	gin.SetMode(gin.TestMode)

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
				BufferSize:    10000,
				FlushInterval: 100 * time.Millisecond,
			})
			defer obs.Close()

			// Add test output
			obs.AddOutput(testOutput)

			// Reset test output
			testOutput.Reset()

			// Create router with middleware
			r := gin.New()
			r.Use(Middleware(Config{
				Observer: obs,
				SkipPaths: []string{
					"/health",
				},
			}))

			// Add test routes
			r.GET("/test", func(c *gin.Context) {
				ctx := GetContext(c)
				obs.Info(ctx, "Test endpoint").
					WithField("path", "/test")
				c.JSON(200, gin.H{"status": "ok"})
			})

			r.GET("/error", func(c *gin.Context) {
				ctx := GetContext(c)
				obs.Error(ctx, "Test error").
					WithField("path", "/error")
				c.JSON(500, gin.H{"error": "test error"})
			})

			r.GET("/health", func(c *gin.Context) {
				c.Status(200)
			})

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.path, nil)
			req.Header.Set("X-Trace-ID", "test-trace")
			req.Header.Set("X-Request-ID", "test-request")

			// Serve request
			r.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

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

			// Check logs
			if tt.shouldLog {
				assert.True(t, testOutput.HasEntries(), "Should have log entries")
				entry := testOutput.LastEntry()
				assert.NotNil(t, entry, "Should have a log entry")
				if entry != nil {
					assert.Equal(t, "HTTP Request", entry.Message, "Wrong message")
					assert.Equal(t, "GET", entry.Data["method"], "Wrong method")
					assert.Equal(t, tt.path, entry.Data["path"], "Wrong path")
					assert.Equal(t, tt.expectedStatus, entry.Data["status"], "Wrong status")
					assert.NotNil(t, entry.Data["duration_ms"], "Missing duration")
				}
			} else {
				assert.False(t, testOutput.HasEntries(), "Should not have log entries")
			}
		})
	}
}

func TestGetContext(t *testing.T) {
	// Create gin context with request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test without context
	ctx := GetContext(c)
	assert.NotNil(t, ctx)

	// Test with context
	testCtx := context.WithValue(context.Background(), "test", "value")
	c.Set("observContext", testCtx)
	ctx = GetContext(c)
	assert.Equal(t, testCtx, ctx)
}

func TestGetObserver(t *testing.T) {
	// Create gin context with request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test without observer
	obs := GetObserver(c)
	assert.Nil(t, obs)

	// Test with observer
	testObs := core.NewObserver(core.Config{})
	c.Set("observer", testObs)
	obs = GetObserver(c)
	assert.Equal(t, testObs, obs)
}
