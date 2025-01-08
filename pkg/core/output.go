package core

// Output represents an output handler interface
type Output interface {
	Write(entry *Entry) error
	Flush() error
	Close() error
}
