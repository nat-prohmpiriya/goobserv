package output

import (
	"sync"

	"github.com/nat-prohmpiriya/goobserv/pkg/core"
)

// TestOutput represents test output handler
type TestOutput struct {
	entries []*core.Entry
	mu      sync.RWMutex
}

// Write writes entries to test output
func (o *TestOutput) Write(entries []*core.Entry) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.entries = append(o.entries, entries...)
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
	return o.entries
}
