package main

import "fmt"

func rectInfo(width, height float64) (area, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return
}

func main() {
    area, perimeter := rectInfo(4.0, 5.0)
    fmt.Println("Area:", area, "Perimeter:", perimeter)
}
