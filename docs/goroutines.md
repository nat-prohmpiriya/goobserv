# Understanding Goroutines in Goobserv

Goroutines are a fundamental concept in Go that enables concurrent execution. This document explains how goroutines are used in the Goobserv project and provides best practices for working with them.

## What are Goroutines?

Goroutines are lightweight threads managed by the Go runtime. They are much more efficient than OS threads:
- Memory usage: ~2KB vs several MB for OS threads
- Can create thousands of goroutines concurrently
- Managed by Go's scheduler rather than the OS scheduler

## Using Goroutines in Goobserv

### Observer Worker

The main usage of goroutines in Goobserv is in the observer's worker, which processes log entries asynchronously:

```go
// Start worker
obs.wg.Add(1)
go obs.worker(cfg.FlushInterval)
```

The worker runs in the background and processes entries from a buffer channel:

```go
func (o *Observer) worker(interval time.Duration) {
    defer o.wg.Done()
    
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    entries := make([]*Entry, 0)

    for {
        select {
        case entry := <-o.buffer:
            entries = append(entries, entry)
        case <-ticker.C:
            if len(entries) > 0 {
                o.flush(entries)
                entries = make([]*Entry, 0)
            }
        case <-o.done:
            if len(entries) > 0 {
                o.flush(entries)
            }
            return
        }
    }
}
```

### Synchronization

We use several mechanisms to synchronize goroutines:

1. **WaitGroup**: Used to wait for goroutines to complete
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // do work
}()
wg.Wait()
```

2. **Mutex**: Used to protect shared resources
```go
var mu sync.Mutex
mu.Lock()
// access shared resource
mu.Unlock()
```

3. **Channels**: Used for communication between goroutines
```go
buffer := make(chan *Entry, bufferSize)
done := make(chan struct{})
```

## Best Practices

1. **Resource Management**
   - Always clean up resources (close channels, stop tickers)
   - Use `defer` for cleanup operations
   - Implement proper shutdown mechanisms

2. **Error Handling**
   - Handle panics in long-running goroutines
   - Log errors appropriately
   - Consider using error channels for error propagation

3. **Race Conditions**
   - Use mutexes to protect shared resources
   - Consider using atomic operations for simple counters
   - Run tests with race detector: `go test -race`

4. **Memory Management**
   - Be careful with goroutine leaks
   - Close channels when they're no longer needed
   - Use buffered channels appropriately

## Common Pitfalls

1. **Closing Channels**
   - Only close channels from the sender side
   - Check if channels are already closed before closing
   - Use sync.Once for safe channel closing

2. **Deadlocks**
   - Avoid circular dependencies between goroutines
   - Be careful with mutex ordering
   - Use timeouts where appropriate

3. **Memory Leaks**
   - Always provide a way to stop goroutines
   - Clean up resources properly
   - Use context for cancellation

## Example: Safe Channel Closing

```go
type Observer struct {
    done   chan struct{}
    closed bool
    mu     sync.Mutex
}

func (o *Observer) Close() error {
    o.mu.Lock()
    if o.closed {
        o.mu.Unlock()
        return nil
    }
    o.closed = true
    close(o.done)
    o.mu.Unlock()
    
    // Wait for worker to finish
    o.wg.Wait()
    return nil
}
```

## Testing

When testing code with goroutines:

1. Use race detector:
```bash
go test -race ./...
```

2. Add appropriate delays in tests:
```go
// Wait for logs to be processed
time.Sleep(100 * time.Millisecond)
```

3. Use test timeouts:
```go
func TestWithTimeout(t *testing.T) {
    done := make(chan bool)
    go func() {
        // do test
        done <- true
    }()
    
    select {
    case <-done:
        // test passed
    case <-time.After(1 * time.Second):
        t.Fatal("test timed out")
    }
}
```

## Further Reading

- [Go Documentation - Goroutines](https://golang.org/doc/effective_go#goroutines)
- [Go Blog - Race Detector](https://blog.golang.org/race-detector)
- [Go Blog - Pipelines and Cancellation](https://blog.golang.org/pipelines)
