package main 

import (
	"fmt"
	"os"
	"strings"
)


func main() { // main.go funktion never takes parameters 

	// if the inputs are not exactly 3, print a usage message
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . input.txt output.txt")
		return
	}

	// reading a file
	fileContents, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Split the text into tokens (words) so each rule can edit them one by one,
	// then join them back with single spaces. strings.Fields also removes extra whitespace.
	text := string(fileContents)
	tokens := strings.Fields(text)
	result := strings.Join(tokens, " ")

	// writing the output file
	err = os.WriteFile(os.Args[2], []byte(result), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}