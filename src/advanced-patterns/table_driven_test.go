package main

import (
	"testing"
)

// Example 1: Basic Table-Driven Tests

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed numbers", -2, 3, 1},
		{"zeros", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Example 2: String Manipulation Tests

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", "olleh"},
		{"palindrome", "racecar", "racecar"},
		{"empty string", "", ""},
		{"single char", "a", "a"},
		{"with spaces", "hello world", "dlrow olleh"},
		{"unicode", "こんにちは", "はちにんこ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Example 3: Error Handling Tests

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

var ErrDivisionByZero = &DivisionError{"division by zero"}

type DivisionError struct {
	msg string
}

func (e *DivisionError) Error() string {
	return e.msg
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a         float64
		b         float64
		expected  float64
		expectErr bool
	}{
		{"normal division", 10.0, 2.0, 5.0, false},
		{"division by zero", 10.0, 0.0, 0.0, true},
		{"negative numbers", -10.0, 2.0, -5.0, false},
		{"decimal result", 10.0, 3.0, 3.3333333333333335, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			
			if tt.expectErr {
				if err == nil {
					t.Errorf("Divide(%f, %f) expected error, got nil", tt.a, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%f, %f) unexpected error: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("Divide(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

// Example 4: Struct Validation Tests

type User struct {
	Name  string
	Email string
	Age   int
}

func (u *User) Validate() error {
	if u.Name == "" {
		return &ValidationError{"name cannot be empty"}
	}
	if u.Email == "" {
		return &ValidationError{"email cannot be empty"}
	}
	if u.Age < 0 || u.Age > 150 {
		return &ValidationError{"age must be between 0 and 150"}
	}
	return nil
}

type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string {
	return e.msg
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		expectErr bool
		errMsg    string
	}{
		{
			name:      "valid user",
			user:      User{Name: "John", Email: "john@example.com", Age: 30},
			expectErr: false,
		},
		{
			name:      "empty name",
			user:      User{Name: "", Email: "john@example.com", Age: 30},
			expectErr: true,
			errMsg:    "name cannot be empty",
		},
		{
			name:      "empty email",
			user:      User{Name: "John", Email: "", Age: 30},
			expectErr: true,
			errMsg:    "email cannot be empty",
		},
		{
			name:      "negative age",
			user:      User{Name: "John", Email: "john@example.com", Age: -1},
			expectErr: true,
			errMsg:    "age must be between 0 and 150",
		},
		{
			name:      "age too high",
			user:      User{Name: "John", Email: "john@example.com", Age: 200},
			expectErr: true,
			errMsg:    "age must be between 0 and 150",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Example 5: Benchmarking with Table-Driven Tests

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func BenchmarkFibonacci(b *testing.B) {
	tests := []struct {
		name  string
		input int
	}{
		{"Fib(5)", 5},
		{"Fib(10)", 10},
		{"Fib(15)", 15},
		{"Fib(20)", 20},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Fibonacci(tt.input)
			}
		})
	}
}

// Example 6: Parallel Tests

func ExpensiveOperation(n int) int {
	result := 0
	for i := 0; i < n; i++ {
		result += i
	}
	return result
}

func TestExpensiveOperation(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"small", 10, 45},
		{"medium", 100, 4950},
		{"large", 1000, 499500},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel
			result := ExpensiveOperation(tt.input)
			if result != tt.expected {
				t.Errorf("ExpensiveOperation(%d) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}
