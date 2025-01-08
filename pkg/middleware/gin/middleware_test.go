package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vongga-platform/goobserv/pkg/core"
	"github.com/vongga-platform/goobserv/pkg/output"
)

func TestMiddleware(t *testing.T) {
	// Switch to test mode
	gin.SetMode(gin.TestMode)

	// Create observer
	obs := core.NewObserver(core.Config{
		BufferSize:    1000,
		FlushInterval: time.Second,
	})
	defer obs.Close()

	// Create test output
	testOutput := &output.TestOutput{}
	obs.AddOutput(testOutput)

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
		obs.Info(ctx, "Test endpoint")
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/error", func(c *gin.Context) {
		ctx := GetContext(c)
		obs.Error(ctx, "Test error")
		c.JSON(500, gin.H{"error": "test error"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
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
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			req.Header.Set("X-Trace-ID", "test-trace")
			req.Header.Set("X-Request-ID", "test-request")

			// Serve request
			r.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

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
	// Create gin context
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test without context
	ctx := GetContext(c)
	assert.NotNil(t, ctx)

	// Test with context
	testCtx := core.NewContext(c.Request.Context())
	c.Set("observContext", testCtx)
	ctx = GetContext(c)
	assert.Equal(t, testCtx, ctx)
}

func TestGetObserver(t *testing.T) {
	// Create gin context
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test without observer
	obs := GetObserver(c)
	assert.Nil(t, obs)

	// Test with observer
	testObs := core.NewObserver(core.Config{})
	c.Set("observer", testObs)
	obs = GetObserver(c)
	assert.Equal(t, testObs, obs)
}
