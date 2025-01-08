package core

import (
	"sync"
	"time"
)

// BufferMetrics represents buffer performance metrics
type BufferMetrics struct {
	entriesAdded    *Counter
	entriesDropped  *Counter
	bufferSize      *Gauge
	flushLatency    *Histogram
	memoryUsage     *Gauge
}

// NewBufferMetrics creates buffer metrics
func NewBufferMetrics(observer *Observer) *BufferMetrics {
	return &BufferMetrics{
		entriesAdded:   observer.Counter("buffer_entries_added_total"),
		entriesDropped: observer.Counter("buffer_entries_dropped_total"),
		bufferSize:     observer.Gauge("buffer_size_current"),
		flushLatency:   observer.Histogram("buffer_flush_duration_ms", []float64{1, 5, 10, 50, 100, 500}),
		memoryUsage:    observer.Gauge("buffer_memory_bytes"),
	}
}

// Pool represents an object pool for entries
type Pool struct {
	entryPool  sync.Pool
	bufferPool sync.Pool
}

// NewPool creates a new object pool
func NewPool() *Pool {
	return &Pool{
		entryPool: sync.Pool{
			New: func() interface{} {
				return &Entry{
					Data: make(map[string]interface{}),
				}
			},
		},
		bufferPool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024)
			},
		},
	}
}

// GetEntry gets an entry from the pool
func (p *Pool) GetEntry() *Entry {
	return p.entryPool.Get().(*Entry)
}

// PutEntry returns an entry to the pool
func (p *Pool) PutEntry(entry *Entry) {
	// Clear entry data
	entry.Message = ""
	entry.Level = 0
	entry.Type = 0
	entry.Context = nil
	for k := range entry.Data {
		delete(entry.Data, k)
	}
	p.entryPool.Put(entry)
}

// GetBuffer gets a buffer from the pool
func (p *Pool) GetBuffer() []byte {
	return p.bufferPool.Get().([]byte)
}

// PutBuffer returns a buffer to the pool
func (p *Pool) PutBuffer(buf []byte) {
	p.bufferPool.Put(buf[:0])
}

// RingBuffer represents a ring buffer for storing entries
type RingBuffer struct {
	entries []*Entry
	size    int
	pos     int
	mu      sync.RWMutex
}

// NewRingBuffer creates a new ring buffer
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		entries: make([]*Entry, size),
		size:    size,
	}
}

// Write writes an entry to the buffer
func (b *RingBuffer) Write(entry *Entry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries[b.pos] = entry
	b.pos = (b.pos + 1) % b.size
}

// ReadAll reads all entries from the buffer
func (b *RingBuffer) ReadAll() []*Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	entries := make([]*Entry, 0, b.size)
	for i := 0; i < b.size; i++ {
		idx := (b.pos - i - 1 + b.size) % b.size
		if b.entries[idx] != nil {
			entries = append(entries, b.entries[idx])
		}
	}

	return entries
}
