package main

import "fmt"

// Define an interface with a method signature
type Greeter interface {
    Greet() string
}

// Define a struct that represents a Person
type Person struct {
    Name string
    Age  int
}

// Implement the Greeter interface for a pointer to a Person
func (p *Person) Greet() string {
    return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}

func printGreeting(g Greeter) {
    fmt.Println(g.Greet())
}

func main() {
    // Create a person instance
    alice := Person{Name: "Alice", Age: 30}

    // Passing a pointer to the interface function
    printGreeting(&alice)

    // Directly using the pointer to call Greet
    fmt.Println((&alice).Greet())

    // Demonstrate the pointer effect
    changeName(&alice, "Alice Wonderland")
    fmt.Println((&alice).Greet())
}

// Function that modifies the struct through a pointer
func changeName(p *Person, newName string) {
    p.Name = newName
}
