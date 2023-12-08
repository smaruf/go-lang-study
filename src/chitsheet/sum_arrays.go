package main

import "fmt"

func sum(numbers []int) int {
 total := 0
 for _, number := range numbers {
 total += number
 }
 return total
}

func main() {
 numbers := []int{1, 2, 3, 4, 5}
 result := sum(numbers)
 fmt.Println("Sum:", result)
}
