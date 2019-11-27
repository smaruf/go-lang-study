type Combination struct {
  given [45]int
  place int
}

func (c *Combination) helper(combinations [][]int, data []int, start int, end int, index int) {
    if (index == data.length) {
        add(combinations, clone(data))
    } else if (start <= end) {
        data[index] = start
        helper(combinations, data, start + 1, end, index + 1)
        helper(combinations, data, start + 1, end, index)
    }
}

func (c *Combination) add(combinations [][]int, data []int) {
    combinations[combinations.length] := clone(data)   
}

func (c *Combination) clone(data []int) []int {
    buffer [data.length]int
    for index :=0; index < data.length; index++ {
        buffer[index] := data[index]
    }
    return buffer
}

func (c *Combination) Generate(n int, r int) [][]int {
    combinations [][r]int
    data [r]int
    helper(combinations, data, 0, n-1, 0)
    return combinations
}
