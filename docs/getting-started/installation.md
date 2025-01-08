# การติดตั้ง Goobserv

## วิธีติดตั้ง

```bash
go get github.com/yourusername/goobserv
```

## ความต้องการของระบบ

- Go 1.16+
- เนื้อที่ว่างสำหรับ Log Files
- สิทธิ์ในการเขียนไฟล์

## การตั้งค่าเริ่มต้น

1. สร้าง Observer:

```go
obs := goobserv.New(goobserv.Config{
    ServiceName: "my-service",
    Environment: "production",
})
```

2. ใช้กับ Web Framework:

```go
// Fiber
app.Use(goobserv.Middleware())

// Echo
e.Use(goobserv.EchoMiddleware())

// Gin
r.Use(goobserv.GinMiddleware())
```

## การตรวจสอบการติดตั้ง

```go
// ทดสอบ logging
obs.Log("info", "test message")

// ทดสอบ metrics
obs.RecordMetric("test_metric", 1)

// ตรวจสอบไฟล์ log
cat /var/log/goobserv/app.log
```

## ปัญหาที่พบบ่อย

1. **ไม่สามารถเขียนไฟล์ได้**
   - ตรวจสอบสิทธิ์การเขียนไฟล์
   - ตรวจสอบพื้นที่ว่าง

2. **Import ไม่ได้**
   - ตรวจสอบ Go version
   - ตรวจสอบ go.mod

3. **Middleware ไม่ทำงาน**
   - ตรวจสอบลำดับ middleware
   - ตรวจสอบการตั้งค่า context
