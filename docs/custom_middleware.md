# Creating Custom Middleware

This guide explains how to create custom middleware for your preferred web framework using Goobserv.

## Core Concepts

When creating middleware, you need to handle these key aspects:

1. Context creation and management
2. Trace and request ID extraction
3. Span management
4. Request logging

## Example: Creating Custom Middleware

Here's a template for creating middleware for any web framework:

```go
func CreateMiddleware(obs *core.Observer, skipPaths []string) YourFrameworkMiddleware {
    return func(c YourFrameworkContext) {
        // 1. Check skip paths
        path := getPath(c) // Get path according to your framework
        for _, p := range skipPaths {
            if p == path {
                return next() // Skip middleware
            }
        }

        // 2. Extract trace and request IDs
        traceID := getHeader(c, "X-Trace-ID")
        requestID := getHeader(c, "X-Request-ID")

        // 3. Create context
        ctx := core.NewContext(context.Background()).
            WithTraceID(traceID).
            WithRequestID(requestID)

        // 4. Start span
        spanCtx := obs.StartSpan(ctx, fmt.Sprintf("%s %s", getMethod(c), path))
        
        // 5. Store context and observer
        setContext(c, "observContext", spanCtx)
        setContext(c, "observer", obs)

        // 6. Record metrics
        start := time.Now()
        defer func() {
            obs.Info(spanCtx, "HTTP Request",
                "method", getMethod(c),
                "path", path,
                "status", getStatus(c),
                "duration_ms", time.Since(start).Milliseconds(),
            )
            obs.EndSpan(spanCtx)
        }()

        // 7. Process request
        next()
    }
}
```

## Framework-Specific Examples

### Echo Framework

```go
func EchoMiddleware(obs *core.Observer, skipPaths []string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            path := c.Path()
            for _, p := range skipPaths {
                if p == path {
                    return next(c)
                }
            }

            traceID := c.Request().Header.Get("X-Trace-ID")
            requestID := c.Request().Header.Get("X-Request-ID")

            ctx := core.NewContext(c.Request().Context()).
                WithTraceID(traceID).
                WithRequestID(requestID)

            spanCtx := obs.StartSpan(ctx, fmt.Sprintf("%s %s", c.Request().Method, path))
            
            c.Set("observContext", spanCtx)
            c.Set("observer", obs)

            start := time.Now()
            defer func() {
                obs.Info(spanCtx, "HTTP Request",
                    "method", c.Request().Method,
                    "path", path,
                    "status", c.Response().Status,
                    "duration_ms", time.Since(start).Milliseconds(),
                )
                obs.EndSpan(spanCtx)
            }()

            return next(c)
        }
    }
}
```

### Chi Framework

```go
func ChiMiddleware(obs *core.Observer, skipPaths []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            path := r.URL.Path
            for _, p := range skipPaths {
                if p == path {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            traceID := r.Header.Get("X-Trace-ID")
            requestID := r.Header.Get("X-Request-ID")

            ctx := core.NewContext(r.Context()).
                WithTraceID(traceID).
                WithRequestID(requestID)

            spanCtx := obs.StartSpan(ctx, fmt.Sprintf("%s %s", r.Method, path))
            
            // Use context to store values
            r = r.WithContext(context.WithValue(r.Context(), "observContext", spanCtx))
            r = r.WithContext(context.WithValue(r.Context(), "observer", obs))

            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
            start := time.Now()
            defer func() {
                obs.Info(spanCtx, "HTTP Request",
                    "method", r.Method,
                    "path", path,
                    "status", ww.Status(),
                    "duration_ms", time.Since(start).Milliseconds(),
                )
                obs.EndSpan(spanCtx)
            }()

            next.ServeHTTP(ww, r)
        })
    }
}
```

## Helper Functions

We provide these helper functions to make it easier to work with the observer:

```go
// GetContext retrieves observer context from your framework's context
func GetContext(c YourFrameworkContext) *core.Context {
    if ctx := getFromContext(c, "observContext"); ctx != nil {
        if obsCtx, ok := ctx.(*core.Context); ok {
            return obsCtx
        }
    }
    return core.NewContext(context.Background())
}

// GetObserver retrieves observer from your framework's context
func GetObserver(c YourFrameworkContext) *core.Observer {
    if obs := getFromContext(c, "observer"); obs != nil {
        if observer, ok := obs.(*core.Observer); ok {
            return observer
        }
    }
    return nil
}
```

## Best Practices

1. Always handle skip paths to allow bypassing the middleware
2. Extract trace and request IDs from headers
3. Create a new context for each request
4. Use spans to track request lifecycle
5. Record request metrics (method, path, status, duration)
6. Clean up resources using defer
7. Provide helper functions for context retrieval

## Testing

Remember to test your middleware thoroughly:

1. Test skip paths functionality
2. Test context creation and retrieval
3. Test span management
4. Test metric recording
5. Test with various request scenarios
6. Test helper functions

See our Gin and Fiber middleware implementations for complete examples of testing.
