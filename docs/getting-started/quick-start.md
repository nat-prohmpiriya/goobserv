# เริ่มต้นใช้งาน Goobserv

## การใช้งานพื้นฐาน

### 1. สร้าง Observer

```go
package main

import (
    "github.com/yourusername/goobserv"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // สร้าง observer
    obs := goobserv.New(goobserv.Config{
        ServiceName: "my-service",
        Environment: "production",
    })
    defer obs.Shutdown()

    // ใช้กับ fiber
    app := fiber.New()
    app.Use(goobserv.Middleware())

    app.Listen(":3000")
}
```

### 2. การใช้งานใน Handler

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    // ดึง context จาก fiber
    ctx := c.Locals("ctx").(*goobserv.Context)
    
    // ทำงานปกติ - goobserv จะ track ให้อัตโนมัติ
    result, err := h.usecase.CreateUser(ctx, input)
    if err != nil {
        return err
    }
    
    return c.JSON(result)
}
```

### 3. การใช้งานใน UseCase

```go
func (u *UseCase) CreateUser(ctx *goobserv.Context, input Input) (*User, error) {
    // เพิ่ม custom log ถ้าต้องการ
    ctx.Log("info", "creating user", map[string]interface{}{
        "email": input.Email,
    })
    
    // เรียก repository
    return u.repo.Create(ctx, input)
}
```

### 4. การใช้งานใน Repository

```go
func (r *Repository) Create(ctx *goobserv.Context, input Input) (*User, error) {
    // วัด performance
    defer ctx.RecordMetric("db_operation", 1)
    
    // ทำงานกับ database
    return r.db.Create(input)
}
```

## ตัวอย่างการใช้งานจริง

### 1. REST API

```go
func main() {
    obs := goobserv.New()
    defer obs.Shutdown()

    app := fiber.New()
    app.Use(goobserv.Middleware())

    app.Post("/users", func(c *fiber.Ctx) error {
        ctx := c.Locals("ctx").(*goobserv.Context)
        
        input := new(CreateUserInput)
        if err := c.BodyParser(input); err != nil {
            return err
        }
        
        user, err := userService.CreateUser(ctx, input)
        if err != nil {
            return err
        }
        
        return c.JSON(user)
    })

    app.Listen(":3000")
}
```

### 2. gRPC Service

```go
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
    // สร้าง goobserv context
    obsCtx := goobserv.FromContext(ctx)
    
    // ทำงานปกติ
    user, err := s.usecase.CreateUser(obsCtx, req)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 3. Background Job

```go
func ProcessQueue(ctx context.Context) {
    obs := goobserv.New()
    defer obs.Shutdown()

    for {
        // สร้าง context สำหรับแต่ละ job
        jobCtx := obs.NewContext()
        
        // ทำงาน
        if err := processJob(jobCtx); err != nil {
            jobCtx.Error("job failed", map[string]interface{}{
                "error": err.Error(),
            })
        }
    }
}
```

## การตรวจสอบผลลัพธ์

### 1. ดู Logs

```bash
# ดู log ทั้งหมด
cat /var/log/goobserv/app.log

# ดู error log
cat /var/log/goobserv/app.log | grep "level\":\"error"
```

### 2. ดู Metrics

```go
// ดู metrics ทั้งหมด
metrics := obs.GetMetrics()
for name, value := range metrics {
    fmt.Printf("%s: %v\n", name, value)
}
```

### 3. ดู Traces

```go
// ดู trace ของ request
traces := obs.GetTraces()
for _, trace := range traces {
    fmt.Printf("Request ID: %s, Duration: %v\n", trace.RequestID, trace.Duration)
}
```

## ขั้นตอนต่อไป

1. ศึกษาการตั้งค่าเพิ่มเติมใน [configuration.md](configuration.md)
2. ศึกษาการทำงานของ Context ใน [context.md](../concepts/context.md)
3. ดูแนวทางการใช้งานที่ดีใน [best-practices](../best-practices/)
