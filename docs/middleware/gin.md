# Gin Middleware

Goobserv มี middleware สำหรับ [Gin Web Framework](https://github.com/gin-gonic/gin) เพื่อให้ใช้งาน observability ได้ง่ายขึ้น

## การติดตั้ง

```go
import (
    "github.com/vongga-platform/goobserv/pkg/middleware/gin"
)
```

## การใช้งาน

```go
// สร้าง observer
obs := core.NewObserver(core.Config{
    BufferSize:    1000,
    FlushInterval: time.Second,
})

// สร้าง gin engine
r := gin.New()

// เพิ่ม middleware
r.Use(ginmw.Middleware(ginmw.Config{
    Observer: obs,
    SkipPaths: []string{
        "/health",
        "/metrics",
    },
}))
```

## Configuration

```go
type Config struct {
    // Observer instance (required)
    Observer *core.Observer

    // Paths ที่จะไม่ถูก log
    SkipPaths []string

    // Function สำหรับดึง trace ID จาก request
    GetTraceID func(*gin.Context) string

    // Function สำหรับดึง request ID จาก request
    GetRequestID func(*gin.Context) string
}
```

## Features

### Auto Logging
- HTTP Method
- Path
- Status Code
- Duration
- Client IP
- User Agent
- Request Size
- Response Size
- Errors

### Context Propagation
- Trace ID
- Request ID
- Custom Fields

### Performance Metrics
- Request Count
- Request Duration
- Error Count

## ตัวอย่าง

```go
r.GET("/hello", func(c *gin.Context) {
    // ดึง observer context
    ctx := ginmw.GetContext(c)

    // Log info
    obs.Info(ctx, "Processing hello request")

    // Start span
    span, ctx := obs.StartSpan(ctx, "process_hello")
    defer obs.EndSpan(span)

    // Return response
    c.JSON(200, gin.H{
        "message": "Hello, World!",
    })
})
```

## Best Practices

1. **Skip Paths**
   - Skip health check endpoints
   - Skip metrics endpoints
   - Skip static files

2. **Trace ID**
   - ใช้ custom function เพื่อดึง trace ID จาก header ที่ต้องการ
   - สร้าง trace ID ใหม่ถ้าไม่มี

3. **Request ID**
   - ใช้ custom function เพื่อดึง request ID จาก header ที่ต้องการ
   - สร้าง request ID ใหม่ถ้าไม่มี

4. **Error Handling**
   - Log errors ด้วย error level
   - เพิ่ม error details ใน log
   - ใช้ error status code ที่เหมาะสม
