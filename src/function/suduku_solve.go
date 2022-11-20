package function

func findNextEmpty(board [][]byte) (byte, byte, error) {
	for j := byte(0); j < 9; j++ {
		for i := byte(0); i < 9; i++ {
			if board[j][i] == '.' {
				return j, i, nil
			}
		}
	}
	return byte(0), byte(0), errors.New("no empty cell to fill in")
}

func ValidRule(board [][]byte, j, i byte) bool {
	num := board[j][i]
	// Check if the number is valid in the row
	for k := byte(0); k < 9; k++ {
		if board[j][k] == num && k != i {
			return false
		}
	}

	// Check if the number is valid in the column
	for k := byte(0); k < 9; k++ {
		if board[k][i] == num && k != j {
			return false
		}
	}

	// Check if the number is valid in the 3x3 grid
	for k := byte(0); k < 3; k++ {
		for l := byte(0); l < 3; l++ {
			if board[j-j%3+k][i-i%3+l] == num && (j-j%3+k != j || i-i%3+l != i) {
				return false
			}
		}
	}
	return true
}

func solveSudoku(board [][]byte) {
	j, i, e := findNextEmpty(board)
	if e != nil {
		return
	}

	var backtracing func(byte, byte) bool
	backtracing = func(j, i byte) bool {

		// fill in avabilable choices
		for k := byte(1); k <= 9; k++ {
			board[j][i] = 48 + k
			if ValidRule(board, j, i) {
				j, i, e := findNextEmpty(board)
				if e != nil {
					return true
				}

				status := backtracing(j, i)
				if status {
					return true
				}
			}
			board[j][i] = '.'
		}

		return false
	}

	// start backtracking
	backtracing(j, i)
}
