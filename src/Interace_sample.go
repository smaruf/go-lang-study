package main

import (
	"fmt"
)

// Define an interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Define a struct that implements the Shape interface
type Rectangle struct {
	Width, Height float64
}

// Implement the Area method for Rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Implement the Perimeter method for Rectangle
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Define another struct that implements the Shape interface
type Circle struct {
	Radius float64
}

// Implement the Area method for Circle
func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

// Implement the Perimeter method for Circle
func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.Radius
}

// A function that takes a Shape interface and calls its methods
func PrintShapeInfo(s Shape) {
	fmt.Printf("Area: %f\n", s.Area())
	fmt.Printf("Perimeter: %f\n", s.Perimeter())
}

func main() {
	r := Rectangle{Width: 10, Height: 5}
	c := Circle{Radius: 7}

	fmt.Println("Rectangle:")
	PrintShapeInfo(r)

	fmt.Println("\nCircle:")
	PrintShapeInfo(c)
}
