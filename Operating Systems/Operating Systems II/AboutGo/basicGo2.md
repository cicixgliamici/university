# Getting Started with Go Programming

This document introduces basic Go programming concepts and examples, including a simple "Hello, World!" program and basic mathematical functions.

## Hello, World!

A "Hello, World!" program is often the first program you write when learning a new language. Here’s how it looks in Go:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Explanation:
- `package main`: Every Go program starts with a `main` package. This indicates the entry point for the program.
- `import "fmt"`: Imports the `fmt` package, which includes functions for formatted I/O, such as `Println`.
- `func main()`: The `main` function is the starting point of the program.
- `fmt.Println`: Prints the message to the console.

Run the program using:
```sh
go run hello.go
```

## Basic Mathematical Functions

Let’s explore how to write and use basic mathematical functions in Go.

### Example: Adding Two Numbers

```go
package main

import "fmt"

func add(a int, b int) int {
    return a + b
}

func main() {
    result := add(5, 7)
    fmt.Println("The sum is:", result)
}
```

### Example: Factorial Function

```go
package main

import "fmt"

func factorial(n int) int {
    if n == 0 {
        return 1
    }
    return n * factorial(n-1)
}

func main() {
    fmt.Println("Factorial of 5 is:", factorial(5))
}
```

### Example: Checking if a Number is Even or Odd

```go
package main

import "fmt"

func isEven(n int) bool {
    return n%2 == 0
}

func main() {
    number := 8
    if isEven(number) {
        fmt.Println(number, "is even.")
    } else {
        fmt.Println(number, "is odd.")
    }
}
```
