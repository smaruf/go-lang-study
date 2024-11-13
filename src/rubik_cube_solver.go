// need to install `brew install cube-solver` and execute this by go
package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	scramble := "D2 R' U2 R U2 R' F R U R' F' R F U' F R' F' R U R"
	result, err := solveCube(scramble)
	if err != nil {
		fmt.Println("Error solving cube:", err)
		return
	}
	fmt.Println("Solution:", result)
}

func solveCube(scramble string) (string, error) {
	// Constructing the command
	cmd := exec.Command("cube-solver", scramble)
	
	// Capturing the output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Running the command
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
