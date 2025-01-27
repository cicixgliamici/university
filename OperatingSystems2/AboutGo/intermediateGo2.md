# Intermediate Go Concepts - Part 2: Error Handling, Flow Control, and Project Management

## Error Handling
Explicit error checking with multiple return values.

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    if result, err := divide(10, 0); err != nil {
        fmt.Println("Error:", err) // Division by zero
    }
}
```

### Custom Errors
Implement the `error` interface for domain-specific errors.

```go
type APIError struct {
    Code    int
    Message string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

func fetchData() error {
    return &APIError{404, "Resource not found"}
}
```

## Defer, Panic, and Recover
Control flow mechanisms for resource management and emergencies.

```go
func readFile(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        panic("File error!")
    }
    defer file.Close() // Always executes last
    
    // Process file...
}

func handlePanic() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    panic("Unexpected problem!")
}
```

## Testing and Benchmarking
Native testing framework examples.

```go
// math_test.go
package main

import "testing"

func TestMultiply(t *testing.T) {
    if multiply(2, 3) != 6 {
        t.Error("Expected 2 * 3 = 6")
    }
}

func BenchmarkMultiply(b *testing.B) {
    for i := 0; i < b.N; i++ {
        multiply(10, 20)
    }
}
```

Run tests with:
```sh
go test -v    # Verbose tests
go test -bench .  # Benchmarks
```

## Package Management
Using Go Modules for dependency management.

```go
// go.mod
module github.com/yourusername/project

go 1.21

require (
    github.com/gorilla/mux v1.8.0
)
```

### Creating Packages
```go
// utils/math.go
package utils

func Square(n int) int {
    return n * n
}

// main.go
package main

import (
    "fmt"
    "github.com/yourusername/project/utils"
)

func main() {
    fmt.Println(utils.Square(5)) // 25
}
```

## Concurrency Patterns
Advanced channel usage with select and timeouts.

```go
func worker(ch chan string) {
    time.Sleep(2 * time.Second)
    ch <- "Work done"
}

func main() {
    ch := make(chan string)
    go worker(ch)
    
    select {
    case res := <-ch:
        fmt.Println(res)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```
