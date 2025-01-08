# หลักการทำงานของ Goobserv

## แนวคิดหลัก

Goobserv ถูกออกแบบมาเพื่อให้ทำ Observability ใน Go Application โดยเฉพาะ Clean Architecture โดยมีหลักการดังนี้:

### 1. Auto Tracking

```go
// ไม่ต้องเขียนโค้ดเพิ่ม แค่ส่ง context
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    return h.usecase.CreateUser(ctx, input)
}

// Goobserv จะ track ให้อัตโนมัติ:
// - Request ID
// - Trace ID
// - Layer (Handler -> UseCase -> Repository)
// - Input/Output
// - Duration
// - Errors
```

### 2. Clean Architecture Support

```plaintext
┌──────────────┐
│   Handler    │ ─── Auto track request/response
└──────────────┘
       │
┌──────────────┐
│   UseCase    │ ─── Auto track business logic
└──────────────┘
       │
┌──────────────┐
│  Repository  │ ─── Auto track database operations
└──────────────┘
```

### 3. Performance First

```go
// Non-blocking logging
type Observer struct {
    buffer    chan Entry      // buffer channel
    logFile   *os.File       // main file
    tempFile  *os.File       // backup file
}

// Async processing
go func() {
    for entry := range buffer {
        writeToFile(entry)
    }
}()
```

## การทำงานของระบบ

### 1. Request Flow

```plaintext
HTTP Request
    │
    ▼
Middleware (สร้าง context)
    │
    ▼
Handler (auto track input/output)
    │
    ▼
UseCase (auto track business logic)
    │
    ▼
Repository (auto track db operations)
    │
    ▼
Response (รวม logs/metrics)
```

### 2. Data Flow

```plaintext
Application
    │
    ▼
Buffer (in-memory)
    │
    ▼
Temp File (backup)
    │
    ▼
Main File (storage)
```

### 3. Error Flow

```plaintext
Error เกิดขึ้น
    │
    ▼
Record ใน Context
    │
    ▼
Propagate ขึ้นไป
    │
    ▼
Log พร้อม Stack Trace
```

## ระบบย่อย

### 1. Context Management

```go
// สร้าง context
ctx := obs.NewContext()

// ส่งต่อไปแต่ละ layer
handler -> usecase -> repository

// รวบรวมข้อมูลตอนจบ
ctx.Complete()
```

### 2. Buffer Management

```go
// เขียนลง buffer
buffer <- entry

// ถ้า buffer เต็ม
if bufferFull {
    writeToTempFile()
}

// flush buffer
flushToMainFile()
```

### 3. Storage Management

```go
// การ rotate file
if fileSize > maxSize {
    rotateFile()
}

// การลบไฟล์เก่า
if fileAge > maxAge {
    deleteOldFiles()
}
```

## การ Recovery

### 1. Application Crash

```plaintext
1. Application พัง
2. Buffer ยังมีข้อมูล
3. Temp File มีข้อมูลสำรอง
4. Restart Application
5. Recovery จาก Temp File
```

### 2. Disk Full

```plaintext
1. Disk เต็ม
2. เขียน Main File ไม่ได้
3. เขียนลง Emergency File
4. Alert Admin
5. Clean up เมื่อมีที่ว่าง
```

## การ Monitor ระบบ

### 1. Health Check

```go
// ตรวจสอบ buffer
if obs.BufferSize() > threshold {
    alert("Buffer nearly full")
}

// ตรวจสอบ disk
if obs.DiskUsage() > threshold {
    alert("Disk space low")
}
```

### 2. Performance Metrics

```go
// ดู metrics
metrics := obs.GetMetrics()

// ตัวอย่าง metrics
- Buffer size
- Write latency
- Error rate
- Request duration
```

## Best Practices

### 1. Configuration

```go
obs := goobserv.New(goobserv.Config{
    // ตั้ง buffer ให้เหมาะสม
    BufferSize: 1000,
    
    // ตั้ง file size ให้เหมาะสม
    MaxFileSize: 100 * 1024 * 1024,
    
    // ตั้ง retention ให้เหมาะสม
    RetentionDays: 7,
})
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
    fmt.Printf("error: %v\n", err)
    return err
}
```

### 3. Resource Management

```go
// ✅ ดี
obs := goobserv.New()
defer obs.Shutdown() // ปิดให้เรียบร้อย

// ❌ ไม่ดี
obs := goobserv.New()
// ไม่มีการ shutdown
```
