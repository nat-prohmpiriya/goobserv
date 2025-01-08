package core

import (
	"sync"
)

// MetricType represents the type of metric
type MetricType int

const (
	CounterMetric MetricType = iota
	GaugeMetric
	HistogramMetric
)

// Metric represents a metric value
type Metric struct {
	Name   string
	Type   MetricType
	Value  float64
	Labels map[string]string
}

// Counter represents a cumulative metric
type Counter struct {
	name   string
	value  float64
	labels map[string]string
	mu     sync.RWMutex
}

// NewCounter creates a new counter
func NewCounter(name string) *Counter {
	return &Counter{
		name:   name,
		labels: make(map[string]string),
	}
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	c.Add(1)
}

// Add adds a value to the counter
func (c *Counter) Add(value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += value
}

// WithLabels adds labels to the counter
func (c *Counter) WithLabels(labels map[string]string) *Counter {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range labels {
		c.labels[k] = v
	}
	return c
}

// Value returns the current value
func (c *Counter) Value() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Gauge represents a metric that can go up and down
type Gauge struct {
	name   string
	value  float64
	labels map[string]string
	mu     sync.RWMutex
}

// NewGauge creates a new gauge
func NewGauge(name string) *Gauge {
	return &Gauge{
		name:   name,
		labels: make(map[string]string),
	}
}

// Set sets the gauge value
func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
}

// WithLabels adds labels to the gauge
func (g *Gauge) WithLabels(labels map[string]string) *Gauge {
	g.mu.Lock()
	defer g.mu.Unlock()
	for k, v := range labels {
		g.labels[k] = v
	}
	return g
}

// Value returns the current value
func (g *Gauge) Value() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

// Histogram represents a metric that samples observations
type Histogram struct {
	name      string
	buckets   []float64
	counts    map[float64]uint64
	sum       float64
	count     uint64
	labels    map[string]string
	mu        sync.RWMutex
}

// NewHistogram creates a new histogram
func NewHistogram(name string, buckets []float64) *Histogram {
	return &Histogram{
		name:      name,
		buckets:   buckets,
		counts:    make(map[float64]uint64),
		labels:    make(map[string]string),
	}
}

// Observe adds a observation
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sum += value
	h.count++

	for _, bucket := range h.buckets {
		if value <= bucket {
			h.counts[bucket]++
		}
	}
}

// WithLabels adds labels to the histogram
func (h *Histogram) WithLabels(labels map[string]string) *Histogram {
	h.mu.Lock()
	defer h.mu.Unlock()
	for k, v := range labels {
		h.labels[k] = v
	}
	return h
}

// Snapshot returns the current values
func (h *Histogram) Snapshot() (uint64, float64, map[float64]uint64) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.count, h.sum, h.counts
}
