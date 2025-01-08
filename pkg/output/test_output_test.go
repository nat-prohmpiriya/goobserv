package output

import (
	"fmt"
	"testing"
	"time"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestTestOutput(t *testing.T) {
	t.Run("Basic operations", func(t *testing.T) {
		output := NewTestOutput()

		// Test initial state
		assert.False(t, output.HasEntries(), "Should not have entries initially")
		assert.Nil(t, output.LastEntry(), "Last entry should be nil initially")
		assert.Empty(t, output.Entries(), "Should not have any entries initially")

		total, kept, dropped := output.Stats()
		assert.Equal(t, 0, total, "Total count should be 0")
		assert.Equal(t, 0, kept, "Kept count should be 0")
		assert.Equal(t, 0, dropped, "Dropped count should be 0")

		// Create test entries
		entries := []*core.Entry{
			{
				Time:      time.Now(),
				Level:     core.LevelInfo,
				Message:   "Test message 1",
				TraceID:   "trace-1",
				SpanID:    "span-1",
				RequestID: "req-1",
				Data: map[string]interface{}{
					"key1": "value1",
				},
			},
			{
				Time:      time.Now(),
				Level:     core.LevelError,
				Message:   "Test message 2",
				TraceID:   "trace-2",
				SpanID:    "span-2",
				RequestID: "req-2",
				Data: map[string]interface{}{
					"key2": "value2",
				},
			},
		}

		// Test Write
		err := output.Write(entries)
		assert.NoError(t, err, "Write should not return error")
		assert.True(t, output.HasEntries(), "Should have entries after write")
		assert.Equal(t, entries[1], output.LastEntry(), "Last entry should match")
		assert.Equal(t, entries, output.Entries(), "All entries should match")

		total, kept, dropped = output.Stats()
		assert.Equal(t, 2, total, "Total count should be 2")
		assert.Equal(t, 2, kept, "Kept count should be 2")
		assert.Equal(t, 0, dropped, "Dropped count should be 0")

		// Test Reset
		output.Reset()
		assert.False(t, output.HasEntries(), "Should not have entries after reset")
		assert.Nil(t, output.LastEntry(), "Last entry should be nil after reset")
		assert.Empty(t, output.Entries(), "Should not have any entries after reset")

		total, kept, dropped = output.Stats()
		assert.Equal(t, 0, total, "Total count should be 0 after reset")
		assert.Equal(t, 0, kept, "Kept count should be 0 after reset")
		assert.Equal(t, 0, dropped, "Dropped count should be 0 after reset")
	})

	t.Run("Max entries", func(t *testing.T) {
		maxEntries := 2
		output := NewTestOutput(WithMaxEntries(maxEntries))

		// Create test entries
		entries := make([]*core.Entry, 5)
		for i := 0; i < 5; i++ {
			entries[i] = &core.Entry{
				Time:    time.Now(),
				Level:   core.LevelInfo,
				Message: fmt.Sprintf("Test message %d", i+1),
			}
		}

		// Write entries
		err := output.Write(entries)
		assert.NoError(t, err, "Write should not return error")

		// Check stats
		total, kept, dropped := output.Stats()
		assert.Equal(t, 5, total, "Total count should be 5")
		assert.Equal(t, maxEntries, kept, "Kept count should match maxEntries")
		assert.Equal(t, 3, dropped, "Should have dropped 3 entries")

		// Check kept entries
		assert.Equal(t, maxEntries, len(output.Entries()), "Should have maxEntries entries")
		assert.Equal(t, entries[4], output.LastEntry(), "Last entry should be the latest")
	})

	t.Run("Concurrent access", func(t *testing.T) {
		output := NewTestOutput()
		entries := []*core.Entry{{
			Time:    time.Now(),
			Level:   core.LevelInfo,
			Message: "Test message",
		}}

		// Run concurrent operations
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					output.Write(entries)
					output.HasEntries()
					output.LastEntry()
					output.Entries()
					output.Stats()
				}
				done <- true
			}()
		}

		// Run reset concurrently
		go func() {
			for i := 0; i < 50; i++ {
				output.Reset()
				time.Sleep(time.Millisecond)
			}
			done <- true
		}()

		// Wait for all goroutines
		for i := 0; i < 11; i++ {
			<-done
		}
	})
}
