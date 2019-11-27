package main

// "Combination structure for generate combinations"
type Combination struct {
  n int
  r int
}

func (c *Combination) helper(combinations [][]int, data []int, start int, end int, index int) {
    if (index == len(data)) {
        c.add(combinations, c.clone(data))
    } else if (start <= end) {
        data[index] = start
        c.helper(combinations, data, start + 1, end, index + 1)
        c.helper(combinations, data, start + 1, end, index)
    }
}

func (c *Combination) add(combinations [][]int, data []int) {
    append(combinations, c.clone(data))   
}

func (c *Combination) clone(data []int) []int {
    buffer := make([]int, len(data))
    for index := 0; index < len(data); index++ {
        buffer[index] = data[index]
    }
    return buffer
}

// "Generate combination of numbers"
func (c *Combination) Generate() [][]int { 
    combinations := [][]int {}
    data := make([]int, c.r)
    c.helper(combinations, data, 0, (c.n - 1), 0)
    return combinations
}
