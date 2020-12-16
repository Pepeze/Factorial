package main

import "testing"

// Run tests by typing "go test" in the terminal.
// Show a more thorough output by typing "go test -v".

func TestMin(t *testing.T) {
	if min(1, 2) != 1 {
		t.Error("Expected the minimum of 1 and 2 to be 1.")
	}
}

func TestTableMin(t *testing.T) {
	var tests = []struct {
		input1   int
		input2   int
		expected int
	}{
		{1, 2, 1},
		{4, 3, 3},
		{10, 11, 10},
		{12, 13, 12},
	}

	for _, test := range tests {
		if output := min(test.input1, test.input2); output != test.expected {
			t.Error("Test failed:", test.input1, test.input2, "inputted,", test.expected, "expected,", output, "outputted.")
		}
	}
}
