package main
import (
    "fmt"
)
func main() {
    list := []int{3, 4, 1, 5, 2}
    iList := insertionSort(list)
    mList := mergeSort(list)
    fmt.Println(iList, mList)
}
func insertionSort(list []int) []int {
    for i, num := range list {
        j := i - 1
        for j >= 0 && num < list[j] {
            list[j+1] = list[j]
            j -= 1
        }
        list[j+1] = num
    }
    return list
}
func mergeSort(list []int) []int {
    if len(list) > 1 {
        mid := len(list) / 2
        left := list[:mid]
        right := list[mid:]
        mergeSort(left)
        mergeSort(right)
        i, j, k := 0, 0, 0
        for i < len(left) && j < len(right) {
            if left[i] < right[j] {
                list[k] = left[i]
                i++
            } else {
                list[k] = right[j]
                j++
            }
            k++
        }
        for i < len(left) {
            list[k] = left[i]
            i++
            k++
        }
        for j < len(right) {
            list[k] = right[j]
            j++
            k++
        }
    }
    return list
}
