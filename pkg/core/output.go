package core

// Output represents output handler interface
type Output interface {
	// Write writes entries to output
	Write(entries []*Entry) error

	// Flush flushes the output
	Flush() error

	// Close closes the output
	Close() error
}
