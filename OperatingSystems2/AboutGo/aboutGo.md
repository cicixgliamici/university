# Introduction to Go (Golang)

Go is a statically typed, compiled programming language designed for simplicity, efficiency, and concurrency. It is widely used for building system-level software, web services, and distributed systems.

## Key Features of Go

1. **Concurrency**:
   - Go provides lightweight threads of execution called **goroutines**.
   - **Channels** are used for communication between goroutines, enabling safe and efficient data exchange.
2. **Static Typing and Simplicity**:
   - The language enforces type safety while maintaining a simple and readable syntax.
3. **Garbage Collection**:
   - Automatic memory management eliminates the need for manual allocation and deallocation.
4. **Built-in Tooling**:
   - Go includes tools for formatting, testing, benchmarking, and more.
5. **Rich Standard Library**:
   - Provides built-in packages for networking, cryptography, HTTP, and more.

## Declaring Variables

In Go, variables can be declared using the `var` keyword or shorthand syntax.

Using the 'var' keyword
`var message string = "Hello, World!"`

Type inference
`var number = 42`

Shorthand syntax
`count := 10`
