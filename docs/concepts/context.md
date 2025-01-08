# การจัดการ Context

## โครงสร้างของ Context

```go
type Context struct {
    context.Context // embed standard context
    
    // ID Fields
    requestID     string
    traceID       string
    correlationID string
    userID        string
    
    // Timing
    startTime     time.Time
    timeout       time.Duration
    
    // Request Info
    method        string
    path          string
    userAgent     string
    clientIP      string
    
    // State
    currentSpan   *SpanData
    spans         []SpanData
    state         RequestState
    
    // Data
    metadata      map[string]interface{}
    metrics       map[string]float64
    errors        []error
}
```

## การสร้าง Context

### 1. จาก HTTP Request

```go
func (obs *Observer) Middleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // สร้าง context ใหม่
        ctx := obs.NewContext(
            WithTimeout(5*time.Second),
            WithRequestInfo(c.Method(), c.Path()),
        )
        
        // เก็บใน fiber
        c.Locals("ctx", ctx)
        
        return c.Next()
    }
}
```

### 2. จาก Background

```go
func ProcessJob() {
    // สร้าง context สำหรับ background job
    ctx := obs.NewContext(
        WithTimeout(30*time.Second),
        WithRequestInfo("BACKGROUND", "job-processor"),
    )
    
    // ทำงาน
    if err := process(ctx); err != nil {
        ctx.RecordError(err)
    }
}
```

## การใช้งาน Context

### 1. การเก็บ Metadata

```go
// เพิ่ม metadata
ctx.SetMetadata("user_id", "123")
ctx.SetMetadata("plan", "premium")

// ดึง metadata
userID := ctx.GetMetadata("user_id")
```

### 2. การจัดการ Error

```go
// บันทึก error
if err != nil {
    ctx.RecordError(err)
    return err
}

// ดู errors ทั้งหมด
if ctx.HasErrors() {
    errors := ctx.Errors()
    for _, err := range errors {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### 3. การวัด Performance

```go
// เริ่มวัดเวลา
start := ctx.StartTimer("db_query")

// ทำงาน
result := db.Query()

// จบการวัดเวลา
ctx.EndTimer(start)
```

## การจัดการ Span

### 1. การสร้าง Span

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    
    // เริ่ม span
    ctx.StartSpan("create_user", "handler")
    defer ctx.EndSpan()
    
    return h.usecase.CreateUser(ctx, input)
}
```

### 2. Nested Spans

```go
func (u *UseCase) CreateUser(ctx *goobserv.Context, input Input) error {
    // เริ่ม span ใหม่
    ctx.StartSpan("business_logic", "usecase")
    defer ctx.EndSpan()
    
    // เรียก repository
    return u.repo.Create(ctx, input)
}
```

### 3. Span Attributes

```go
func (r *Repository) Create(ctx *goobserv.Context, input Input) error {
    // เริ่ม span พร้อม attributes
    ctx.StartSpan("db_operation", "repository",
        map[string]interface{}{
            "table": "users",
            "operation": "insert",
        },
    )
    defer ctx.EndSpan()
    
    return r.db.Create(input)
}
```

## Context Propagation

### 1. HTTP Headers

```go
// ส่ง context ไป service อื่น
headers := ctx.ToHeaders()
req.Headers = headers

// รับ context จาก headers
ctx := obs.FromHeaders(headers)
```

### 2. gRPC Metadata

```go
// ส่ง context ผ่าน gRPC
md := ctx.ToMetadata()
ctx = metadata.NewOutgoingContext(ctx, md)

// รับ context จาก metadata
md, _ := metadata.FromIncomingContext(ctx)
ctx = obs.FromMetadata(md)
```

## Best Practices

### 1. Context Timeout

```go
// ✅ ดี
ctx := obs.NewContext(
    WithTimeout(5*time.Second),
)

// ❌ ไม่ดี
ctx := obs.NewContext() // ไม่มี timeout
```

### 2. Error Handling

```go
// ✅ ดี
if err != nil {
    ctx.RecordError(err)
    return err
}

// ❌ ไม่ดี
if err != nil {
    return err // ไม่ได้ record error
}
```

### 3. Metadata

```go
// ✅ ดี
ctx.SetMetadata("user_id", userID)
ctx.SetMetadata("action", "create")

// ❌ ไม่ดี
ctx.SetMetadata("data", map[string]interface{}{
    "user": user, // อาจมีข้อมูลเยอะเกินไป
})
```
