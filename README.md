# Go Text Formatter

> A command-line text auto-corrector in Go that reads a text file, applies formatting rules (hex/bin to decimal, case changes, punctuation spacing, quotes, a→an), and writes the result to a new file.

## Example

```
Input:  it (cap) was a amazing day : 1E (hex) files and 10 (bin) bugs , but I WAS VERY FINE (low, 3) ,and she said ' this is so exciting (up, 2) ' and I thought ... what a honest wow (up) !?
Output: It was an amazing day: 30 files and 2 bugs, but I was very fine, and she said 'this is SO EXCITING' and I thought... what an honest WOW!?
```

## Rules

| Rule | Example |
|---|---|
| `(hex)` / `(bin)` → decimal | `1E (hex)` → `30`, `10 (bin)` → `2` |
| `(up)`, `(low)`, `(cap)` | `go (up)` → `GO` |
| Numbered: `(up, n)` etc. | `so exciting (up, 2)` → `SO EXCITING` |
| Punctuation `.,!?:;` | `there ,and` → `there, and`, `BAMM !!` → `BAMM!!` |
| Quotes `' '` | `' awesome '` → `'awesome'` |
| a → an before a vowel or h | `a amazing` → `an amazing` |

## Usage

```bash
git clone https://github.com/zinebnadak/go-text-formatter.git
cd go-text-formatter
go run . sample.txt results.txt
cat result.txt
```

## Tests

```bash
go test -v
```

## Tech

- Language: Go
- Dependencies: [Go's standard library](https://pkg.go.dev/std) only

## What I learned personally

- Reading and writing files with `os`
- String and number manipulation (`strings`, `strconv`)
- Why the order of processing steps matters
- [Unit testing](https://go.dev/doc/tutorial/add-a-test) in Go

---

Built by Zineb at [grit:lab](https://gritlab.ax), Åland (01-edu peer-learning program).