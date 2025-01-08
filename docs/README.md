# Goobserv SDK

Goobserv เป็น SDK สำหรับทำ Observability ใน Go Application โดยเฉพาะสำหรับ Clean Architecture

## 📚 สารบัญ

### 🚀 เริ่มต้นใช้งาน
- [การติดตั้ง](getting-started/installation.md)
- [เริ่มต้นใช้งาน](getting-started/quick-start.md)
- [การตั้งค่า](getting-started/configuration.md)

### 📖 แนวคิดหลัก
- [สถาปัตยกรรม](concepts/architecture.md)
- [การจัดการ Context](concepts/context.md)
- [หลักการทำงาน](concepts/observability.md)

### ⚡ ฟีเจอร์
- [ระบบ Logging](features/logging.md)
- [ระบบ Tracing](features/tracing.md)
- [ระบบ Metrics](features/metrics.md)

### 👍 แนวทางการใช้งานที่ดี
- [การจัดการ Performance](best-practices/performance.md)
- [การจัดการ Error](best-practices/error-handling.md)
- [การจัดการความปลอดภัย](best-practices/security.md)

## 🎯 จุดเด่น

- **Auto Context Tracking**: ติดตาม Request ID, Trace ID และ Layer tracking อัตโนมัติ
- **Performance Metrics**: วัดประสิทธิภาพทั้ง Request, Layer และ Database
- **Error Handling**: จัดการ Error พร้อม Stack trace และ Context
- **Storage System**: ระบบจัดเก็บข้อมูลพร้อม Buffer และ Backup

## 🛠 ตัวอย่างการใช้งาน

```go
func main() {
    // สร้าง Observer
    obs := goobserv.New()
    defer obs.Shutdown()

    // ใช้กับ Fiber
    app := fiber.New()
    app.Use(goobserv.Middleware())
    
    app.Listen(":3000")
}
```

## 🤝 การสนับสนุน

- [GitHub Issues](https://github.com/yourusername/goobserv/issues)
- [เอกสารประกอบ](https://github.com/yourusername/goobserv/docs)
- [ตัวอย่างโค้ด](https://github.com/yourusername/goobserv/examples)

## 📝 License

MIT License
