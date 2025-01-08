# เปรียบเทียบ Zap vs Goobserv

## 1. การใช้งานพื้นฐาน

### Zap
```go
logger, _ := zap.NewProduction()
defer logger.Sync()

// ต้องระบุ type ของแต่ละ field
logger.Info("user created",
    zap.String("user_id", "123"),
    zap.String("email", "user@example.com"),
    zap.Int("age", 25),
)

// ต้อง clone logger เพื่อเพิ่ม fields
userLogger := logger.With(
    zap.String("user_id", "123"),
)
userLogger.Info("processing user")
```

### Goobserv
```go
obs := goobserv.New(goobserv.Config{})

// ใช้ map interface{} - flexible แต่ไม่มี type safety
ctx.Info("user created", map[string]interface{}{
    "user_id": "123",
    "email": "user@example.com",
    "age": 25,
})

// Context propagation อัตโนมัติ
ctx.Info("processing user") // มี user_id อัตโนมัติ
```

## 2. Performance

### Zap
```go
// Zero allocation
logger.Info("process completed",
    zap.Int64("count", count),
    zap.Duration("duration", duration),
)

// Type safety ทำให้ไม่มี runtime allocation
type User struct {
    Name string `json:"name"`
}
logger.Info("user", zap.Object("user", user))
```

### Goobserv
```go
// ใช้ object pool ลด allocation
ctx.Info("process completed", map[string]interface{}{
    "count": count,
    "duration": duration,
})

// Flexible แต่มี allocation
type User struct {
    Name string `json:"name"`
}
ctx.Info("user", map[string]interface{}{
    "user": user,
})
```

## 3. Features

### Zap
```go
// จุดเด่น:
- Zero allocation
- Type safety
- High performance
- Flexible encoders
- Sampling
- Stacktraces

// ข้อจำกัด:
- ไม่มี context propagation
- ไม่มี auto tracking
- ต้อง setup เยอะ
- ไม่มี unified observability
```

### Goobserv
```go
// จุดเด่น:
- Auto context propagation
- Unified observability (logs, metrics, traces)
- Auto tracking
- Buffer management
- Multiple outputs
- Easy setup

// ข้อจำกัด:
- Performance ด้อยกว่า Zap
- ไม่มี type safety
- Memory usage สูงกว่า
```

## 4. Use Cases

### Zap เหมาะกับ:
```go
// 1. High-performance systems
logger.Info("processing",
    zap.Int64("items", 1000000),
    zap.Duration("duration", 100*time.Millisecond),
)

// 2. ต้องการ type safety
type Config struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}
logger.Info("config", zap.Object("config", config))

// 3. Resource-constrained environments
logger.With(
    zap.String("service", "api"),
).Info("started")
```

### Goobserv เหมาะกับ:
```go
// 1. Distributed systems
ctx.Info("request processed", map[string]interface{}{
    "trace_id": ctx.TraceID(),
    "duration": duration,
})

// 2. Microservices
ctx.Counter("requests_total").Inc()
ctx.Histogram("response_time").Observe(duration)

// 3. Development speed > Performance
ctx.Info("processing", map[string]interface{}{
    "data": complexData, // ไม่ต้องกำหนด schema
})
```

## 5. Configuration

### Zap
```go
// ต้อง config เยอะ
config := zap.Config{
    Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
    Encoding:         "json",
    EncoderConfig:    zap.NewProductionEncoderConfig(),
    OutputPaths:      []string{"stdout", "/var/log/app.log"},
    ErrorOutputPaths: []string{"stderr"},
}
logger, _ := config.Build()
```

### Goobserv
```go
// Simple config
obs := goobserv.New(goobserv.Config{
    Outputs: []OutputConfig{
        {Type: "stdout"},
        {Type: "file", Path: "/var/log/app.log"},
    },
})
```

## 6. Integration

### Zap
```go
// ต้อง integrate แยกกัน
logger.Info("request received")
metrics.Counter("requests").Inc()
span.AddEvent("processing")

// ต้องเขียน middleware เอง
func LogMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger := logger.With(
            zap.String("request_id", requestID),
        )
        next.ServeHTTP(w, r)
    })
}
```

### Goobserv
```go
// Unified observability
ctx.Info("request received") // auto track ทั้ง log, metric, trace

// Middleware มีให้
app.Use(goobserv.Middleware())
```

## แนะนำการเลือกใช้

### ใช้ Zap เมื่อ:
```go
// 1. ต้องการ performance สูงสุด
// 2. มี resource จำกัด
// 3. ต้องการ type safety
// 4. ไม่ต้องการ context propagation
// 5. ทีมมีประสบการณ์กับ Go สูง
```

### ใช้ Goobserv เมื่อ:
```go
// 1. ต้องการ setup ง่าย
// 2. ต้องการ unified observability
// 3. ทำ distributed tracing
// 4. Development speed สำคัญกว่า performance
// 5. ทีมมีประสบการณ์หลากหลาย
```
