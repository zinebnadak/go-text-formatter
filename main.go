package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// the program needs exactly 2 arguments: input file and output file
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . input.txt output.txt")
		return
	}

	// read the input file
	fileContents, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// apply every rule to the text
	result := processText(string(fileContents))

	// write the output file (creates it if it doesn't exist)
	err = os.WriteFile(os.Args[2], []byte(result), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// The test calls processText directly with a string and compares the result to what to expect. It skips main entirely: no files or os.Args
func processText(text string) string {
	text = strings.ReplaceAll(text, "(up, ", "(up,") // "(up, 2)" -> "(up,2)" so it stays one token
	text = strings.ReplaceAll(text, "(low, ", "(low,")
	text = strings.ReplaceAll(text, "(cap, ", "(cap,")

	tokens := strings.Fields(text) // split into words, drops extra spaces
	tokens = convertBases(tokens) // hex & bin
	tokens = changeCase(tokens)    // up, low, cap (+ numbered)
	tokens = removeTags(tokens)    // tags out, AFTER applying them
	tokens = fixPunctuation(tokens)
	tokens = fixQuotes(tokens)
	tokens = fixArticles(tokens)
	return strings.Join(tokens, " ")
}

// convertBases replaces the word before (hex) or (bin) with its decimal value.
// Hex is base 16, bin is base 2.
func convertBases(tokens []string) []string {
	for index, word := range tokens {

		// hex
		if word == "(hex)" && index > 0 {
			number, err := strconv.ParseInt(tokens[index-1], 16, 64) // read as base 16, fits in int64
			if err == nil {                                          // only if it really was hex
				tokens[index-1] = strconv.FormatInt(number, 10) // number back to text, in base 10
			}
		}

		// bin
		if word == "(bin)" && index > 0 {
			number, err := strconv.ParseInt(tokens[index-1], 2, 64) // read as base 2
			if err == nil {
				tokens[index-1] = strconv.FormatInt(number, 10)
			}
		}
	}
	return tokens
}

// changeCase applies (up), (low), (cap) to the word before the tag,
// and numbered tags like (up,2) to the n words before it.
func changeCase(tokens []string) []string {
	for index, word := range tokens {

		// up
		if word == "(up)" && index > 0 { // found the tag, and there IS a word before it
			tokens[index-1] = strings.ToUpper(tokens[index-1]) // go -> GO
		}

		// low
		if word == "(low)" && index > 0 {
			tokens[index-1] = strings.ToLower(tokens[index-1]) // SHOUT -> shout
		}

		// cap
		if word == "(cap)" && index > 0 {
			previousWord := tokens[index-1]
			tokens[index-1] = strings.ToUpper(previousWord[:1]) + strings.ToLower(previousWord[1:]) // bRIDGE -> Bridge
		}

		// numbered tags like "(up,2)"
		if strings.HasPrefix(word, "(up,") || strings.HasPrefix(word, "(low,") || strings.HasPrefix(word, "(cap,") {
			comma := strings.Index(word, ",")                   // position of the comma
			mode := word[1:comma]                               // "(up,2)" -> "up"
			n, err := strconv.Atoi(word[comma+1 : len(word)-1]) // "(up,2)" -> 2
			if err == nil {
				for j := 1; j <= n && index-j >= 0; j++ { // the n words before; stop at the start
					tokens[index-j] = applyCase(tokens[index-j], mode)
				}
			}
		}
	}
	return tokens
}

// applyCase changes one word according to mode: "up", "low" or "cap".
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

// removeTags returns the tokens without any tag words.
func removeTags(tokens []string) []string {
	var keep []string // empty slice to fill
	for _, word := range tokens {
		if word == "(hex)" || word == "(bin)" || word == "(up)" || word == "(low)" || word == "(cap)" ||
			strings.HasPrefix(word, "(up,") || strings.HasPrefix(word, "(low,") || strings.HasPrefix(word, "(cap,") { // is it a tag?
			continue // yes: skip it
		}
		keep = append(keep, word) // no: keep it
	}
	return keep
}

// fixPunctuation attaches .,!?:; to the previous word.
// "there ,and" -> "there, and"   "BAMM !!" -> "BAMM!!"
func fixPunctuation(tokens []string) []string {
	var out []string

	for _, word := range tokens {
		rest := strings.TrimLeft(word, ".,!?:;")   // word without leading punctuation
		punctuation := word[:len(word)-len(rest)] // the punctuation that was removed

		if punctuation != "" && len(out) > 0 { // there is punctuation AND a previous word
			out[len(out)-1] += punctuation // attach it to the end of the previous word
			word = rest                    // keep only what is left
		}

		if word != "" { // something left? keep it
			out = append(out, word)
		}
	}
	return out
}

// fixQuotes attaches ' marks to the words inside them.
// "' awesome '" -> "'awesome'"
func fixQuotes(tokens []string) []string {
	var newTokens []string
	insideQuote := false        // are we between two ' marks right now?
	addQuoteToNextWord := false // should the next word get a ' in front?

	for _, word := range tokens {
		if word == "'" {
			if !insideQuote { // opening quote
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

// fixArticles turns "a" into "an" when the next word starts with a vowel or h.
func fixArticles(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ { // stop one early: the last word has no next word
		next := strings.ToLower(tokens[i+1])         // lowercase, so "Amazing" counts too
		if strings.ContainsAny(next[:1], "aeiouh") { // first letter is a vowel or h?
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