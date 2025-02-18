# Data Structures in Go

This markdown file provides an overview of some of the fundamental data structures in the Go programming language, including examples with detailed comments. You can use this file as a study guide and reference for understanding how these data structures work in Go.

---

## Table of Contents
- [Arrays](#arrays)
- [Slices](#slices)
- [Maps](#maps)
- [Structs](#structs)
- [Custom Data Structures](#custom-data-structures)
  - [Linked List](#linked-list)
  - [Stack](#stack)
  - [Queue](#queue)
- [Conclusion](#conclusion)

---

## Arrays

An array in Go is a fixed-size collection of elements of the same type.

```go
package main

import "fmt"

func main() {
    // Declare an array of 5 integers. The size is fixed and cannot be changed.
    var arr [5]int

    // Initialize the array with values.
    arr = [5]int{1, 2, 3, 4, 5}

    // Print the entire array.
    fmt.Println("Array:", arr)

    // Accessing array elements by index (indexes start at 0).
    fmt.Println("First element:", arr[0])
    fmt.Println("Last element:", arr[len(arr)-1])
}
```

**Explanation:**
- **Fixed Size:** Arrays have a fixed size defined at compile time.
- **Type Specific:** The size is part of the array's type, so `[5]int` and `[10]int` are considered different types.
- **Usage:** Arrays are useful when the number of elements is known and constant.

---

## Slices

Slices are more flexible than arrays. They are built on top of arrays and provide dynamic sizing.

```go
package main

import "fmt"

func main() {
    // Create a slice of integers using a slice literal.
    slice := []int{10, 20, 30, 40, 50}

    // Print the slice.
    fmt.Println("Slice:", slice)

    // Append a new element to the slice. The underlying array may be reallocated.
    slice = append(slice, 60)
    fmt.Println("Slice after append:", slice)

    // Slicing a slice: create a sub-slice from index 1 to 3 (3 is non-inclusive).
    subSlice := slice[1:4]
    fmt.Println("Sub-slice:", subSlice)

    // Slices have both a length and a capacity.
    fmt.Println("Length:", len(slice))
    fmt.Println("Capacity:", cap(slice))
}
```

**Explanation:**
- **Dynamic Size:** Slices can grow or shrink as needed.
- **Underlying Array:** They reference an underlying array, and when capacity is exceeded, a new array is allocated.
- **Convenience:** Slices are the most commonly used data structure for collections in Go.

---

## Maps

Maps are unordered collections of key-value pairs, implemented as hash tables.

```go
package main

import "fmt"

func main() {
    // Create a map with string keys and integer values.
    myMap := make(map[string]int)

    // Insert key-value pairs.
    myMap["apple"] = 5
    myMap["banana"] = 10
    myMap["cherry"] = 7

    // Print the entire map.
    fmt.Println("Map:", myMap)

    // Retrieve a value by key.
    appleCount := myMap["apple"]
    fmt.Println("Apple count:", appleCount)

    // Check if a key exists.
    if value, exists := myMap["banana"]; exists {
        fmt.Println("Banana count:", value)
    } else {
        fmt.Println("Banana not found")
    }

    // Delete a key-value pair.
    delete(myMap, "cherry")
    fmt.Println("Map after deletion:", myMap)
}
```

**Explanation:**
- **Unordered:** Maps do not maintain any order.
- **Fast Lookups:** They provide efficient retrieval of values by keys.
- **Safety:** Use the two-value assignment to safely check for key existence.

---

## Structs

Structs in Go allow you to create complex data types by grouping together fields.

```go
package main

import "fmt"

// Define a struct to represent a Person.
type Person struct {
    Name string  // Name of the person
    Age  int     // Age of the person
}

func main() {
    // Create an instance of Person using a struct literal.
    person := Person{
        Name: "Alice",
        Age:  30,
    }

    // Access struct fields using dot notation.
    fmt.Println("Person Name:", person.Name)
    fmt.Println("Person Age:", person.Age)

    // Modify a struct field.
    person.Age = 31
    fmt.Println("Updated Age:", person.Age)
}
```

**Explanation:**
- **Composite Type:** Structs allow you to combine multiple fields into a single type.
- **Initialization:** Use struct literals for clear and concise initialization.
- **No Inheritance:** Unlike some object-oriented languages, structs do not support inheritance, but you can achieve similar patterns with composition.

---

## Custom Data Structures

Go also allows you to implement more complex data structures. Below are examples of a linked list, a stack, and a queue.

### Linked List

A simple implementation of a singly linked list.

```go
package main

import "fmt"

// Node represents a single element in a linked list.
type Node struct {
    Value int   // Data stored in the node
    Next  *Node // Pointer to the next node
}

// printList traverses the linked list and prints its elements.
func printList(head *Node) {
    current := head
    for current != nil {
        fmt.Printf("%d -> ", current.Value)
        current = current.Next
    }
    fmt.Println("nil")
}

func main() {
    // Create nodes.
    node1 := &Node{Value: 1}
    node2 := &Node{Value: 2}
    node3 := &Node{Value: 3}

    // Link the nodes.
    node1.Next = node2
    node2.Next = node3

    // Print the linked list starting from the head node.
    printList(node1)
}
```

**Explanation:**
- **Nodes:** Each node contains data and a pointer to the next node.
- **Traversal:** The list is traversed by following the pointers from one node to the next.
- **Simplicity:** This is a basic example; more advanced linked lists may include features like double links or tail pointers.

---

### Stack

A stack follows a Last-In-First-Out (LIFO) order.

```go
package main

import "fmt"

// Stack represents a simple stack data structure using a slice.
type Stack struct {
    elements []int
}

// Push adds an element to the top of the stack.
func (s *Stack) Push(value int) {
    s.elements = append(s.elements, value)
}

// Pop removes and returns the top element of the stack.
// It returns false if the stack is empty.
func (s *Stack) Pop() (int, bool) {
    if len(s.elements) == 0 {
        // Stack is empty.
        return 0, false
    }
    // Get the last element.
    index := len(s.elements) - 1
    value := s.elements[index]
    // Remove the element from the slice.
    s.elements = s.elements[:index]
    return value, true
}

func main() {
    // Create a new stack.
    stack := Stack{}

    // Push elements onto the stack.
    stack.Push(10)
    stack.Push(20)
    stack.Push(30)

    fmt.Println("Stack after pushes:", stack.elements)

    // Pop an element from the stack.
    if value, ok := stack.Pop(); ok {
        fmt.Println("Popped:", value)
    }
    fmt.Println("Stack after pop:", stack.elements)
}
```

**Explanation:**
- **LIFO:** The last element added is the first one to be removed.
- **Operations:** `Push` adds an element to the top, and `Pop` removes it.
- **Implementation:** Using a slice makes it simple to implement dynamic behavior.

---

### Queue

A queue follows a First-In-First-Out (FIFO) order.

```go
package main

import "fmt"

// Queue represents a simple queue data structure using a slice.
type Queue struct {
    elements []int
}

// Enqueue adds an element to the end of the queue.
func (q *Queue) Enqueue(value int) {
    q.elements = append(q.elements, value)
}

// Dequeue removes and returns the first element of the queue.
// It returns false if the queue is empty.
func (q *Queue) Dequeue() (int, bool) {
    if len(q.elements) == 0 {
        // Queue is empty.
        return 0, false
    }
    // Get the first element.
    value := q.elements[0]
    // Remove the first element by re-slicing.
    q.elements = q.elements[1:]
    return value, true
}

func main() {
    // Create a new queue.
    queue := Queue{}

    // Enqueue elements.
    queue.Enqueue(100)
    queue.Enqueue(200)
    queue.Enqueue(300)

    fmt.Println("Queue after enqueues:", queue.elements)

    // Dequeue an element from the queue.
    if value, ok := queue.Dequeue(); ok {
        fmt.Println("Dequeued:", value)
    }
    fmt.Println("Queue after dequeue:", queue.elements)
}
```

**Explanation:**
- **FIFO:** The first element added is the first one removed.
- **Operations:** `Enqueue` adds to the end and `Dequeue` removes from the front.
- **Implementation:** Using slices provides a simple way to manage dynamic collections.

---

## Conclusion

In this markdown file, we explored several common data structures in Go:

- **Arrays:** Fixed-size collections.
- **Slices:** Dynamic, flexible arrays built on top of fixed arrays.
- **Maps:** Unordered collections of key-value pairs, ideal for fast lookups.
- **Structs:** Custom composite types to group related data.
- **Custom Data Structures:** Examples of a linked list, stack, and queue to demonstrate more complex behaviors.
