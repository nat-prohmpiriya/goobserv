# แนวทางการจัดการความปลอดภัย

## การป้องกันข้อมูลสำคัญ

### 1. Sensitive Data

```go
// กำหนดข้อมูลที่ต้อง mask
type SecurityConfig struct {
    // Fields ที่ต้อง mask
    SensitiveKeys []string
    
    // Pattern สำหรับ mask
    MaskPattern string
    
    // Custom masking function
    MaskFunc func(string) string
}

// ตัวอย่าง
obs := goobserv.New(goobserv.Config{
    Security: SecurityConfig{
        SensitiveKeys: []string{
            "password",
            "token",
            "credit_card",
            "secret",
        },
        MaskPattern: "***",
    },
})
```

### 2. Auto Masking

```go
// ระบบจะ mask ข้อมูลอัตโนมัติ
func (h *Handler) Login(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    
    input := struct {
        Email    string `json:"email"`
        Password string `json:"password"` // จะถูก mask อัตโนมัติ
    }{}
    
    if err := c.BodyParser(&input); err != nil {
        return err
    }
    
    // Log จะเป็น:
    // {
    //     "email": "user@example.com",
    //     "password": "***"
    // }
    
    return h.usecase.Login(ctx, input)
}
```

## การจัดการ Access

### 1. Log Access Control

```go
type AccessConfig struct {
    // ระดับการเข้าถึง
    Level AccessLevel
    
    // ผู้ที่สามารถเข้าถึง
    AllowedUsers []string
    
    // IP ที่อนุญาต
    AllowedIPs []string
}

// การใช้งาน
obs := goobserv.New(goobserv.Config{
    Access: AccessConfig{
        Level: AccessLevelRestricted,
        AllowedUsers: []string{
            "admin",
            "monitoring",
        },
        AllowedIPs: []string{
            "10.0.0.0/8",
            "192.168.0.0/16",
        },
    },
})
```

### 2. Audit Logging

```go
// บันทึกการเข้าถึง
type AuditLog struct {
    Time      time.Time
    User      string
    Action    string
    Resource  string
    IP        string
    UserAgent string
}

// การใช้งาน
ctx.AuditLog("view_logs", map[string]interface{}{
    "files": []string{"app.log"},
    "lines": 100,
})
```

## การเข้ารหัสข้อมูล

### 1. Log Encryption

```go
type EncryptionConfig struct {
    // เปิด/ปิดการเข้ารหัส
    Enabled bool
    
    // Algorithm
    Algorithm string
    
    // Key management
    KeyFile string
    
    // Rotation
    KeyRotationInterval time.Duration
}

// การใช้งาน
obs := goobserv.New(goobserv.Config{
    Encryption: EncryptionConfig{
        Enabled: true,
        Algorithm: "AES-256",
        KeyFile: "/etc/goobserv/keys/log.key",
        KeyRotationInterval: 30 * 24 * time.Hour,
    },
})
```

### 2. Transport Security

```go
type TransportConfig struct {
    // TLS config
    TLS *tls.Config
    
    // Certificate
    CertFile string
    KeyFile  string
    
    // Mutual TLS
    ClientCAs *x509.CertPool
}
```

## Data Protection

### 1. Retention Policy

```go
type RetentionConfig struct {
    // ระยะเวลาเก็บข้อมูล
    Duration time.Duration
    
    // การลบข้อมูล
    DeleteMethod DeleteMethod
    
    // Backup ก่อนลบ
    BackupBeforeDelete bool
}

// การใช้งาน
obs := goobserv.New(goobserv.Config{
    Retention: RetentionConfig{
        Duration: 90 * 24 * time.Hour, // 90 วัน
        DeleteMethod: DeleteMethodSecure,
        BackupBeforeDelete: true,
    },
})
```

### 2. Data Sanitization

```go
type SanitizeConfig struct {
    // Fields ที่ต้อง sanitize
    Fields []string
    
    // Rules
    Rules map[string]SanitizeRule
    
    // Custom sanitizer
    Sanitizer func(string) string
}

// การใช้งาน
obs := goobserv.New(goobserv.Config{
    Sanitize: SanitizeConfig{
        Fields: []string{"email", "phone"},
        Rules: map[string]SanitizeRule{
            "email": EmailRule,
            "phone": PhoneRule,
        },
    },
})
```

## Best Practices

### 1. Sensitive Data

```go
// ✅ ดี
obs := goobserv.New(goobserv.Config{
    Security: SecurityConfig{
        SensitiveKeys: []string{
            "password",
            "token",
            "credit_card",
        },
    },
})

// ❌ ไม่ดี
obs := goobserv.New(goobserv.Config{
    // ไม่ได้กำหนด sensitive data
})
```

### 2. Access Control

```go
// ✅ ดี
obs := goobserv.New(goobserv.Config{
    Access: AccessConfig{
        Level: AccessLevelRestricted,
        AllowedUsers: []string{"admin"},
        AllowedIPs: []string{"10.0.0.0/8"},
    },
})

// ❌ ไม่ดี
obs := goobserv.New(goobserv.Config{
    Access: AccessConfig{
        Level: AccessLevelPublic, // ไม่ควรเปิดให้เข้าถึงได้ทั้งหมด
    },
})
```

### 3. Encryption

```go
// ✅ ดี
obs := goobserv.New(goobserv.Config{
    Encryption: EncryptionConfig{
        Enabled: true,
        Algorithm: "AES-256",
        KeyRotationInterval: 30 * 24 * time.Hour,
    },
})

// ❌ ไม่ดี
obs := goobserv.New(goobserv.Config{
    Encryption: EncryptionConfig{
        Enabled: false, // ไม่ได้เข้ารหัสข้อมูล
    },
})
```
