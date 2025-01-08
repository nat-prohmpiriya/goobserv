# Outputs

## รูปแบบ Output ที่รองรับ

### 1. Terminal Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "stdout",
            Format: "text",  // "text" หรือ "json"
            ColorEnabled: true,
            PrettyPrint: true,
        },
    },
})
```

### 2. File Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "file",
            Path: "/var/log/app.log",
            Format: "json",
            Rotation: &RotationConfig{
                MaxSize: 100,    // MB
                MaxAge: 7,       // วัน
                MaxBackups: 5,   // จำนวนไฟล์
                Compress: true,
            },
        },
    },
})
```

### 3. MongoDB Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "mongodb",
            URI: "mongodb://localhost:27017",
            Database: "logs",
            Collection: "app_logs",
            BatchSize: 100,
            FlushInterval: 5 * time.Second,
        },
    },
})
```

### 4. Elasticsearch Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "elasticsearch",
            Addresses: []string{"http://localhost:9200"},
            Index: "app-logs",
            BatchSize: 100,
            FlushInterval: 5 * time.Second,
            RetryConfig: &RetryConfig{
                MaxRetries: 3,
                InitialInterval: time.Second,
                MaxInterval: 30 * time.Second,
            },
        },
    },
})
```

### 5. Kafka Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "kafka",
            Brokers: []string{"localhost:9092"},
            Topic: "app-logs",
            BatchSize: 100,
            CompressionType: "snappy",
            RequiredAcks: -1,  // all
        },
    },
})
```

### 6. Redis Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "redis",
            Addr: "localhost:6379",
            Key: "app:logs",
            DataType: "stream",  // "list" หรือ "stream"
            MaxLen: 10000,
            TTL: 24 * time.Hour,
        },
    },
})
```

### 7. HTTP Webhook Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "webhook",
            URL: "https://api.example.com/logs",
            Method: "POST",
            Headers: map[string]string{
                "Authorization": "Bearer token",
                "Content-Type": "application/json",
            },
            RetryConfig: &RetryConfig{
                MaxRetries: 3,
                InitialInterval: time.Second,
            },
        },
    },
})
```

### 8. Custom Output
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "custom",
            Handler: func(entry Entry) error {
                // Custom logic
                return nil
            },
        },
    },
})
```

## การใช้งานหลาย Outputs

### Development Environment
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "stdout",
            Format: "text",
            ColorEnabled: true,
        },
        {
            Type: "file",
            Path: "app.log",
            Format: "json",
        },
    },
})
```

### Production Environment
```go
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "file",
            Path: "/var/log/app.log",
            Format: "json",
            Rotation: &RotationConfig{
                MaxSize: 100,
                MaxAge: 7,
                Compress: true,
            },
        },
        {
            Type: "elasticsearch",
            Addresses: []string{"http://es:9200"},
            Index: "app-logs",
        },
        {
            Type: "kafka",
            Brokers: []string{"kafka:9092"},
            Topic: "app-logs",
        },
    },
})
```

## Best Practices

### 1. Output Selection
```go
// ✅ ดี: เลือก outputs ตามความเหมาะสม
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "file",    // สำหรับ audit และ backup
            Path: "/var/log/app.log",
        },
        {
            Type: "elasticsearch", // สำหรับ search และ analysis
            Addresses: []string{"http://es:9200"},
        },
    },
})

// ❌ ไม่ดี: ใช้ outputs มากเกินความจำเป็น
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        { Type: "file" },
        { Type: "elasticsearch" },
        { Type: "kafka" },
        { Type: "redis" },
        { Type: "webhook" },
    },
})
```

### 2. Error Handling
```go
// ✅ ดี: มี error handling ที่ดี
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "elasticsearch",
            RetryConfig: &RetryConfig{
                MaxRetries: 3,
                InitialInterval: time.Second,
            },
        },
        {
            Type: "file", // fallback
            Path: "/var/log/app.log",
        },
    },
})

// ❌ ไม่ดี: ไม่มี error handling
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "elasticsearch",
            // ไม่มี retry config
        },
    },
})
```

### 3. Performance
```go
// ✅ ดี: ใช้ batch และ async
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "elasticsearch",
            BatchSize: 100,
            FlushInterval: 5 * time.Second,
            Async: true,
        },
    },
})

// ❌ ไม่ดี: sync และ no batching
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "elasticsearch",
            BatchSize: 1,
            Async: false,
        },
    },
})
```

### 4. Security
```go
// ✅ ดี: มีการป้องกันที่ดี
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "elasticsearch",
            TLS: &TLSConfig{
                CertFile: "/path/to/cert",
                KeyFile: "/path/to/key",
            },
            Authentication: &AuthConfig{
                Username: "user",
                Password: "pass",
            },
        },
    },
})

// ❌ ไม่ดี: ไม่มีการป้องกัน
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {
            Type: "elasticsearch",
            // ไม่มี TLS และ authentication
        },
    },
})
```
