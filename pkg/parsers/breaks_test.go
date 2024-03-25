package parsers_test

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/parsers"
)

func TestThematicBreak(t *testing.T) {
	tests := []struct {
		input    string
		notFound bool
	}{
		// These test cases follow the examples from the spec, in order.
		{
			input: "---\n",
		},
		{
			input: "___\n",
		},
		{
			input: "***\n",
		},
		{
			input:    "+++\n",
			notFound: true,
		},
		{
			input:    "===\n",
			notFound: true,
		},
		{
			input:    "--\n",
			notFound: true,
		},
		{
			input:    "**\n",
			notFound: true,
		},
		{
			input:    "__\n",
			notFound: true,
		},
		{
			input: " ***\n",
		},
		{
			input: "  ***\n",
		},
		{
			input: "   ***\n",
		},
		{
			input:    "    ***\n",
			notFound: true,
		},
		// Skip the foo one as it's more general than just thematic breaks.
		{
			input: "_____________________________________\n",
		},
		{
			input: " - - -\n",
		},
		{
			input: " **  * ** * ** * **\n",
		},
		{
			input: "-     -      -      -\n",
		},
		{
			input: "- - - -    \n",
		},
		{
			input:    "_ _ _ _ a\n",
			notFound: true,
		},
		{
			input:    "a------\n",
			notFound: true,
		},
		{
			input:    "---a---\n",
			notFound: true,
		},
		{
			input:    "*-*\n",
			notFound: true,
		},
		// Skip blank lines before and after as this is more general than just thematic breaks.

	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := parsers.ThematicBreak.Parse(in)
			if test.notFound {
				if found {
					t.Fatal("expected not found")
				}
				// Ensure we haven't consumed any input.
				if in.Index() != 0 {
					t.Fatalf("expected index to be 0, but got %d", in.Index())
				}
				return
			}

			if result != test.input {
				t.Fatalf("expected %#v, but got %#v", test.input, result)
			}
		})
	}

}
