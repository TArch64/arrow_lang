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
			name: "basic/define_int",
			text: "def a = 1",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
			},
		},
		{
			name: "basic/define_negative_int",
			text: "def a = -1",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(-1),
			},
		},
		{
			name: "basic/define_float",
			text: "def a = 1.123",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralFloat(1.123),
			},
		},
		{
			name: "basic/define_negative_float",
			text: "def a = -1.123",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralFloat(-1.123),
			},
		},
		{
			name: "basic/free_variable",
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
			name:     "edge_cases/empty_input",
			text:     "",
			expected: nil,
		},
		{
			name:     "edge_cases/whitespace_only",
			text:     "   \t\n  ",
			expected: nil,
		},
		{
			name: "edge_cases/single_identifier",
			text: "variable",
			expected: []Token{
				NewIdentifier("variable"),
			},
		},
		{
			name: "numbers/zero_integer",
			text: "def a = 0",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(0),
			},
		},
		{
			name: "numbers/large_integer",
			text: "def num = 999999999",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("num"),
				NewOperatorAssign(),
				NewLiteralInt(999999999),
			},
		},
		{
			name: "numbers/float_zero_fractional",
			text: "def val = 1.0",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("val"),
				NewOperatorAssign(),
				NewLiteralFloat(1.0),
			},
		},
		{
			name: "numbers/float_leading_zero",
			text: "def small = 0.123",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("small"),
				NewOperatorAssign(),
				NewLiteralFloat(0.123),
			},
		},
		{
			name: "numbers/high_precision_float",
			text: "def pi = 3.14159265359",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("pi"),
				NewOperatorAssign(),
				NewLiteralFloat(3.14159265359),
			},
		},
		{
			name: "identifiers/with_numbers",
			text: "def var123 = 5",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("var123"),
				NewOperatorAssign(),
				NewLiteralInt(5),
			},
		},
		{
			name: "identifiers/with_underscores",
			text: "def my_variable = 10",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("my_variable"),
				NewOperatorAssign(),
				NewLiteralInt(10),
			},
		},
		{
			name: "identifiers/similar_to_keywords",
			text: "def define = 1 def definition = 2",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("define"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewKeywordDefine(),
				NewIdentifier("definition"),
				NewOperatorAssign(),
				NewLiteralInt(2),
			},
		},
		{
			name: "variables/assign_to_new_variable",
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
		{
			name: "expressions/sum_integers",
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
			name: "expressions/sum_literal_and_variable",
			text: `
				def a = 1
				def b = a + 2`,
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewKeywordDefine(),
				NewIdentifier("b"),
				NewOperatorAssign(),
				NewIdentifier("a"),
				NewOperatorPlus(),
				NewLiteralInt(2),
			},
		},
		{
			name: "expressions/chained_additions",
			text: "def result = 1 + 2 + 3",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("result"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewOperatorPlus(),
				NewLiteralInt(2),
				NewOperatorPlus(),
				NewLiteralInt(3),
			},
		},
		{
			name: "expressions/mixed_int_float_addition",
			text: "def result = 1 + 2.5",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("result"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewOperatorPlus(),
				NewLiteralFloat(2.5),
			},
		},
		{
			name: "whitespace/newline_at_eof",
			text: "def a = 1\n\n",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("a"),
				NewOperatorAssign(),
				NewLiteralInt(1),
			},
		},
		{
			name: "whitespace/mixed_tabs_newlines",
			text: "def\ta\t=\n1\n+\n2",
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
			name: "complex/multiple_definitions_and_operations",
			text: "def x = 1 def y = 2.5 def z = x + y",
			expected: []Token{
				NewKeywordDefine(),
				NewIdentifier("x"),
				NewOperatorAssign(),
				NewLiteralInt(1),
				NewKeywordDefine(),
				NewIdentifier("y"),
				NewOperatorAssign(),
				NewLiteralFloat(2.5),
				NewKeywordDefine(),
				NewIdentifier("z"),
				NewOperatorAssign(),
				NewIdentifier("x"),
				NewOperatorPlus(),
				NewIdentifier("y"),
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
