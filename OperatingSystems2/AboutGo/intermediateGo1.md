# Intermediate Go Concepts - Part 1: Structs, Interfaces, and Data Structures

## Structs and Methods

### Defining Structs
Structs are composite types that group fields of different data types.

```go
type Person struct {
    Name string
    Age  int
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    fmt.Println(p.Name) // Output: Alice
}
```

### Methods
Methods are functions with a receiver argument bound to a struct.

```go
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", p.Name)
}

func main() {
    p := Person{"Bob", 25}
    fmt.Println(p.Greet()) // Output: Hello, I'm Bob
}
```

## Interfaces
Interfaces define behavior contracts implicitly implemented by types.

```go
type Shape interface {
    Area() float64
}

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func printArea(s Shape) {
    fmt.Println("Area:", s.Area())
}

func main() {
    printArea(Circle{Radius: 5}) // Output: Area: 78.5398...
}
```

## Working with Slices
Dynamic sequences with flexible operations.

```go
func main() {
    // Create and modify slices
    nums := []int{1, 2, 3}
    nums = append(nums, 4, 5)
    
    // Slice operations
    sub := nums[1:3] // [2, 3]
    
    // Iteration
    for i, v := range nums {
        fmt.Printf("Index %d: %d\n", i, v)
    }
}
```

## Maps
Key-value collections for efficient lookups.

```go
func main() {
    colors := map[string]string{
        "red":   "#ff0000",
        "green": "#00ff00",
    }
    
    // Add/delete entries
    colors["blue"] = "#0000ff"
    delete(colors, "red")
    
    // Check existence
    if hex, exists := colors["green"]; exists {
        fmt.Println("Green hex:", hex) // #00ff00
    }
}
```
