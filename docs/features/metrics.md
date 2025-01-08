# ระบบ Metrics

## ประเภทของ Metrics

### 1. Auto Metrics

```go
// Metrics ที่เก็บอัตโนมัติ:
- request_duration_seconds
- request_count_total
- error_count_total
- db_operation_duration_seconds
- memory_usage_bytes
- goroutine_count
```

### 2. Custom Metrics

```go
// Counter
ctx.Counter("user_created_total").Inc()

// Gauge
ctx.Gauge("active_users").Set(100)

// Histogram
ctx.Histogram("response_size_bytes").Observe(len(response))

// Summary
ctx.Summary("request_duration_seconds").Observe(duration)
```

## การใช้งาน Metrics

### 1. Request Metrics

```go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx := c.Locals("ctx").(*goobserv.Context)
    
    // เพิ่ม counter
    ctx.Counter("api_requests_total", 
        map[string]string{
            "method": "POST",
            "path": "/users",
        },
    ).Inc()
    
    // วัด duration
    defer ctx.Timer("api_duration_seconds").Record()
    
    return h.usecase.CreateUser(ctx, input)
}
```

### 2. Business Metrics

```go
func (u *UseCase) CreateUser(ctx *goobserv.Context, input Input) error {
    // นับจำนวน user ตาม plan
    ctx.Counter("users_created_total",
        map[string]string{
            "plan": input.Plan,
        },
    ).Inc()
    
    // เก็บขนาดข้อมูล
    ctx.Histogram("user_data_size_bytes").Observe(len(input.Data))
    
    return u.repo.Create(ctx, input)
}
```

### 3. Database Metrics

```go
func (r *Repository) Create(ctx *goobserv.Context, input Input) error {
    // วัด db latency
    timer := ctx.Timer("db_operation_seconds",
        map[string]string{
            "operation": "insert",
            "table": "users",
        },
    )
    defer timer.Record()
    
    // นับ queries
    ctx.Counter("db_queries_total").Inc()
    
    return r.db.Create(input)
}
```

## การตั้งค่า Metrics

### 1. Metric Types

```go
type MetricConfig struct {
    // Counter
    Counters []CounterDefinition
    
    // Gauge
    Gauges []GaugeDefinition
    
    // Histogram
    Histograms []HistogramDefinition
    
    // Summary
    Summaries []SummaryDefinition
}

// ตัวอย่าง
obs := goobserv.New(goobserv.Config{
    Metric: MetricConfig{
        Counters: []CounterDefinition{
            {
                Name: "api_requests_total",
                Help: "Total API requests",
                Labels: []string{"method", "path"},
            },
        },
        Histograms: []HistogramDefinition{
            {
                Name: "response_size_bytes",
                Help: "Response size distribution",
                Buckets: []float64{100, 1000, 10000},
            },
        },
    },
})
```

### 2. Labels

```go
// Global labels
obs := goobserv.New(goobserv.Config{
    Metric: MetricConfig{
        DefaultLabels: map[string]string{
            "service": "user-service",
            "version": "1.0.0",
            "env": "production",
        },
    },
})

// Metric specific labels
ctx.Counter("api_requests_total",
    map[string]string{
        "method": "POST",
        "status": "success",
    },
).Inc()
```

## การดู Metrics

### 1. Prometheus Format

```go
// /metrics endpoint
func MetricsHandler(c *fiber.Ctx) error {
    metrics := obs.GetMetrics()
    return c.SendString(metrics.ToPrometheus())
}

// Output:
# HELP api_requests_total Total API requests
# TYPE api_requests_total counter
api_requests_total{method="POST",path="/users"} 100

# HELP api_duration_seconds API request duration
# TYPE api_duration_seconds histogram
api_duration_seconds_bucket{le="0.1"} 90
api_duration_seconds_bucket{le="0.5"} 95
api_duration_seconds_bucket{le="1.0"} 99
api_duration_seconds_bucket{le="+Inf"} 100
```

### 2. JSON Format

```go
// Get metrics as JSON
metrics := obs.GetMetrics()
jsonMetrics := metrics.ToJSON()

// Output:
{
    "counters": {
        "api_requests_total": {
            "value": 100,
            "labels": {
                "method": "POST",
                "path": "/users"
            }
        }
    },
    "histograms": {
        "api_duration_seconds": {
            "count": 100,
            "sum": 45.6,
            "buckets": [
                {"le": 0.1, "count": 90},
                {"le": 0.5, "count": 95},
                {"le": 1.0, "count": 99},
                {"le": "+Inf", "count": 100}
            ]
        }
    }
}
```

## Best Practices

### 1. Naming

```go
// ✅ ดี
request_duration_seconds
http_requests_total
db_connections_active

// ❌ ไม่ดี
duration
requests
connections
```

### 2. Labels

```go
// ✅ ดี
ctx.Counter("api_requests_total",
    map[string]string{
        "method": "POST",    // ค่าที่เป็นไปได้น้อย
        "status": "success", // ค่าที่เป็นไปได้น้อย
    },
)

// ❌ ไม่ดี
ctx.Counter("api_requests_total",
    map[string]string{
        "user_id": userID,      // ค่าที่เป็นไปได้เยอะเกินไป
        "request_id": reqID,    // ค่าที่เป็นไปได้เยอะเกินไป
    },
)
```

### 3. Performance

```go
// ✅ ดี
timer := ctx.Timer("operation_duration")
defer timer.Record()

// ❌ ไม่ดี
start := time.Now()
// ... operation
duration := time.Since(start)
ctx.Histogram("operation_duration").Observe(duration.Seconds())
```
