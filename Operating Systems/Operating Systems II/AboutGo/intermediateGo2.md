# Intermediate Go Concepts - Part 2: Error Handling, Flow Control, and Project Management

## Error Handling in Depth
Go's explicit error handling encourages robust code through multiple return values.

```go
// Basic error handling pattern
func divide(a, b float64) (float64, error) {
    if b == 0 {
        // Return zero value and meaningful error
        return 0, fmt.Errorf("division by zero in divide(%f, %f)", a, b)
    }
    return a / b, nil
}

func main() {
    // Handle errors immediately after function call
    result, err := divide(10, 0)
    if err != nil {
        // Use %v for general errors, %+v for stack traces (if available)
        fmt.Printf("Operation failed: %v\n", err)
        return
    }
    fmt.Println("Result:", result)
    
    // Wrapping errors with context
    _, err = os.Open("nonexistent.txt")
    if err != nil {
        wrappedErr := fmt.Errorf("file open failed: %w", err)
        fmt.Println(wrappedErr) // Preserves original error
    }
}
```

### Advanced Error Handling Techniques
```go
// Custom error type with stack trace
import "runtime/debug"

type DetailedError struct {
    Code     int
    Message  string
    Stack    string
    Original error
}

func (e *DetailedError) Error() string {
    return fmt.Sprintf("[%d] %s (orig: %v)", e.Code, e.Message, e.Original)
}

func newDetailedError(code int, msg string, err error) *DetailedError {
    return &DetailedError{
        Code:     code,
        Message:  msg,
        Stack:    string(debug.Stack()),
        Original: err,
    }
}

// Error type checking
func handleError(err error) {
    switch e := err.(type) {
    case *APIError:
        fmt.Println("API Failure:", e.Code)
    case *DetailedError:
        fmt.Println("Debug info:", e.Stack)
    default:
        fmt.Println("Generic error:", err)
    }
}
```

## Mastering Defer, Panic, and Recover
Control flow tools for resource management and exceptional situations.

```go
// Defer execution order (LIFO)
func deferExample() {
    defer fmt.Println("First registered")
    defer fmt.Println("Second registered")
    fmt.Println("Main logic")
    // Output:
    // Main logic
    // Second registered
    // First registered
}

// Resource cleanup pattern
func writeFile(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close() // Ensure file is always closed
    
    writer := bufio.NewWriter(file)
    defer writer.Flush() // Flush before closing
    
    _, err = writer.WriteString("data")
    return err
}

// Panic/Recover in web server middleware
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Recovered from panic: %v", r)
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

## Comprehensive Testing Strategies
Go's built-in testing framework supports various testing approaches.

```go
// Table-driven tests
func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a       float64
        b       float64
        want    float64
        wantErr bool
    }{
        {"normal division", 6, 3, 2, false},
        {"divide by zero", 5, 0, 0, true},
        {"decimal result", 5, 2, 2.5, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := divide(tt.a, tt.b)
            if (err != nil) != tt.wantErr {
                t.Errorf("divide() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("divide() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Benchmark with setup
func BenchmarkComplexOperation(b *testing.B) {
    // Expensive setup
    config := loadConfig()
    b.ResetTimer() // Exclude setup time
    
    for i := 0; i < b.N; i++ {
        complexOperation(config)
    }
}

// TestMain for global setup/teardown
func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    teardown()
    os.Exit(code)
}
```

## Advanced Package Management
Go Modules best practices and patterns.

```go
// go.mod with version constraints
module github.com/yourproject/service

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.7 // indirect
)

replace github.com/local/module => ../local-module // Local development

// Private repository setup
require (
    git.example.com/internal/auth v0.1.0
)

// go get commands:
// - go get -u ./...          # Update all dependencies
// - go get package@version   # Specific version
// - go mod tidy              # Cleanup unused dependencies
```

### Professional Package Design
```go
// Internal package structure
// project/
// ├── internal/
// │   └── analytics/         # Only accessible to parent module
// ├── pkg/
// │   └── utils/             # Reusable public utilities
// └── cmd/
//     └── server/            # Main application

// Package documentation example
// Package utils provides general-purpose utility functions.
//
// Usage:
//  import "github.com/yourproject/pkg/utils"
package utils

// Square returns the square of an integer.
// Example:
//  utils.Square(5) // 25
func Square(n int) int {
    return n * n
}
```

## Production-Grade Concurrency Patterns
Advanced channel patterns for real-world applications.

```go
// Worker pool pattern
func workerPool() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go func(id int) {
            for job := range jobs {
                fmt.Printf("Worker %d processing job %d\n", id, job)
                results <- job * 2
            }
        }(w)
    }
    
    // Send jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Collect results
    for r := 1; r <= 5; r++ {
        fmt.Println("Result:", <-results)
    }
}

// Context-based cancellation
func longOperation(ctx context.Context) error {
    select {
    case <-time.After(10 * time.Second):
        return nil // Normal completion
    case <-ctx.Done():
        return ctx.Err() // Cancellation
    }
}

// Fan-in pattern
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    // Start output goroutine
    wg.Add(len(cs))
    for _, c := range cs {
        go func(c <-chan int) {
            for n := range c {
                out <- n
            }
            wg.Done()
        }(c)
    }
    
    // Close channel when all done
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

// Channel ownership guidelines
// 1. Owner creates the channel
// 2. Owner writes to channel or passes write capability
// 3. Owner closes the channel
// 4. Consumers only read from channel
```

## Best Practices and Pro Tips

1. **Error Handling**:
   - Use `errors.Is` and `errors.As` for error inspection
   - Wrap errors with context using `%w`
   - Create sentinel errors for API boundaries: `var ErrNotFound = errors.New("not found")`

2. **Concurrency**:
   - Use `sync.Once` for one-time initialization
   - Prefer `sync.WaitGroup` for coordinating goroutines
   - Use `atomic` operations for simple counters
   - Always run tests with `-race` flag: `go test -race ./...`

3. **Project Structure**:
   - Keep business logic in `internal` packages
   - Separate commands in `cmd` directory
   - Use `config` package for configuration loading
   - Implement `internal/mocks` for testing interfaces

4. **Performance**:
   - Preallocate slice capacity with `make([]T, 0, capacity)`
   - Use `strings.Builder` for string concatenation
   - Profile with `pprof`: 
     ```sh
     go tool pprof -http :8080 http://localhost:6060/debug/pprof/heap
     ```

5. **Tooling**:
   - Format code with `gofmt`
   - Lint with `golangci-lint`
   - Generate docs with `go doc`
   - Create binaries with `go build -ldflags="-s -w"` for smaller size
