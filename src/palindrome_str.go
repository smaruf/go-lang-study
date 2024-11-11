package main

import (
	"fmt"
	"strings"
)

// isPalindrome checks if a given string is a palindrome.
func isPalindrome(s string) bool {
    // Normalize the string: Remove spaces and convert to lower case.
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToLower(s)

	// Initialize two pointers, one at the start and one at the end of the string.
	i, j := 0, len(s)-1

	// While the two pointers haven't met in the middle,
	// check if the characters at the pointers are the same.
	while i < j {
		if s[i] != s[j] {
			return false // If they differ, it's not a palindrome.
		}
		i++ // Move the start pointer towards the center.
		j-- // Move the end pointer towards the center.
	}

	// If all characters matched correctly, it's a palindrome.
	return true
}

func main() {
	testStrings := []string{"racecar", "hello", "A man a plan a canal Panama"}
	for _, str := range testStrings {
		fmt.Printf("Is '%s' a palindrome? %v\n", str, isPalindrome(str))
	}
}
