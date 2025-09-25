package main

import (
	"fmt"
	"os"
	"strconv"
)

// GCD struct for gcd generation
type GCD struct {
	first  int64
	second int64
}

// Calculate calculates the greatest common divisor using Euclidean algorithm
func (numbers *GCD) Calculate() int64 {
	// Handle special cases
	if numbers.first == 0 {
		return numbers.second
	}
	if numbers.second == 0 {
		return numbers.first
	}
	
	divisor := numbers.getMin()
	dividend := numbers.getMax()
	if dividend%divisor == 0 {
		return divisor
	}

	modulus := dividend % divisor
	newNumbers := GCD{modulus, divisor}
	return newNumbers.Calculate()
}

func (numbers *GCD) getMin() int64 {
	if numbers.first < numbers.second {
		return numbers.first
	}
	return numbers.second
}

func (numbers *GCD) getMax() int64 {
	if numbers.first > numbers.second {
		return numbers.first
	}
	return numbers.second
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run calculate_gcd.go <first_number> <second_number>")
		fmt.Println("Example: go run calculate_gcd.go 48 18")
		os.Exit(1)
	}

	first, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing first number: %v\n", err)
		os.Exit(1)
	}

	second, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing second number: %v\n", err)
		os.Exit(1)
	}

	gcd := GCD{first: first, second: second}
	result := gcd.Calculate()

	fmt.Printf("GCD of %d and %d is: %d\n", first, second, result)
}
