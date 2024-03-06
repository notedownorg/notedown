package leaf

import (
	"testing"

	"github.com/a-h/parse"
)

func TestThematicBreak(t *testing.T) {
	tests := []struct {
		input    string
		expected ThematicBreak
		notFound bool
	}{
		// These test cases follow the examples from the spec, in order.
		{
			input:    "---",
			expected: ThematicBreak{Character: ThematicBreakCharacterDash},
		},
		{
			input:    "___",
			expected: ThematicBreak{Character: ThematicBreakCharacterUnderscore},
		},
		{
			input:    "***",
			expected: ThematicBreak{Character: ThematicBreakCharacterAsterisk},
		},
		{
			input:    "+++",
			notFound: true,
		},
		{
			input:    "===",
			notFound: true,
		},
		{
			input:    "--",
			notFound: true,
		},
		{
			input:    "**",
			notFound: true,
		},
		{
			input:    "__",
			notFound: true,
		},
		{
			input:    " ***",
			expected: ThematicBreak{Character: ThematicBreakCharacterAsterisk},
		},
		{
			input:    "  ***",
			expected: ThematicBreak{Character: ThematicBreakCharacterAsterisk},
		},
		{
			input:    "   ***",
			expected: ThematicBreak{Character: ThematicBreakCharacterAsterisk},
		},
		{
			input:    "    ***",
			notFound: true,
		},
		// Skip the foo one as it's more general than just thematic breaks.
		{
			input:    "_____________________________________",
			expected: ThematicBreak{Character: ThematicBreakCharacterUnderscore},
		},
		{
			input:    " - - -",
			expected: ThematicBreak{Character: ThematicBreakCharacterDash},
		},
		{
			input:    " **  * ** * ** * **",
			expected: ThematicBreak{Character: ThematicBreakCharacterAsterisk},
		},
		{
			input:    "-     -      -      -",
			expected: ThematicBreak{Character: ThematicBreakCharacterDash},
		},
		{
			input:    "- - - -    ",
			expected: ThematicBreak{Character: ThematicBreakCharacterDash},
		},
        {
            input:    "_ _ _ _ a",
            notFound: true,
        },
        {
            input:    "a------",
            notFound: true,
        },
        {
            input:    "---a---",
            notFound: true,
        },
        {
            input:    "*-*",
            notFound: true,
        },
        // Skip blank lines before and after as this is more general than just thematic breaks.

	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, _, _ := thematicBreakParser.Parse(in)
			if result != test.expected {
				t.Fatalf("expected %#v, but got %#v", test.expected, result)
			}
			if test.notFound {
				if result != (ThematicBreak{}) {
					t.Fatalf("expected not found, but got %#v", result)
				}
				// Ensure we haven't consumed any input.
				if in.Index() != 0 {
					t.Fatalf("expected index to be 0, but got %d", in.Index())
				}
			}
		})
	}

}
