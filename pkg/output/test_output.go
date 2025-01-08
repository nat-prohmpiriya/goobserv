package output

import (
	"sync"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

const (
	defaultMaxEntries = 1000  // Default maximum number of entries to keep
)

// TestOutput represents test output handler
type TestOutput struct {
	entries     []*core.Entry
	mu          sync.RWMutex
	maxEntries  int  // Maximum number of entries to keep
	totalCount  int  // Total number of entries received
	droppedCount int // Number of entries dropped
}

// NewTestOutput creates a new test output handler
func NewTestOutput(opts ...TestOutputOption) *TestOutput {
	o := &TestOutput{
		maxEntries: defaultMaxEntries,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// TestOutputOption represents an option for test output
type TestOutputOption func(*TestOutput)

// WithMaxEntries sets maximum number of entries to keep
func WithMaxEntries(max int) TestOutputOption {
	return func(o *TestOutput) {
		if max > 0 {
			o.maxEntries = max
		}
	}
}

// Write writes entries to test output
func (o *TestOutput) Write(entries []*core.Entry) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Update counters
	o.totalCount += len(entries)

	// If we need to drop entries
	if len(o.entries) + len(entries) > o.maxEntries {
		// Calculate how many entries to keep from the new batch
		spaceLeft := o.maxEntries - len(o.entries)
		if spaceLeft > 0 {
			// Add as many entries as we can from the end
			startIdx := len(entries) - spaceLeft
			o.entries = append(o.entries, entries[startIdx:]...)
			o.droppedCount += startIdx
		} else {
			// Remove oldest entries to make space for newest entries
			numToKeep := len(entries)
			if numToKeep > o.maxEntries {
				numToKeep = o.maxEntries
			}
			startIdx := len(entries) - numToKeep
			o.entries = append(o.entries[:0], entries[startIdx:]...)
			o.droppedCount += len(o.entries) + startIdx
		}
	} else {
		// Add all entries
		o.entries = append(o.entries, entries...)
	}

	return nil
}

// Flush flushes the output
func (o *TestOutput) Flush() error {
	return nil
}

// Close closes the output
func (o *TestOutput) Close() error {
	return nil
}

// Reset resets test output
func (o *TestOutput) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.entries = nil
	o.totalCount = 0
	o.droppedCount = 0
}

// HasEntries returns true if test output has entries
func (o *TestOutput) HasEntries() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.entries) > 0
}

// LastEntry returns last entry
func (o *TestOutput) LastEntry() *core.Entry {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if len(o.entries) == 0 {
		return nil
	}
	return o.entries[len(o.entries)-1]
}

// Entries returns all entries
func (o *TestOutput) Entries() []*core.Entry {
	o.mu.RLock()
	defer o.mu.RUnlock()
	result := make([]*core.Entry, len(o.entries))
	copy(result, o.entries)
	return result
}

// Stats returns output statistics
func (o *TestOutput) Stats() (total, kept, dropped int) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.totalCount, len(o.entries), o.droppedCount
}
