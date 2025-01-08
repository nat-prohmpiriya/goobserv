# แนวทางการจัดการ Performance

## Buffer Management

### 1. Buffer Size

```go
// ตั้งค่า buffer size ให้เหมาะสม
obs := goobserv.New(goobserv.Config{
    // คำนวณจาก:
    // - จำนวน requests ต่อวินาที
    // - ขนาดของแต่ละ entry
    // - memory ที่มี
    BufferSize: 1000,
    
    // ระยะเวลาในการ flush
    FlushInterval: 5 * time.Second,
})
```

### 2. Batch Processing

```go
type Buffer struct {
    entries    []Entry
    batchSize  int
    flushChan  chan struct{}
}

func (b *Buffer) process() {
    for {
        if len(b.entries) >= b.batchSize {
            b.flush()
        }
        
        select {
        case entry := <-b.entryChan:
            b.entries = append(b.entries, entry)
        case <-b.flushChan:
            b.flush()
        }
    }
}
```

## Memory Management

### 1. Object Pooling

```go
// ใช้ sync.Pool สำหรับ objects ที่ใช้บ่อย
var entryPool = sync.Pool{
    New: func() interface{} {
        return &Entry{
            Data: make(map[string]interface{}),
        }
    },
}

// การใช้งาน
entry := entryPool.Get().(*Entry)
defer entryPool.Put(entry)
```

### 2. Memory Limits

```go
type MemoryConfig struct {
    // จำกัดขนาด entry
    MaxEntrySize int
    
    // จำกัด memory usage
    MaxMemoryUsage int64
    
    // cleanup เมื่อใช้ memory เยอะ
    MemoryThreshold float64
}
```

## Disk I/O

### 1. File Rotation

```go
type FileConfig struct {
    // ขนาดไฟล์สูงสุด
    MaxFileSize int64
    
    // จำนวนไฟล์สูงสุด
    MaxFiles int
    
    // อายุไฟล์สูงสุด
    MaxAge time.Duration
}

// ตัวอย่าง
obs := goobserv.New(goobserv.Config{
    File: FileConfig{
        MaxFileSize: 100 * 1024 * 1024, // 100MB
        MaxFiles: 5,
        MaxAge: 7 * 24 * time.Hour,     // 7 วัน
    },
})
```

### 2. Compression

```go
type CompressionConfig struct {
    // เปิด/ปิด compression
    Enabled bool
    
    // ระดับการบีบอัด
    Level int
    
    // ประเภทไฟล์ที่จะบีบ
    FilePattern []string
}
```

## Concurrency

### 1. Worker Pool

```go
type WorkerPool struct {
    workers  int
    jobChan  chan Job
    doneChan chan struct{}
}

func NewWorkerPool(workers int) *WorkerPool {
    pool := &WorkerPool{
        workers:  workers,
        jobChan:  make(chan Job),
        doneChan: make(chan struct{}),
    }
    
    for i := 0; i < workers; i++ {
        go pool.worker()
    }
    
    return pool
}
```

### 2. Rate Limiting

```go
type RateLimit struct {
    // จำนวน requests ต่อวินาที
    RPS int
    
    // burst size
    Burst int
    
    // timeout
    Timeout time.Duration
}
```

## Sampling

### 1. Request Sampling

```go
type SamplingConfig struct {
    // อัตราการ sampling
    Rate float64
    
    // paths ที่ต้อง sample เสมอ
    AlwaysSample []string
    
    // paths ที่ไม่ต้อง sample
    NeverSample []string
}
```

### 2. Adaptive Sampling

```go
type AdaptiveSampling struct {
    // อัตราเริ่มต้น
    InitialRate float64
    
    // อัตราต่ำสุด
    MinRate float64
    
    // อัตราสูงสุด
    MaxRate float64
    
    // ปรับตามโหลด
    AdjustByLoad bool
}
```

## Best Practices

### 1. Buffer Configuration

```go
// ✅ ดี
obs := goobserv.New(goobserv.Config{
    BufferSize: 1000,
    FlushInterval: 5 * time.Second,
})

// ❌ ไม่ดี
obs := goobserv.New(goobserv.Config{
    BufferSize: 1000000,    // ใช้ memory เยอะเกินไป
    FlushInterval: 100 * time.Millisecond, // flush บ่อยเกินไป
})
```

### 2. File Management

```go
// ✅ ดี
obs := goobserv.New(goobserv.Config{
    File: FileConfig{
        MaxFileSize: 100 * 1024 * 1024, // 100MB
        MaxFiles: 5,
        MaxAge: 7 * 24 * time.Hour,
    },
})

// ❌ ไม่ดี
obs := goobserv.New(goobserv.Config{
    File: FileConfig{
        MaxFileSize: 1024 * 1024 * 1024, // 1GB
        MaxFiles: 100,                    // เก็บไฟล์เยอะเกินไป
        MaxAge: 365 * 24 * time.Hour,    // เก็บนานเกินไป
    },
})
```

### 3. Memory Usage

```go
// ✅ ดี
entry := entryPool.Get().(*Entry)
defer entryPool.Put(entry)

// ❌ ไม่ดี
entry := &Entry{
    Data: make(map[string]interface{}),
}
```
