# แนวทางการจัดการ Error

## หลักการจัดการ Error

### 1. Error Recording

```go
// บันทึก error ใน context
func (s *Service) CreateUser(ctx *goobserv.Context, input Input) error {
    if err := validate(input); err != nil {
        // บันทึก error พร้อม metadata
        ctx.RecordError(err, map[string]interface{}{
            "input": input,
            "validation": "failed",
        })
        return err
    }
    return nil
}
```

### 2. Error Context

```go
// เพิ่ม context ให้ error
type ErrorContext struct {
    Error     error
    Stack     []string
    Metadata  map[string]interface{}
    Timestamp time.Time
    Layer     string
}

// การใช้งาน
ctx.RecordError(err, map[string]interface{}{
    "user_id": userID,
    "action": "create_user",
    "layer": "usecase",
})
```

### 3. Error Propagation

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    
    // Error จะถูก propagate และ track อัตโนมัติ
    if err := h.usecase.CreateUser(ctx, input); err != nil {
        return err
    }
    
    return nil
}
```

## การจัดการ Error ในแต่ละ Layer

### 1. Handler Layer

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    
    // Validate input
    input := new(CreateUserInput)
    if err := c.BodyParser(input); err != nil {
        ctx.RecordError(err, map[string]interface{}{
            "error_type": "input_validation",
        })
        return fiber.NewError(400, "Invalid input")
    }
    
    // Call usecase
    user, err := h.usecase.CreateUser(ctx, input)
    if err != nil {
        // Error จาก usecase จะมี context แล้ว
        return err
    }
    
    return c.JSON(user)
}
```

### 2. UseCase Layer

```go
func (u *UseCase) CreateUser(ctx *goobserv.Context, input Input) (*User, error) {
    // Business validation
    if err := u.validateBusiness(input); err != nil {
        ctx.RecordError(err, map[string]interface{}{
            "error_type": "business_validation",
            "rules": "business_rules",
        })
        return nil, err
    }
    
    // Call repository
    user, err := u.repo.Create(ctx, input)
    if err != nil {
        // Error จาก repository จะมี context แล้ว
        return nil, err
    }
    
    return user, nil
}
```

### 3. Repository Layer

```go
func (r *Repository) Create(ctx *goobserv.Context, input Input) (*User, error) {
    // Database operation
    result := r.db.Create(&input)
    if result.Error != nil {
        ctx.RecordError(result.Error, map[string]interface{}{
            "error_type": "database",
            "operation": "create",
            "table": "users",
        })
        return nil, result.Error
    }
    
    return &input, nil
}
```

## Error Types

### 1. Custom Errors

```go
// กำหนด error types
type ErrorType string

const (
    ErrorTypeValidation ErrorType = "validation"
    ErrorTypeBusiness   ErrorType = "business"
    ErrorTypeDatabase   ErrorType = "database"
    ErrorTypeSystem     ErrorType = "system"
)

// Custom error
type AppError struct {
    Type    ErrorType
    Message string
    Cause   error
}

func (e *AppError) Error() string {
    return e.Message
}
```

### 2. Error Wrapping

```go
// Wrap error with context
func WrapError(err error, errType ErrorType, msg string) error {
    return &AppError{
        Type:    errType,
        Message: msg,
        Cause:   err,
    }
}

// การใช้งาน
if err := validate(input); err != nil {
    return WrapError(err, ErrorTypeValidation, "invalid input")
}
```

## Error Analysis

### 1. Error Aggregation

```go
// รวม errors ตาม type
type ErrorAggregation struct {
    Type     ErrorType
    Count    int
    Examples []ErrorContext
}

func AnalyzeErrors(ctx *goobserv.Context) []ErrorAggregation {
    errors := ctx.GetErrors()
    // Aggregate errors
    return aggregations
}
```

### 2. Error Metrics

```go
// เก็บ metrics ของ errors
ctx.Counter("errors_total",
    map[string]string{
        "type": string(err.Type),
        "layer": ctx.CurrentLayer(),
    },
).Inc()
```

## Best Practices

### 1. Error Recording

```go
// ✅ ดี
if err != nil {
    ctx.RecordError(err, map[string]interface{}{
        "error_type": "validation",
        "field": "email",
    })
    return err
}

// ❌ ไม่ดี
if err != nil {
    return err // ไม่ได้ record error
}
```

### 2. Error Context

```go
// ✅ ดี
ctx.RecordError(err, map[string]interface{}{
    "user_id": userID,
    "action": "create_user",
    "layer": "usecase",
})

// ❌ ไม่ดี
ctx.RecordError(err) // ไม่มี context
```

### 3. Error Types

```go
// ✅ ดี
return WrapError(err, ErrorTypeValidation, "invalid email format")

// ❌ ไม่ดี
return errors.New("error") // ไม่มี type และ context
```
