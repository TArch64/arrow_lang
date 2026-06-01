package token

import (
	"slices"
	"strings"
	"testing"

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
			name: "define variable with literal int",
			text: "def a = 1",

			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := slices.Collect(Read(strings.NewReader(tc.text)))
			if diff := cmp.Diff(result, tc.expected); diff != "" {
				t.Error(diff)
			}
		})
	}
}
