package main 

import (
	"fmt"
	"os"
)


func main() { // main.go funktion never takes parameters 

	// if the inputs are not exactly 3, print a usage message
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . input.txt output.txt")
		return
	}

	

}