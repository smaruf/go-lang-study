package main

import "fmt"

// Define an interface that includes a method for demonstrating the greeting.
type Greeter interface {
    Greet() string
}

// Define a struct that represents a Person.
type Person struct {
    Name string
    Age  int
}

// Implement the Greeter interface with a method that uses a pointer receiver.
func (p *Person) Greet() string {
    return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}

// Function that returns a closure, which is a nested function that can access and modify Person.
func personModifier(p *Person) func(string, int) {
    return func(newName string, newAge int) {
        p.Name = newName // modify struct field through pointer
        p.Age = newAge   // modify struct field through pointer
    }
}

func main() {
    alice := &Person{Name: "Alice", Age: 30} // Create an instance of Person

    fmt.Println(alice.Greet()) // Use the method that satisfies the interface

    // Get the closure that can modify alice.
    updateAlice := personModifier(alice)

    // Use the closure to modify alice.
    updateAlice("Alice Wonderland", 35)

    fmt.Println(alice.Greet()) // Check the updated details
}
