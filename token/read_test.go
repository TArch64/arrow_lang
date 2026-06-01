package token

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	type testCase struct {
		name     string
		text     string
		expected string
	}

	var testCases = []testCase{
		{
			name:     "define variable with literal int value",
			text:     "def a = 1",
			expected: "Keyword(def) Identifier(a) Operator(=) Int(1)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result strings.Builder
			for token := range Read(strings.NewReader(tc.text)) {
				if result.Len() != 0 {
					result.WriteRune(' ')
				}
				result.WriteString(token.String())
			}
			if tc.expected != result.String() {
				t.Errorf("\nexpected: %s\ngot     : %s", tc.expected, result.String())
			}
		})
	}
}
