# Introduction to Java and Object-Oriented Programming (OOP)

## What is Java?
- **Java** is a high-level, class-based, object-oriented programming language designed for platform independence.  
- Key features:  
  - **Write Once, Run Anywhere (WORA)**: Compiled Java code runs on any device with a Java Virtual Machine (JVM).  
  - **Automatic Memory Management**: Uses garbage collection to handle memory allocation/deallocation.  
  - **Strongly Typed**: Variables must be declared with specific data types.  
  - **Rich Standard Library**: Provides pre-built modules for I/O, networking, data structures, and more.  

---

## Core Concepts of Object-Oriented Programming (OOP)

### 1. **Classes and Objects**  
- **Class**: A blueprint/template for creating objects (e.g., `Car`).  
- **Object**: An instance of a class (e.g., `Car myCar = new Car();`).  

### 2. **Four Pillars of OOP**  
| Pillar          | Description                                                                 | Java Example                          |  
|-----------------|-----------------------------------------------------------------------------|---------------------------------------|  
| **Encapsulation** | Bundling data (fields) and methods into a single unit (class). Restricts direct access using `private`/`public` modifiers. | `private int speed; public int getSpeed() { return speed; }` |  
| **Inheritance**   | Deriving a new class (`subclass`) from an existing class (`superclass`). Promotes code reuse. | `class ElectricCar extends Car { ... }` |  
| **Polymorphism**  | Objects of different classes treated as objects of a common superclass. Achieved via method overriding/interfaces. | `@Override void startEngine() { ... }` |  
| **Abstraction**   | Hiding complex implementation details. Achieved via abstract classes/interfaces. | `abstract class Vehicle { abstract void move(); }` |  

---

## Why Java for OOP?  
- **Built for OOP**: Everything in Java is a class/object (except primitives like `int`).  
- **Modularity**: Code is organized into reusable classes and packages.  
- **Security**: Access modifiers (`public`, `private`, `protected`) enforce encapsulation.  
- **Multi-threading**: Native support for concurrent programming.  

---

## Basic Java Syntax Example  
```java
// Class definition
public class Dog {
    // Fields (encapsulated)
    private String name;
    
    // Constructor
    public Dog(String name) {
        this.name = name;
    }
    
    // Method
    public void bark() {
        System.out.println(name + " says: Woof!");
    }
}

// Creating an object
public class Main {
    public static void main(String[] args) {
        Dog myDog = new Dog("Buddy");
        myDog.bark(); // Output: Buddy says: Woof!
    }
}
```
## Key Takeaways
- Java’s OOP structure promotes clean, reusable, and scalable code.
- Mastering classes, objects, and the four OOP pillars is essential for Java development.
- Modern Java frameworks (Spring, Hibernate) and Android development rely heavily on OOP principles
