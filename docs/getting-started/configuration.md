# การตั้งค่า Goobserv

## การตั้งค่าพื้นฐาน

```go
type Config struct {
    // ข้อมูล Service
    ServiceName    string            // ชื่อ service
    ServiceVersion string            // version
    Environment    string            // dev/staging/prod
    Labels         map[string]string // custom labels
    
    // การจัดเก็บข้อมูล
    StoragePath    string        // ที่เก็บ log
    RetentionDays  int           // เก็บกี่วัน
    MaxFileSize    int64         // ขนาดไฟล์ max
    
    // Performance
    BufferSize     int           // buffer size
    FlushInterval  time.Duration // flush ทุกกี่วินาที
    SamplingRate   float64       // อัตราการเก็บ sample
}
```

### ตัวอย่างการตั้งค่า

```go
obs := goobserv.New(goobserv.Config{
    // ข้อมูลพื้นฐาน
    ServiceName: "user-service",
    ServiceVersion: "1.0.0",
    Environment: "production",
    Labels: map[string]string{
        "team": "backend",
        "region": "asia",
    },
    
    // การจัดเก็บ
    StoragePath: "/var/log/goobserv",
    RetentionDays: 7,
    MaxFileSize: 100 * 1024 * 1024, // 100MB
    
    // Performance
    BufferSize: 1000,
    FlushInterval: 5 * time.Second,
    SamplingRate: 1.0, // เก็บทุก request
})
```

## การตั้งค่าขั้นสูง

### 1. Log Configuration

```go
type LogConfig struct {
    // ระดับของ Log
    Level string // debug, info, warn, error
    
    // Format
    Format string // json, text
    
    // Fields ที่ต้องการ mask
    SensitiveKeys []string
    
    // Custom Format
    Formatter func(Entry) string
}

// ตัวอย่าง
obs := goobserv.New(goobserv.Config{
    Log: LogConfig{
        Level: "info",
        Format: "json",
        SensitiveKeys: []string{
            "password",
            "token",
        },
    },
})
```

### 2. Metric Configuration

```go
type MetricConfig struct {
    // Prefix สำหรับ metrics
    Prefix string
    
    // Default Labels
    DefaultLabels map[string]string
    
    // Custom Metrics
    CustomMetrics []MetricDefinition
}

// ตัวอย่าง
obs := goobserv.New(goobserv.Config{
    Metric: MetricConfig{
        Prefix: "myapp",
        DefaultLabels: map[string]string{
            "version": "1.0.0",
        },
    },
})
```

### 3. Trace Configuration

```go
type TraceConfig struct {
    // Sampling
    SampleRate float64
    
    // Paths ที่ต้องเก็บเสมอ
    AlwaysSamplePaths []string
    
    // Custom Span Attributes
    DefaultAttributes map[string]string
}

// ตัวอย่าง
obs := goobserv.New(goobserv.Config{
    Trace: TraceConfig{
        SampleRate: 0.1, // เก็บ 10%
        AlwaysSamplePaths: []string{
            "/api/users",
            "/api/orders",
        },
    },
})
```

## การตั้งค่า Storage

### 1. File Storage

```go
type FileStorageConfig struct {
    // การ Rotate
    MaxSize    int64 // bytes
    MaxAge     int   // วัน
    MaxBackups int   // จำนวนไฟล์
    
    // Compression
    Compress bool
    
    // Path Pattern
    PathPattern string // e.g., "app.%Y%m%d.log"
}
```

### 2. Buffer Storage

```go
type BufferConfig struct {
    // ขนาด Buffer
    Size int
    
    // Flush Policy
    FlushInterval time.Duration
    FlushSize     int
    
    // Emergency
    EmergencyPath string
}
```

## Environment Variables

สามารถตั้งค่าผ่าน Environment Variables ได้:

```bash
# Service Info
GOOBSERV_SERVICE_NAME=user-service
GOOBSERV_SERVICE_VERSION=1.0.0
GOOBSERV_ENVIRONMENT=production

# Storage
GOOBSERV_STORAGE_PATH=/var/log/goobserv
GOOBSERV_RETENTION_DAYS=7
GOOBSERV_MAX_FILE_SIZE=104857600

# Performance
GOOBSERV_BUFFER_SIZE=1000
GOOBSERV_FLUSH_INTERVAL=5s
GOOBSERV_SAMPLING_RATE=1.0
```

## ตัวอย่างการตั้งค่าทั้งหมด

```go
func main() {
    obs := goobserv.New(goobserv.Config{
        // ข้อมูล Service
        ServiceName: "user-service",
        ServiceVersion: "1.0.0",
        Environment: "production",
        Labels: map[string]string{
            "team": "backend",
            "region": "asia",
        },
        
        // Storage
        StoragePath: "/var/log/goobserv",
        RetentionDays: 7,
        MaxFileSize: 100 * 1024 * 1024,
        
        // Performance
        BufferSize: 1000,
        FlushInterval: 5 * time.Second,
        SamplingRate: 1.0,
        
        // Log
        Log: LogConfig{
            Level: "info",
            Format: "json",
            SensitiveKeys: []string{"password", "token"},
        },
        
        // Metric
        Metric: MetricConfig{
            Prefix: "myapp",
            DefaultLabels: map[string]string{
                "version": "1.0.0",
            },
        },
        
        // Trace
        Trace: TraceConfig{
            SampleRate: 0.1,
            AlwaysSamplePaths: []string{
                "/api/users",
                "/api/orders",
            },
        },
        
        // File Storage
        FileStorage: FileStorageConfig{
            MaxSize: 100 * 1024 * 1024,
            MaxAge: 7,
            MaxBackups: 5,
            Compress: true,
            PathPattern: "app.%Y%m%d.log",
        },
        
        // Buffer
        Buffer: BufferConfig{
            Size: 1000,
            FlushInterval: 5 * time.Second,
            FlushSize: 100,
            EmergencyPath: "/var/log/goobserv/emergency",
        },
    })
    defer obs.Shutdown()
    
    // เริ่มใช้งาน
    app := fiber.New()
    app.Use(goobserv.Middleware())
    app.Listen(":3000")
}
