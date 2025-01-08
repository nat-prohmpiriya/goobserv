# สถาปัตยกรรมของ Goobserv

## โครงสร้างหลัก

### 1. Observer
- เป็น Core Component หลัก
- จัดการการทำงานทั้งหมด
- รับผิดชอบการเก็บและส่งข้อมูล

```go
type Observer struct {
    buffer    chan Entry     // buffer สำหรับเก็บข้อมูล
    logFile   *os.File      // ไฟล์หลัก
    tempFile  *os.File      // ไฟล์สำรอง
}
```

### 2. Context
- เก็บข้อมูล Request
- ติดตาม Layer
- จัดการ Error

```go
type Context struct {
    RequestID string
    TraceID   string
    StartTime time.Time
    Layer     string
}
```

### 3. Buffer System
- ระบบ Buffer แบบ Non-blocking
- มีระบบ Backup
- Auto Flush

```go
// การทำงานของ Buffer
buffer := make(chan Entry, 1000)
go func() {
    for entry := range buffer {
        // เขียนลง temp file ก่อน
        writeTempFile(entry)
        // เขียนลง main file
        writeMainFile(entry)
        // ลบจาก temp file
        cleanTempFile()
    }
}()
```

## การไหลของข้อมูล

1. **Request เข้ามา**
   ```
   HTTP Request -> Middleware -> Create Context
   ```

2. **ผ่าน Layer**
   ```
   Handler -> UseCase -> Repository
   ```

3. **เก็บข้อมูล**
   ```
   Buffer -> Temp File -> Main File
   ```

## ระบบความปลอดภัย

1. **Backup System**
   - เก็บข้อมูลชั่วคราว
   - กู้คืนเมื่อมีปัญหา

2. **Error Handling**
   - เก็บ Stack Trace
   - เก็บ Context

3. **Data Protection**
   - Redact ข้อมูลสำคัญ
   - Sanitize Input

## การ Scale

1. **Buffer Size**
   ```go
   type Config struct {
       BufferSize    int
       FlushInterval time.Duration
   }
   ```

2. **File Rotation**
   ```go
   type FileConfig struct {
       MaxSize    int64
       MaxAge     int
       MaxBackups int
   }
   ```

3. **Sampling**
   ```go
   type SamplingConfig struct {
       Rate     float64
       Priority []string
   }
   ```

## แนวทางการพัฒนาต่อ

1. **Distributed Tracing**
   - OpenTelemetry Integration
   - Jaeger Support

2. **Metrics Export**
   - Prometheus Integration
   - Custom Exporters

3. **Log Aggregation**
   - ELK Stack Support
   - Custom Aggregators
