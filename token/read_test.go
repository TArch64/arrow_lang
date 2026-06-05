package token

import (
	"slices"
	"strings"
	"testing"

	"arrow_lang/testutil"

	"github.com/google/go-cmp/cmp"
)

func TestRead(t *testing.T) {
	type testCase struct {
		name     string
		text     string
		expected []Token
	}

	var testCases = []testCase{
		{
			name: "define variable with int",
			text: "def a = 1",

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
			},
		},
		{
			name: "define variable with negative int",
			text: "def a = -1",

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(-1),
			},
		},
		{
			name: "define variable with float",
			text: "def a = 1.123",

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralFloat(1.123),
			},
		},
		{
			name: "define variable with negative float",
			text: "def a = -1.123",

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralFloat(-1.123),
			},
		},
		{
			name: "with newline at eof",
			text: "def a = 1\n\n",

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
			},
		},
		{
			name: "free variable",

			text: `
				def a = 1
				free a`,

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewKeywordFree(),
				NewIdentifier("a"),
			},
		},
		{
			name: "define variable with sum ints",
			text: `def a = 1 + 2`,

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewOperatorPlus(),
				NewLiteralInt(2),
			},
		},
		{
			name: "assign variable to new variable",

			text: `
				def a = 1
				def b = a`,

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewKeywordDefine(),
				NewIdentifier("b"),
				NewOperatorAssign(),
				NewIdentifier("a"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			text := testutil.Dedent(tc.text)
			result := slices.Collect(Read(strings.NewReader(text)))
			if diff := cmp.Diff(result, tc.expected); diff != "" {
				t.Error(diff)
			}
		})
	}
}
