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
	tokens = convertBases(tokens) // our function
	result := strings.Join(tokens, " ")

	// writing the output file
	err = os.WriteFile(os.Args[2], []byte(result), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}




// Hex is base 16 and bin is base 2
func convertBases(tokens []string) []string {
	for index, word := range tokens {

		// hex
		if word == "(hex)" && index > 0 {
			number, err := strconv.ParseInt(tokens[index-1], 16, 64) // The 16 tells ParseInt how to read the digits, stored in 64 bits.
			if err == nil { // err == nil: success, so we use the result
				tokens[index-1] = strconv.FormatInt(number, 10) // replace it. FormatInt converts the number back into text, written in base 10 (normal decimal)
			}
		}

		// bin
		if word == "(bin)" && index > 0 { 
			number, err := strconv.ParseInt(tokens[index-1], 2, 64) // read as base 2
			if err == nil { 
				tokens[index-1] = strconv.FormatInt(number, 10) // replace it.
			}
		}

	}


}