# Fiber Middleware

Goobserv มี middleware สำหรับ [Fiber Web Framework](https://github.com/gofiber/fiber) เพื่อให้ใช้งาน observability ได้ง่ายขึ้น

## การติดตั้ง

```go
import (
    "github.com/vongga-platform/goobserv/pkg/middleware/fiber"
)
```

## การใช้งาน

```go
// สร้าง observer
obs := core.NewObserver(core.Config{
    BufferSize:    1000,
    FlushInterval: time.Second,
})

// สร้าง fiber app
app := fiber.New()

// เพิ่ม middleware
app.Use(fibermw.Middleware(fibermw.Config{
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
    GetTraceID func(*fiber.Ctx) string

    // Function สำหรับดึง request ID จาก request
    GetRequestID func(*fiber.Ctx) string
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

## Error Handling

```go
app := fiber.New(fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        // ดึง observer context
        ctx := fibermw.GetContext(c)

        // Log error
        obs.Error(ctx, "Request error").
            WithError(err)

        // Return error response
        code := fiber.StatusInternalServerError
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
        }
        return c.Status(code).JSON(fiber.Map{
            "error": err.Error(),
        })
    },
})
```

## ตัวอย่าง

```go
app.Get("/hello", func(c *fiber.Ctx) error {
    // ดึง observer context
    ctx := fibermw.GetContext(c)

    // Log info
    obs.Info(ctx, "Processing hello request")

    // Start span
    span, ctx := obs.StartSpan(ctx, "process_hello")
    defer obs.EndSpan(span)

    // Return response
    return c.JSON(fiber.Map{
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
   - ใช้ custom error handler
   - Log errors ด้วย error level
   - เพิ่ม error details ใน log
   - ใช้ error status code ที่เหมาะสม
