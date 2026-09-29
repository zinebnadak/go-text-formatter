package main 

import (
	"fmt"
	"os"
	"strings"
	"strconv"
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

	text = strings.ReplaceAll(text, "(up, ", "(up,")   // for numbered tags. "(up, 2)" -> "(up,2)"
	text = strings.ReplaceAll(text, "(low, ", "(low,")
	text = strings.ReplaceAll(text, "(cap, ", "(cap,")

	tokens := strings.Fields(text)
	tokens = convertBases(tokens) // hex & bin
	tokens = changeCase(tokens) // up, low, cap
	tokens = removeTags(tokens) // remove tags ALWAYS RUNS AFTER
	tokens = fixPunctuation(tokens) // .,!?:; spacing
	tokens = fixQuotes(tokens) // ' ' quotes
	tokens = fixAns(tokens) // a -> an
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
	return tokens
}


func removeTags(tokens []string) []string {
	var keep []string 	// empty slice to fill 
	for _,  word := range tokens {
			if word == "(hex)" || word == "(bin)" || word == "(up)" || word == "(low)" || word == "(cap)" || strings.HasPrefix(word, "(up,") || strings.HasPrefix(word, "(low,") || strings.HasPrefix(word, "(cap,") {  // is it a tag???
			continue // skip if yes
		}
		keep = append(keep, word) // keep if no
	}
	return keep // the list "tokens" but without tags
}

// changeCase applies (up), (low), (cap) to the word before the tag, and numbered tags like (up,2) to the n words before it.

func changeCase(tokens []string) []string {
	for index, word := range tokens { // go through every token

		// up
		if word == "(up)" && index > 0 { // found (tag), and there IS a word before it
			tokens[index-1] = strings.ToUpper(tokens[index-1]) // go -> GO
		}

		// low
		if word == "(low)" && index > 0 {
			tokens[index-1] = strings.ToLower(tokens[index-1]) // SHOUT -> shout
		}

		// cap
		if word == "(cap)" && index > 0 {
			previous_word := tokens[index-1]
			tokens[index-1] = strings.ToUpper(previous_word[:1]) + strings.ToLower(previous_word[1:]) // bRIDGE -> Bridge
		}

		// numbered tags like "(up,2)"
		if strings.HasPrefix(word, "(up,") || strings.HasPrefix(word, "(low,") || strings.HasPrefix(word, "(cap,") {
			comma := strings.Index(word, ",") // position of the comma
			mode := word[1:comma] // "(up,2)" to "up"
			n, err := strconv.Atoi(word[comma+1 : len(word)-1]) // "(up,2)" to 2
			if err == nil {
				for j := 1; j <= n && index-j >= 0; j++ { // the n words before; stop at the start
					tokens[index-j] = applyCase(tokens[index-j], mode)
				}
			}
		}
	}
	return tokens
}


// to identify which case to apply in numbered tags
func applyCase(word string, mode string) string {
	switch mode {
	case "up": // if mode == "up"
		return strings.ToUpper(word)
	case "low": // else if mode == "low"
		return strings.ToLower(word)
	case "cap": // else if mode == "cap"
		return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}
	return word // unknown mode: leave it unchanged
}



// function attaches .,!?:; to the previous word.
func fixPunctuation(tokens []string) []string {
	var out []string // result list

	for _, word := range tokens {
		rest := strings.TrimLeft(word, ".,!?:;")   // word without punctuation
		punctuation := word[:len(word)-len(rest)]   // the punctuation rmeoved

		if punctuation != "" && len(out) > 0 { // there is punctuation AND a previous word
			out[len(out)-1] += punctuation // attach it to the end of the previous word
			word = rest              // keep only what is  left
		}

		if word != "" { // if something left, keep it
			out = append(out, word)
		}
	}
	return out
}

// function attaches ' marks to the words inside them. so for ex "' awesome '" becomes only "'awesome'"
func fixQuotes(tokens []string) []string {
	var newTokens []string
	insideQuote := false        // are we between two ' marks right nowww??
	addQuoteToNextWord := false // should the next word get a ' in front??

	for _, word := range tokens {
		if word == "'" {
			if !insideQuote { // eg false
				addQuoteToNextWord = true
			} else if len(newTokens) > 0 { // closing quote: glue to previous word
				newTokens[len(newTokens)-1] += "'"
			}
			insideQuote = !insideQuote // flip: opening <-> closing
			continue                   // the lone ' is never kept as its own token
		}
		if addQuoteToNextWord { // first word after an opening quote
			word = "'" + word
			addQuoteToNextWord = false
		}
		newTokens = append(newTokens, word)
	}
	return newTokens
}


// function turns "a" into "an" when the next word starts with a vowel or h
func fixAns(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ { // stop one step early, becaus the last word has no next word
		next := strings.ToLower(tokens[i+1])  // lowercase
		if strings.ContainsAny(next[:1], "aeiouh") { // first letter is a vowel or h????
			if tokens[i] == "a" {
				tokens[i] = "an"
			}
			if tokens[i] == "A" {
				tokens[i] = "An" 
			}
		}
	}
	return tokens
}