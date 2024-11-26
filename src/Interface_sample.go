package main

import (
    "fmt"
    "math"
)

// Shape interface defines methods for shapes
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle struct implements the Shape interface
type Rectangle struct {
    Width, Height float64
}

// Area method calculates the area of a rectangle
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Perimeter method calculates the perimeter of a rectangle
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle struct implements the Shape interface
type Circle struct {
    Radius float64
}

// Area method calculates the area of a circle
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Perimeter method calculates the perimeter of a circle
func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// PrintShapeInfo prints the area and perimeter of a shape
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %f\n", s.Area())
    fmt.Printf("Perimeter: %f\n", s.Perimeter())
}

func main() {
    r := Rectangle{
