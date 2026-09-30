// test file (file's name with _test.go for go test command)

package main

import "testing"

// TestProcessText checks the subject's examples and our own edge cases.
func TestProcessText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// examples from the subject
		{"it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.",
			"It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair."},
		{"Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			"Simply add 66 and 2 and you will see the result is 68."},
		{"There is no greater agony than bearing a untold story inside you.",
			"There is no greater agony than bearing an untold story inside you."},
		{"Punctuation tests are ... kinda boring ,what do you think ?",
			"Punctuation tests are... kinda boring, what do you think?"},

		// edge cases from the notes
		{"(up) hello", "hello"},                // tag first: do nothing, remove tag
		{"hello world (up, 5)", "HELLO WORLD"}, // count too big: clamp
	}

	for _, tc := range tests { // run every case
		got := processText(tc.input)
		if got != tc.want {
			t.Errorf("\ninput: %q\n got:  %q\n want: %q", tc.input, got, tc.want)
		}
	}
}