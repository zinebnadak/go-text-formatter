// test file (file's name with _test.go for go test command)

package main

import "testing"

// TestProcessText checks every audit example and our own edge cases.
func TestProcessText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// audit examples
		{"If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?",
			"If I make you breakfast in bed just say thank you instead of: How did you get in MY HOUSE?"},
		{"I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure",
			"I have to pack 5 outfits. Packed 26 just to be sure"},
		{"Don not be sad ,because sad backwards is das . And das not good",
			"Don not be sad, because sad backwards is das. And das not good"},
		{"harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '",
			"Harold Wilson: 'I am an optimist, but an optimist who carries a raincoat.'"},
		{"it (cap) was the best of times, it was the worst of times (up) , it was the age of foolishness (cap, 6) , IT WAS THE (low, 3) winter of despair.",
			"It was the best of times, it was the worst of TIMES, It Was The Age Of Foolishness, it was the winter of despair."},

		// our edge cases from the notes
		{"(up) hello", "hello"},                  // tag first: do nothing, remove tag
		{"hello world (up, 5)", "HELLO WORLD"},   // count too big: clamp
	}

	for _, tc := range tests { // run every case
		got := processText(tc.input)
		if got != tc.want {
			t.Errorf("\ninput: %q\n got:  %q\n want: %q", tc.input, got, tc.want)
		}
	}
}