package main

import "testing"

func TestGCD_Calculate(t *testing.T) {
	tests := []struct {
		name     string
		first    int64
		second   int64
		expected int64
	}{
		{"Basic case", 48, 18, 6},
		{"Same numbers", 42, 42, 42},
		{"One is multiple of other", 15, 5, 5},
		{"Prime numbers", 17, 13, 1},
		{"Large numbers", 1071, 462, 21},
		{"Zero case", 10, 0, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcd := GCD{first: tt.first, second: tt.second}
			result := gcd.Calculate()
			if result != tt.expected {
				t.Errorf("GCD.Calculate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func BenchmarkGCD_Calculate(b *testing.B) {
	gcd := GCD{first: 1071, second: 462}
	for i := 0; i < b.N; i++ {
		gcd.Calculate()
	}
}