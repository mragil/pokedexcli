package repl_test

import (
	"pokedexcli/internal/repl"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello     world",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := repl.CleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("expected: %v, got: %v", c.expected, actual)
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if expectedWord != word {
				t.Errorf("expected: %v, got: %v", expectedWord, word)
			}
		}
	}
}
