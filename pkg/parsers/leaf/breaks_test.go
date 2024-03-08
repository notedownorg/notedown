package leaf_test

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/parsers/leaf"
)

func TestThematicBreak(t *testing.T) {
	tests := []struct {
		input    string
		expected leaf.ThematicBreak
		notFound bool
	}{
		// These test cases follow the examples from the spec, in order.
		{
			input:    "---\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterDash},
		},
		{
			input:    "___\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterUnderscore},
		},
		{
			input:    "***\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterAsterisk},
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
			input:    " ***\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterAsterisk},
		},
		{
			input:    "  ***\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterAsterisk},
		},
		{
			input:    "   ***\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterAsterisk},
		},
		{
			input:    "    ***\n",
			notFound: true,
		},
		// Skip the foo one as it's more general than just thematic breaks.
		{
			input:    "_____________________________________\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterUnderscore},
		},
		{
			input:    " - - -\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterDash},
		},
		{
			input:    " **  * ** * ** * **\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterAsterisk},
		},
		{
			input:    "-     -      -      -\n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterDash},
		},
		{
			input:    "- - - -    \n",
			expected: leaf.ThematicBreak{Character: leaf.ThematicBreakCharacterDash},
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
			result, found, _ := leaf.ThematicBreakParser.Parse(in)
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

			if result != test.expected {
				t.Fatalf("expected %#v, but got %#v", test.expected, result)
			}
		})
	}

}
