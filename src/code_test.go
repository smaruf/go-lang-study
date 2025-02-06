package main

import (
	"fmt"
	"sort"
	"strings"
)

// Sorting Algorithms
func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}

func insertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j = j - 1
		}
		arr[j+1] = key
	}
}

// Matrix Manipulations
func rotateMatrix(matrix [][]int) [][]int {
	n := len(matrix)
	for x := 0; x < n/2; x++ {
		for y := x; y < n-x-1; y++ {
			temp := matrix[x][y]
			matrix[x][y] = matrix[y][n-1-x]
			matrix[y][n-1-x] = matrix[n-1-x][n-1-y]
			matrix[n-1-x][n-1-y] = matrix[n-1-y][x]
			matrix[n-1-y][x] = temp
		}
	}
	return matrix
}

// String Manipulations
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func isPalindrome(s string) bool {
	s = strings.ToLower(strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s))
	return s == reverseString(s)
}

// Searching Algorithms
func binarySearch(arr []int, target int) int {
	low, high := 0, len(arr)-1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] < target {
			low = mid + 1
		} else if arr[mid] > target {
			high = mid - 1
		} else {
			return mid
		}
	}
	return -1
}

// Set Manipulations
func union(set1, set2 map[int]bool) map[int]bool {
	result := make(map[int]bool)
	for key := range set1 {
		result[key] = true
	}
	for key := range set2 {
		result[key] = true
	}
	return result
}

func intersection(set1, set2 map[int]bool) map[int]bool {
	result := make(map[int]bool)
	for key := range set1 {
		if set2[key] {
			result[key] = true
		}
	}
	return result
}

func main() {
	// Example usage of bubble sort
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	bubbleSort(arr)
	fmt.Println("Sorted Array (Bubble Sort):", arr)

	// Example usage of matrix rotation
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	rotatedMatrix := rotateMatrix(matrix)
	fmt.Println("Rotated matrix:")
	for _, row := range rotatedMatrix {
		fmt.Println(row)
	}

	// Example usage of string manipulations
	fmt.Println("Reversed String:", reverseString("hello"))
	fmt.Println("Is Palindrome:", isPalindrome("A man, a plan, a canal, Panama"))

	// Example usage of set manipulations
	set1 := map[int]bool{1: true, 2: true, 3: true}
	set2 := map[int]bool{2: true, 3: true, 4: true}
	fmt.Println("Union of sets:", union(set1, set2))
	fmt.Println("Intersection of sets:", intersection(set1, set2))
}
