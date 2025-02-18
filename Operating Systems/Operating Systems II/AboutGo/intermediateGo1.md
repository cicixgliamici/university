# Intermediate Go Concepts - Part 1: Structs, Interfaces, and Data Structures

## Structs and Methods

### Defining Structs
Structs are composite types that group together fields of different data types, similar to classes in other languages but without inheritance.

```go
// Person struct represents a person with name and age
type Person struct {
    Name string // Public field (starts with uppercase)
    Age  int    // Public field
}

func main() {
    // Struct initialization with field names (explicit)
    p1 := Person{Name: "Alice", Age: 30}
    
    // Struct initialization without field names (order must match declaration)
    p2 := Person{"Bob", 25}
    
    // Accessing struct fields
    fmt.Println(p1.Name) // Output: Alice
    fmt.Println(p2.Age)  // Output: 25

    // Anonymous structs (useful for one-time use)
    temp := struct {
        ID   int
        Desc string
    }{1001, "Temporary item"}
    
    fmt.Println(temp) // Output: {1001 Temporary item}
}
```

### Methods
Methods are functions with a receiver argument that binds them to a specific type. They can operate on either value or pointer receivers.

```go
// Value receiver (works on a copy of the struct)
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s (age %d)", p.Name, p.Age)
}

// Pointer receiver (can modify the original struct)
func (p *Person) Birthday() {
    p.Age++ // Modifies the actual struct
}

func main() {
    person := Person{"Charlie", 40}
    
    fmt.Println(person.Greet()) // Output: Hello, I'm Charlie (age 40)
    
    person.Birthday()
    fmt.Println(person.Age) // Output: 41 (age modified by pointer receiver)
    
    // Methods can be called on both values and pointers
    (&person).Birthday()
    fmt.Println(person.Age) // Output: 42
}
```

## Interfaces
Interfaces define behavior through method signatures. Types implicitly implement interfaces by implementing all their methods.

```go
// Shape interface requires Area() method
type Shape interface {
    Area() float64
    Perimeter() float64 // Adding another method requirement
}

// Circle implements Shape interface
type Circle struct {
    Radius float64
}

// Area calculation for Circle
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Perimeter calculation for Circle
func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Rectangle also implements Shape
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// printShapeInfo demonstrates interface usage
func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    shapes := []Shape{
        Circle{Radius: 5},
        Rectangle{Width: 3, Height: 4},
    }
    
    for _, shape := range shapes {
        printShapeInfo(shape)
    }
    // Output:
    // Area: 78.54, Perimeter: 31.42
    // Area: 12.00, Perimeter: 14.00
}
```

## Working with Slices
Slices are dynamic, flexible views into arrays. They grow automatically as needed.

```go
func main() {
    // Creating slices
    emptySlice := make([]int, 0, 5) // Type, length (0), capacity (5)
    numbers := []int{1, 2, 3}       // Literal initialization
    
    // Modifying slices
    numbers = append(numbers, 4)    // Add single element
    numbers = append(numbers, 5, 6) // Add multiple elements
    
    // Slice operations
    firstTwo := numbers[:2]   // [1, 2] (inclusive start, exclusive end)
    lastThree := numbers[3:]  // [4, 5, 6]
    middle := numbers[2:4]    // [3, 4]
    
    // Important: Slices share underlying array!
    middle[0] = 100
    fmt.Println(numbers) // [1 2 100 4 5 6]
    
    // Safe copy using copy() function
    copied := make([]int, len(middle))
    copy(copied, middle)
    copied[0] = 200
    fmt.Println(numbers) // Unchanged: [1 2 100 4 5 6]
    
    // Iteration with index and value
    fmt.Println("Slice contents:")
    for index, value := range numbers {
        fmt.Printf("Index %d: %d\n", index, value)
    }
    
    // Slice capacity demonstration
    fmt.Printf("Length: %d, Capacity: %d\n", len(numbers), cap(numbers))
    numbers = append(numbers, 7, 8, 9)
    fmt.Printf("New Capacity: %d\n", cap(numbers)) // Capacity doubles when exceeded
}
```

## Maps
Maps are unordered key-value collections with fast lookups. Keys must be comparable types.

```go
func main() {
    // Map declaration and initialization
    colorCodes := map[string]string{
        "red":   "#ff0000",
        "green": "#00ff00",
        "blue":  "#0000ff",
    }
    
    // Adding new entries
    colorCodes["white"] = "#ffffff"
    colorCodes["black"] = "#000000"
    
    // Deleting entries
    delete(colorCodes, "green")
    
    // Checking existence
    if code, exists := colorCodes["red"]; exists {
        fmt.Println("Red code:", code) // #ff0000
    }
    
    // Safe access for non-existing keys
    fmt.Println("Yellow code:", colorCodes["yellow"]) // Returns zero value: ""
    
    // Iteration (order not guaranteed!)
    fmt.Println("\nAll colors:")
    for color, code := range colorCodes {
        fmt.Printf("%-6s -> %s\n", color, code)
    }
    
    // Map as counter
    wordCount := make(map[string]int)
    text := "the quick brown fox jumps over the lazy dog"
    
    for _, word := range strings.Fields(text) {
        wordCount[word]++
    }
    
    fmt.Println("\nWord counts:")
    fmt.Println(wordCount) // map[brown:1 dog:1 fox:1 jumps:1 lazy:1 over:1 quick:1 the:2 ...]
    
    // Note: Maps are reference types!
    modified := colorCodes
    modified["red"] = "#ff0001"
    fmt.Println("Original red:", colorCodes["red"]) // #ff0001
}
```

## Best Practices and Notes

1. **Struct Design**:
   - Use constructor functions for complex initialization
   ```go
   func NewPerson(name string, age int) *Person {
       return &Person{name, age}
   }
   ```
   - Keep related methods in the same file as the struct

2. **Interface Tips**:
   - Prefer small, focused interfaces (io.Reader/Writer are good examples)
   - Use interface{} (empty interface) sparingly for maximum flexibility

3. **Slice Performance**:
   - Preallocate capacity with make() when size is known
   - Reuse slices when possible to reduce garbage collection

4. **Map Considerations**:
   - Use struct{} as value type for sets: `set := make(map[string]struct{})`
   - Not safe for concurrent use - sync.Map for concurrent access
