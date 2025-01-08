# Data Management

## Sampling Strategies

### 1. Log Sampling
```go
obs := goobserv.New(goobserv.Config{
    Sampling: &SamplingConfig{
        // เก็บตาม Log Level
        Error: 1.0,    // เก็บ 100%
        Warning: 0.5,  // เก็บ 50%
        Info: 0.1,     // เก็บ 10%
        Debug: 0.01,   // เก็บ 1%

        // เก็บตาม Service Priority
        HighPriority: 1.0,   // เก็บทุก log
        MediumPriority: 0.5, // เก็บครึ่งนึง
        LowPriority: 0.1,    // เก็บ 10%
    },
})
```

### 2. Metric Sampling
```go
obs := goobserv.New(goobserv.Config{
    Metrics: &MetricsConfig{
        // Include สำคัญ
        Include: []string{
            "http_requests_total",
            "request_duration_ms",
            "error_count",
            "cpu_usage",
            "memory_usage",
        },

        // Exclude ไม่จำเป็น
        Exclude: []string{
            "debug_*",
            "internal_*",
            "test_*",
        },

        // Aggregate ตาม Interval
        Interval: map[string]time.Duration{
            "high_resolution": time.Minute,
            "medium_resolution": 5 * time.Minute,
            "low_resolution": 15 * time.Minute,
        },
    },
})
```

### 3. Trace Sampling
```go
obs := goobserv.New(goobserv.Config{
    Tracing: &TracingConfig{
        // เก็บตาม Condition
        Error: 1.0,                           // Error traces
        Slow: 1.0,                           // Slow traces
        Normal: 0.01,                        // Normal traces

        // Thresholds
        SlowRequestThreshold: time.Second,    // > 1s = slow
        ErrorPattern: []string{               // Error patterns
            "connection refused",
            "timeout",
            "5xx",
        },
    },
})
```

## Retention Policies

### 1. Log Retention
```go
obs := goobserv.New(goobserv.Config{
    Retention: &RetentionConfig{
        Logs: map[string]time.Duration{
            // ตาม Priority
            "error": 30 * 24 * time.Hour,    // 30 วัน
            "warn": 7 * 24 * time.Hour,      // 7 วัน
            "info": 2 * 24 * time.Hour,      // 2 วัน
            "debug": 6 * time.Hour,          // 6 ชั่วโมง

            // ตาม Service
            "payment": 90 * 24 * time.Hour,  // 90 วัน
            "user": 30 * 24 * time.Hour,     // 30 วัน
            "default": 7 * 24 * time.Hour,   // 7 วัน
        },
    },
})
```

### 2. Metric Retention
```go
obs := goobserv.New(goobserv.Config{
    Retention: &RetentionConfig{
        Metrics: map[string]RetentionRule{
            "raw": {
                Duration: 7 * 24 * time.Hour,    // เก็บ raw data 7 วัน
                Resolution: time.Minute,          // resolution 1 นาที
            },
            "hourly": {
                Duration: 30 * 24 * time.Hour,   // เก็บ hourly data 30 วัน
                Resolution: time.Hour,            // resolution 1 ชั่วโมง
            },
            "daily": {
                Duration: 365 * 24 * time.Hour,  // เก็บ daily data 1 ปี
                Resolution: 24 * time.Hour,       // resolution 1 วัน
            },
        },
    },
})
```

### 3. Trace Retention
```go
obs := goobserv.New(goobserv.Config{
    Retention: &RetentionConfig{
        Traces: map[string]time.Duration{
            "error": 7 * 24 * time.Hour,     // error traces 7 วัน
            "slow": 3 * 24 * time.Hour,      // slow traces 3 วัน
            "normal": 24 * time.Hour,        // normal traces 1 วัน
        },
    },
})
```

## Storage Optimization

### 1. Compression
```go
obs := goobserv.New(goobserv.Config{
    Storage: &StorageConfig{
        Compression: &CompressionConfig{
            // Algorithm
            Type: "zstd",      // zstd, gzip, lz4
            Level: 3,          // 1-9

            // Threshold
            MinSize: 1024,     // compress if > 1KB
            
            // Types
            CompressLogs: true,
            CompressMetrics: true,
            CompressTraces: true,
        },
    },
})
```

### 2. Aggregation
```go
obs := goobserv.New(goobserv.Config{
    Storage: &StorageConfig{
        Aggregation: &AggregationConfig{
            // Metrics
            Metrics: map[string]AggregationRule{
                "request_duration": {
                    Functions: []string{"avg", "p95", "p99"},
                    Interval: time.Hour,
                },
                "error_count": {
                    Functions: []string{"sum"},
                    Interval: time.Hour,
                },
            },

            // Traces
            Traces: map[string]AggregationRule{
                "duration": {
                    Functions: []string{"avg", "max"},
                    Interval: time.Hour,
                },
            },
        },
    },
})
```

### 3. Indexing
```go
obs := goobserv.New(goobserv.Config{
    Storage: &StorageConfig{
        Indexes: &IndexConfig{
            // Logs
            Logs: []string{
                "timestamp",
                "level",
                "service",
                "trace_id",
            },

            // Metrics
            Metrics: []string{
                "timestamp",
                "name",
                "service",
            },

            // Traces
            Traces: []string{
                "trace_id",
                "service",
                "duration",
            },
        },
    },
})
```

## Cost Management

### 1. Storage Costs
```go
// 1. Data Volume
- ใช้ sampling ลดปริมาณข้อมูล
- compress ข้อมูล
- aggregate ข้อมูลเก่า

// 2. Retention
- กำหนด retention ตาม priority
- ลบข้อมูลที่ไม่จำเป็น
- archive ข้อมูลเก่า

// 3. Storage Type
- ใช้ storage ที่เหมาะสม
- แยก hot/cold storage
- ใช้ object storage สำหรับข้อมูลเก่า
```

### 2. Processing Costs
```go
// 1. CPU Usage
- batch processing
- async processing
- ใช้ efficient algorithms

// 2. Memory Usage
- ใช้ buffer pool
- limit buffer size
- garbage collection tuning

// 3. Network Usage
- compress ก่อนส่ง
- batch requests
- local aggregation
```

## Best Practices

### 1. Environment-based Configuration
```go
switch env := os.Getenv("ENV"); env {
case "production":
    // เน้น reliability และ cost
    - เก็บ error ทั้งหมด
    - sample warning/info
    - retention สั้น
    
case "staging":
    // เน้น testing และ monitoring
    - เก็บ error/warning ทั้งหมด
    - sample info/debug
    - retention ปานกลาง
    
case "development":
    // เน้น debugging
    - เก็บ logs ทั้งหมด
    - retention สั้น
    - ไม่ต้อง optimize มาก
}
```

### 2. Dynamic Configuration
```go
// 1. Traffic-based
if isHighTraffic() {
    decreaseSamplingRate()
    increaseBufferSize()
    enableBatchProcessing()
}

// 2. Resource-based
if isHighMemory() {
    flushBuffer()
    enableCompression()
    decreaseRetention()
}

// 3. Cost-based
if isHighCost() {
    increaseSampling()
    enableAggregation()
    decreaseRetention()
}
```

### 3. Monitoring
```go
// 1. Storage Usage
- monitor disk usage
- monitor data growth
- alert on high usage

// 2. Performance
- monitor latency
- monitor throughput
- monitor resource usage

// 3. Costs
- monitor storage costs
- monitor processing costs
- monitor network costs
```
