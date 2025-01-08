# ระบบ Logging

## การทำงาน

ระบบ Logging ใน Goobserv ถูกออกแบบให้ทำงานอัตโนมัติและมีประสิทธิภาพ

### 1. Auto Logging

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    // จะ log อัตโนมัติ:
    // - Input
    // - Output
    // - Error (ถ้ามี)
    // - Duration
    return h.usecase.CreateUser(ctx, input)
}
```

### 2. Log Format

```json
{
    "request_id": "123",
    "trace_id": "456",
    "time": "2025-01-08T21:44:28+07:00",
    "level": "info",
    "message": "create user",
    "layer": "handler",
    "input": {
        "name": "test",
        "email": "test@example.com"
    },
    "output": {
        "id": "789",
        "status": "success"
    },
    "duration": 0.123
}
```

### 3. Log Levels

```go
// ระดับของ Log
ctx.Debug("debug message")   // สำหรับ developer
ctx.Info("info message")    // ข้อมูลทั่วไป
ctx.Warn("warn message")    // คำเตือน
ctx.Error("error message")  // error ที่เกิดขึ้น
```

## การจัดการ Log

### 1. Log Storage

```go
type StorageConfig struct {
    Path        string        // ที่เก็บไฟล์
    MaxSize     int64        // ขนาดสูงสุด
    MaxAge      int          // อายุสูงสุด
    MaxBackups  int          // จำนวน backup
}
```

### 2. Log Rotation

```go
// ตั้งค่า rotation
obs := goobserv.New(goobserv.Config{
    LogRotation: LogRotationConfig{
        MaxSize:    100,    // MB
        MaxAge:     7,      // วัน
        MaxBackups: 5,      // ไฟล์
    },
})
```

### 3. Emergency Logging

```go
// เมื่อ buffer เต็มหรือมีปัญหา
type EmergencyConfig struct {
    Path     string    // ที่เก็บ emergency log
    MaxSize  int64     // ขนาดสูงสุด
}
```

## การใช้งานขั้นสูง

### 1. Custom Logger

```go
type CustomLogger struct {
    goobserv.Logger
}

func (l *CustomLogger) Log(level string, msg string, data map[string]interface{}) {
    // จัดการ log ตามต้องการ
}
```

### 2. Log Filtering

```go
obs := goobserv.New(goobserv.Config{
    LogFilter: func(entry LogEntry) bool {
        // กรอง log ตามเงื่อนไข
        return entry.Level != "debug"
    },
})
```

### 3. Sensitive Data

```go
obs := goobserv.New(goobserv.Config{
    SensitiveKeys: []string{
        "password",
        "token",
        "credit_card",
    },
})
```

## Best Practices

### 1. การเขียน Log Message

```go
// ✅ ดี
ctx.Info("user created successfully", map[string]interface{}{
    "user_id": "123",
    "plan": "premium",
})

// ❌ ไม่ดี
ctx.Info("user 123 created with premium plan")
```

### 2. การจัดการ Error

```go
// ✅ ดี
if err != nil {
    ctx.Error("failed to create user", map[string]interface{}{
        "error": err.Error(),
        "user_id": userID,
    })
    return err
}

// ❌ ไม่ดี
if err != nil {
    ctx.Error(err.Error())
    return err
}
```

### 3. Performance

```go
// ✅ ดี
if ctx.IsDebugEnabled() {
    ctx.Debug("expensive operation result", map[string]interface{}{
        "result": expensiveOperation(),
    })
}

// ❌ ไม่ดี
ctx.Debug("expensive operation result", map[string]interface{}{
    "result": expensiveOperation(), // ทำงานทุกครั้ง
})
```
