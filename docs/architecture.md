# Goobserv Architecture

## Core Components

### Observer
The central component that manages logging and observability:
```go
obs := core.NewObserver(core.Config{
    BufferSize:    1000,
    FlushInterval: time.Second,
})
```

### Context
Carries request-specific information:
```go
ctx := core.NewContext(context.Background()).
    WithTraceID("trace-123").
    WithRequestID("req-456")
```

### Entry
Represents a single log entry:
```go
type Entry struct {
    Time      time.Time
    Level     Level
    Message   string
    TraceID   string
    SpanID    string
    RequestID string
    Data      map[string]interface{}
}
```

## Data Flow

1. **Request Handling**
   ```
   HTTP Request -> Middleware -> Application Handler
        |              |              |
        |              |              |
        v              v              v
   Extract IDs    Create Context   Use Context
   ```

2. **Logging Flow**
   ```
   Application Code -> Observer -> Buffer -> Output Handlers
           |              |          |            |
           |              |          |            |
           v              v          v            v
      Log Message    Create Entry  Queue    Write to Output
   ```

3. **Trace Flow**
   ```
   Start Request -> Create Span -> Process -> End Span
        |              |            |           |
        |              |            |           |
        v              v            v           v
   Generate IDs    Start Timer   Log Info   Record Duration
   ```

## Usage Example

```go
func main() {
    // 1. Create observer
    obs := core.NewObserver(core.Config{
        BufferSize:    1000,
        FlushInterval: time.Second,
    })
    defer obs.Close()

    // 2. Add outputs
    obs.AddOutput(&output.StdoutOutput{})

    // 3. Create middleware
    r := gin.New()
    r.Use(middleware.Middleware(middleware.Config{
        Observer: obs,
    }))

    // 4. Use in handlers
    r.GET("/api", func(c *gin.Context) {
        ctx := middleware.GetContext(c)
        obs.Info(ctx, "Processing request",
            "path", "/api",
            "method", "GET",
        )
        // ... handle request
    })
}
```

## Key Features

1. **Asynchronous Logging**
   - Uses buffered channels
   - Non-blocking write operations
   - Periodic flushing

2. **Context Management**
   - Trace ID propagation
   - Request ID tracking
   - Span management

3. **Middleware Integration**
   - Built-in support for Gin and Fiber
   - Extensible for other frameworks
   - Skip path configuration

4. **Flexible Output**
   - Multiple output handlers
   - Custom output support
   - Structured logging

## Testing Strategy

1. **Unit Tests**
   - Core functionality
   - Middleware behavior
   - Output handlers

2. **Integration Tests**
   - Framework integration
   - End-to-end request flow
   - Concurrent operations

3. **Test Utilities**
   - Test output handler
   - Mock observer
   - Helper functions

## Best Practices

1. **Observer Usage**
   - Create one observer per application
   - Close observer on shutdown
   - Configure appropriate buffer size

2. **Context Management**
   - Always use context from middleware
   - Propagate trace IDs
   - Create spans for operations

3. **Logging**
   - Use appropriate log levels
   - Include relevant context
   - Structure log data

4. **Error Handling**
   - Log errors with context
   - Include error details
   - Use appropriate log level


## Log Format
```json
{
    "request_id": "123",
    "trace_id": "456",
    "correlation_id": "789",
    "user_id": "user_123",
    "start_time": "2025-01-08T20:11:30+07:00",
    "end_time": "2025-01-08T20:11:30.500+07:00",
    "duration": 0.5,
    "state": "success",
    "headers": {

    },
   "client": {
      "ip": "203.0.113.195",
      "user_agent": {
         "browser": "Chrome",
         "browser_version": "120.0.0",
         "os": "MacOS",
         "os_version": "14.2.1",
         "device_type": "desktop",
         "is_mobile": false
      },
   },

    "spans": [
        {
            "function": "CreateUser",
            "layer": "handler",
            "start_time": "...",
            "end_time": "...",
            "input": {...},
            "output": {...},
            "span_id": "1",
            "parent_span_id": "0",
            "duration": 0.2
        },
        {
            "function": "CreateUser", // function name
            "layer": "usecase", // package name
            "start_time": "...",
            "end_time": "...",
            "input": {...},
            "output": {...},
            "span_id": "2",
            "parent_span_id": "1",
            "duration": 0.3
        },
        {
            "function": "CreateUser", // function name
            "layer": "usecase", // package name
            "start_time": "...",
            "end_time": "...",
            "input": {...},
            "output": {...},
            "level": "debug",
            "span_id": "3",
            "parent_span_id": "2",
            "duration": 0.3
        },
        {
            "function": "CreateUser",
            "layer": "repository",
            "start_time": "...",
            "end_time": "...",
            "input": {...},
            "output": {...},
            "span_id": "4",
            "parent_span_id": "3",
            "duration": 0.1
        }
    ],
    "metrics": {
        "service_duration": 0.4
    }
}
```

### new format 

```json
{
    // ข้อมูลพื้นฐานของ Request
    "request_id": "123",
    "trace_id": "456",
    "user_id": "user_123",
    "start_time": "2025-01-08T20:11:30+07:00",
    "end_time": "2025-01-08T20:11:30.500+07:00",
    "duration": 0.5,
    "state": "success",
    "method": "POST",
    "original_path": "/users?page=1",

    // Spans (การทำงานแต่ละขั้นตอน)
    "spans": [
        // auto create by lib
        {
            "function": "handler.CreateUser", // package.function 
            "start_time": "2025-01-08T20:11:30.000+07:00",
            "end_time": "2025-01-08T20:11:30.200+07:00",
            "duration": 0.2,
            "input": {...},
            "output": {...},
            "span_id": "1",
        }
        {
            "function": "usecase.CreateUser", // package.function
            "start_time": "2025-01-08T20:11:30.000+07:00",
            "end_time": "2025-01-08T20:11:30.200+07:00",
            "duration": 0.2,
            "input": {...},
            "output": {...},
            "span_id": "2",
        }
      // manual create by dev
        {
            "function": "usecase.CreateUser", // package.function
            "start_time": "2025-01-08T20:11:30.000+07:00",
            "end_time": "2025-01-08T20:11:30.200+07:00",
            "duration": 0.2,
            "input": {},
            "output": {},
            "envent": {
               "level": "debug", // info, warn, error, debug
               "message": "User id...."
            },
            "span_id": "3",
        }
        // manual create by system
        {
            "function": "repository.CreateUser", // package.function
            "start_time": "2025-01-08T20:11:30.000+07:00",
            "end_time": "2025-01-08T20:11:30.200+07:00",
            "duration": 0.2,
            "input": {...},
            "output": {...},
            "span_id": "4",
        }
        // spans อื่นๆ...
    ],

    // ข้อมูล Error (มีเมื่อ state เป็น error)
    "error": {
        "code": "USER_ALREADY_EXISTS",
        "message": "User with email already exists",
        "stack_trace": "...", // เฉพาะ development
        "details": {}
    }
}
```