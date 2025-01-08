package output

import (
	"sync"

	"github.com/vongga-platform/goobserv/pkg/core"
)

// TestOutput is an output handler for testing
type TestOutput struct {
	entries []*core.Entry
	mu      sync.RWMutex
}

// Write writes an entry to the test output
func (o *TestOutput) Write(entry *core.Entry) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.entries = append(o.entries, entry)
	return nil
}

// Flush flushes the test output
func (o *TestOutput) Flush() error {
	return nil
}

// Close closes the test output
func (o *TestOutput) Close() error {
	return nil
}

// Reset resets the test output
func (o *TestOutput) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.entries = nil
}

// HasEntries returns true if the output has entries
func (o *TestOutput) HasEntries() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.entries) > 0
}

// LastEntry returns the last entry
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
	return o.entries
}
