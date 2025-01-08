# ระบบ Tracing

## การทำงานของ Tracing

### 1. Auto Tracing

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    
    // Auto trace:
    // - Layer transitions
    // - Function calls
    // - Database operations
    return h.usecase.CreateUser(ctx, input)
}
```

### 2. Trace Format

```json
{
    "trace_id": "abc123",
    "request_id": "req123",
    "start_time": "2025-01-08T22:01:15+07:00",
    "duration": 0.123,
    "spans": [
        {
            "name": "CreateUser",
            "layer": "handler",
            "start_time": "2025-01-08T22:01:15+07:00",
            "duration": 0.050,
            "input": {...},
            "output": {...}
        },
        {
            "name": "CreateUser",
            "layer": "usecase",
            "start_time": "2025-01-08T22:01:15.050+07:00",
            "duration": 0.040,
            "input": {...},
            "output": {...}
        },
        {
            "name": "Create",
            "layer": "repository",
            "start_time": "2025-01-08T22:01:15.090+07:00",
            "duration": 0.030,
            "input": {...},
            "output": {...}
        }
    ]
}
```

## การใช้งาน Tracing

### 1. Manual Tracing

```go
func ComplexOperation(ctx *goobserv.Context) error {
    // เริ่ม span ใหม่
    ctx.StartSpan("complex_operation", "usecase")
    defer ctx.EndSpan()
    
    // ทำงาน
    result := process()
    
    // เพิ่ม attributes
    ctx.AddAttribute("result_size", len(result))
    
    return nil
}
```

### 2. Nested Spans

```go
func (s *Service) ProcessOrder(ctx *goobserv.Context, order Order) error {
    // Parent span
    ctx.StartSpan("process_order", "usecase")
    defer ctx.EndSpan()
    
    // Child span 1
    ctx.StartSpan("validate_order", "usecase")
    if err := s.validateOrder(order); err != nil {
        ctx.EndSpan()
        return err
    }
    ctx.EndSpan()
    
    // Child span 2
    ctx.StartSpan("save_order", "usecase")
    if err := s.repo.SaveOrder(ctx, order); err != nil {
        ctx.EndSpan()
        return err
    }
    ctx.EndSpan()
    
    return nil
}
```

### 3. Error Tracing

```go
func (r *Repository) SaveUser(ctx *goobserv.Context, user User) error {
    ctx.StartSpan("save_user", "repository")
    defer ctx.EndSpan()
    
    if err := r.db.Create(user); err != nil {
        // บันทึก error ใน span
        ctx.RecordError(err)
        return err
    }
    
    return nil
}
```

## Trace Propagation

### 1. HTTP Headers

```go
// ส่ง trace ไป service อื่น
headers := ctx.ToHeaders()
req.Headers = headers

// รับ trace จาก headers
ctx := obs.FromHeaders(headers)
```

### 2. gRPC Metadata

```go
// ส่ง trace ผ่าน gRPC
md := ctx.ToMetadata()
ctx = metadata.NewOutgoingContext(ctx, md)

// รับ trace จาก metadata
md, _ := metadata.FromIncomingContext(ctx)
ctx = obs.FromMetadata(md)
```

## การวิเคราะห์ Trace

### 1. Trace Analysis

```go
// ดู trace ทั้งหมด
traces := obs.GetTraces()

// วิเคราะห์ performance
for _, trace := range traces {
    // หา slow operations
    for _, span := range trace.Spans {
        if span.Duration > threshold {
            fmt.Printf("Slow operation: %s\n", span.Name)
        }
    }
}
```

### 2. Performance Analysis

```go
// วิเคราะห์ latency
type LayerLatency struct {
    Layer    string
    Average  float64
    P95      float64
    P99      float64
}

func AnalyzeLatency(traces []Trace) []LayerLatency {
    // คำนวณ latency แต่ละ layer
    return latencies
}
```

### 3. Error Analysis

```go
// วิเคราะห์ errors
func AnalyzeErrors(traces []Trace) map[string]int {
    errors := make(map[string]int)
    
    for _, trace := range traces {
        for _, span := range trace.Spans {
            if span.Error != "" {
                errors[span.Error]++
            }
        }
    }
    
    return errors
}
```

## Best Practices

### 1. Span Naming

```go
// ✅ ดี
ctx.StartSpan("create_user", "usecase")
ctx.StartSpan("save_user", "repository")

// ❌ ไม่ดี
ctx.StartSpan("function1", "usecase")
ctx.StartSpan("do_something", "repository")
```

### 2. Error Handling

```go
// ✅ ดี
if err != nil {
    ctx.RecordError(err)
    ctx.EndSpan()
    return err
}

// ❌ ไม่ดี
if err != nil {
    return err // ไม่ได้ record error
}
```

### 3. Attributes

```go
// ✅ ดี
ctx.AddAttribute("user_id", userID)
ctx.AddAttribute("order_status", status)

// ❌ ไม่ดี
ctx.AddAttribute("data", largeObject) // ข้อมูลเยอะเกินไป
```
