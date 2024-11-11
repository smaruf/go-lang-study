package main

import (
    "fmt"
    "math"
)

// Function type that takes a float64 and returns a float64
type operation func(float64) float64

// Function that applies an operation to a float64 value
func applyOperation(value float64, op operation) float64 {
    return op(value)
}

// Function to double the value
func double(value float64) float64 {
    return value * 2
}

// Function to calculate the square root of the value
func squareRoot(value float64) float64 {
    return math.Sqrt(value)
}

func main() {
    value := 16.0

    // Pass the double function as an argument
    doubled := applyOperation(value, double)
    fmt.Println("Doubled:", doubled)

    // Pass the squareRoot function as an argument
    sqrt := applyOperation(value, squareRoot)
    fmt.Println("Square root:", sqrt)
}
