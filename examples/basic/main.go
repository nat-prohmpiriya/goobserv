package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	"github.com/nat-prohmpiriya/goobserv/pkg/output"
)

func main() {
	// Create observer with buffer configuration
	obs := core.NewObserver(core.Config{
		BufferSize:     1000,          // Store up to 1000 entries
		FlushInterval:  time.Second,   // Flush every second
	})
	defer obs.Close()

	// Add stdout output
	stdoutConfig := output.StdoutConfig{
		Colored: true,
	}
	stdout := output.NewStdoutOutput(stdoutConfig)
	obs.AddOutput(stdout)

	// Create context
	ctx := context.Background()

	// Create metrics
	requestCounter := obs.Counter("http_requests_total")
	requestDuration := obs.Histogram("request_duration_ms", []float64{10, 50, 100, 500})
	activeRequests := obs.Gauge("active_requests")

	// Get buffer metrics
	bufferSize := obs.Gauge("buffer_size_current")
	bufferMemory := obs.Gauge("buffer_memory_bytes")
	bufferDropped := obs.Counter("buffer_entries_dropped_total")
	bufferLatency := obs.Histogram("buffer_flush_duration_ms", []float64{1, 5, 10, 50, 100, 500})

	// Simulate multiple requests
	for i := 0; i < 5; i++ {
		// Log buffer metrics
		obs.Info(ctx, "Buffer metrics").
			WithField("size", bufferSize.Value()).
			WithField("memory_bytes", bufferMemory.Value()).
			WithField("dropped_entries", bufferDropped.Value())

		// Start root span
		rootSpan, ctx := obs.StartSpan(ctx, fmt.Sprintf("process_request_%d", i))

		// Log start
		obs.Info(ctx, "Processing request started").
			WithField("request_id", i).
			WithField("method", "GET").
			WithField("path", "/api/users")

		// Update metrics
		requestCounter.Inc()
		activeRequests.Set(float64(i + 1))

		// Simulate work
		time.Sleep(50 * time.Millisecond)

		// Start child span
		dbSpan, ctx := obs.StartSpan(ctx, fmt.Sprintf("database_query_%d", i))
		
		// Simulate database work
		time.Sleep(25 * time.Millisecond)
		
		// Simulate error
		if rand.Float64() < 0.3 {
			err := fmt.Errorf("database connection failed for request %d", i)
			obs.Error(ctx, "Database error occurred").
				WithError(err).
				WithField("request_id", i)
			dbSpan.SetStatus(core.SpanStatusError)
		}
		
		obs.EndSpan(dbSpan)

		// Update metrics
		requestDuration.Observe(float64(time.Since(rootSpan.StartTime).Milliseconds()))
		activeRequests.Set(float64(4 - i))

		// Get buffer latency snapshot
		count, sum, buckets := bufferLatency.Snapshot()
		avgLatency := float64(0)
		if count > 0 {
			avgLatency = sum / float64(count)
		}

		// Format buckets for logging
		bucketStats := make(map[string]uint64)
		for bucket, count := range buckets {
			bucketStats[fmt.Sprintf("le_%.0f", bucket)] = count
		}

		// Log completion
		obs.Info(ctx, "Processing request completed").
			WithField("request_id", i).
			WithField("flush_latency_avg_ms", avgLatency).
			WithField("flush_latency_count", count).
			WithField("flush_latency_buckets", bucketStats)

		obs.EndSpan(rootSpan)

		// Let the buffer do its work
		time.Sleep(200 * time.Millisecond)
	}

	// Final flush to ensure all data is written
	if err := obs.Flush(); err != nil {
		panic(err)
	}
}
