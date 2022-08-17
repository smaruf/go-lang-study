package main
import (
    "fmt"
    "math/rand"
    "testing"
)
func BenchmarkInsertionSort(b *testing.B) {
    inputSize := []int{10, 100, 1000, 10000, 100000}
    for _, size := range inputSize {
        b.Run(fmt.Sprintf("input_size_%d", size), func(b *testing.B) {
            testList := make([]int, size)
            for i := 0; i < size; i++ {
                testList[i] = rand.Intn(size)
            }
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                insertionSort(testList)
            }
        })
    }
}
func BenchmarkMergeSort(b *testing.B) {
    inputSize := []int{10, 100, 1000, 10000, 100000}
    for _, size := range inputSize {
        b.Run(fmt.Sprintf("input_size_%d", size), func(b *testing.B) {
            testList := make([]int, size)
            for i := 0; i < size; i++ {
                testList[i] = rand.Intn(size)
            }
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                mergeSort(testList)
            }
        })
    }
}
