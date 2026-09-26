# Input file: the program reads it and never changes it.
# Output file: the program writes the modified text into it.

# Command looks like go run . sample.txt result.txt


18 hours 

**Day 1: Foundation**
- 1h: finish understanding (rule order, edge cases)
- 1h: plain-English plan for the whole pipeline
- 1h: `main.go` with args check, read the file, write the file
- 3h: split text into words, then `(hex)` and `(bin)`

**Day 2: The rules**
- 3h: `(up)`, `(low)`, `(cap)`, including the numbered `(up, 2)` version
- 1h: a→an
- 2h: punctuation spacing, including `...` and `!?`

**Day 3: Finish and defend**
- 2h: `'` quotes (the trickiest rule, so budget for it)
- 2h: a test file covering every example in the subject
- 1h: `gofmt`, cleanup, README
- 1h: mock audit, explaining every function out loud to yourself or a friend

**BUFFER DAY**

Tests & mock audit.

### Instructions
- (hex)
- (bin)

For (low), (up), (cap) if a number appears next to it, like so: (low, <number>) it turns the previously specified number of words in lowercase, uppercase or capitalized accordingly:
- (up)
- (low)
- (cap)

- Every instance of the punctuations ., ,, !, ?, : and ; should be close to the previous word and with space apart from the next one (if there are groups of punctuation like: ... or !? the program should format the text like "Thinking ..." to "Thinking..."
- The punctuation mark ' will always be found with another instance of it and they should be placed to the right and left of the word in the middle of them, without any spaces. If there are more than one word between the two ' ' marks, the program should place the marks next to the corresponding words
- Every instance of a should be turned into an if the next word begins with a vowel (a, e, i, o, u)

### Notes
- Input: first filename with text that has tags, and second filename that the modified text should to into
- Output: new textfile with modified text

- Order matters : 

All tags must be applied and removed before a→an, quotes and puctuation:
ex. In A (cap) owl, if a→an runs first, the next word is (cap), not owl. ( isn't a vowel, so A stays A, and you end up with A owl. Wrong.

ex. Punctuation attaches ! to the previous token, so the tag becomes "(up)!". Your tag code looks for exactly "(up)", so it misses it and go never gets uppercased.

- Remove all of the tags for the output text. The number is in the second token; strip ) from it
Problem: a numbered tag does not live in one single token Because we have spaces here (up, 2) and it then becomes two tokens: "(up," and "2)" need to strip out )

- If the tag asks for 5 words, but only 2 come before it the program should uppercase only the words that exist (clamp 5 down to 2)
- (up) as the very first word of the file, with no word before it the program should do nothing and remove the tag.

### Plan in plain english 
(reading a file to writing a file)

1. take the command line arguments
2. read the file
3. split the text into tokens (words with indexes)
4. apply all tags
5. remove the tags
6. punctuation
7. quotes
8. a to an
9. join tokens into a single string
10. write (also creates) the file