package main

import "fmt"

func main() {
    var x int = 10     // Declare an integer variable
    var p *int = &x    // Declare a pointer to an integer and assign it the address of x

    fmt.Println("Value of x:", x)          // Output: Value of x: 10
    fmt.Println("Address of x:", &x)       // Output: Address of x: 0x...
    fmt.Println("Value of p:", p)          // Output: Value of p: 0x...
    fmt.Println("Value pointed to by p:", *p) // Output: Value pointed to by p: 10

    *p = 20    // Modify the value at the memory address pointed to by p
    fmt.Println("New value of x:", x)  // Output: New value of x: 20
}
