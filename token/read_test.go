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
		expected func() []Token
	}

	var testCases = []testCase{
		{
			name: "basic/define_int",
			text: "def a = 1",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(1),
				}
			},
		},
		{
			name: "basic/define_negative_int",
			text: "def a = -1",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(-1),
				}
			},
		},
		{
			name: "basic/define_float",
			text: "def a = 1.123",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralFloat(1.123),
				}
			},
		},
		{
			name: "basic/define_negative_float",
			text: "def a = -1.123",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralFloat(-1.123),
				}
			},
		},
		{
			name: "basic/free_variable",
			text: `
				def a = 1
				free a`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewKeywordFree(),
					NewIdentifier("a"),
				}
			},
		},
		{
			name:     "edge_cases/empty_input",
			text:     "",
			expected: func() []Token { return nil },
		},
		{
			name:     "edge_cases/whitespace_only",
			text:     "   \t\n  ",
			expected: func() []Token { return nil },
		},
		{
			name: "edge_cases/single_identifier",
			text: "variable",
			expected: func() []Token {
				return []Token{
					NewIdentifier("variable"),
				}
			},
		},
		{
			name: "numbers/zero_integer",
			text: "def a = 0",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(0),
				}
			},
		},
		{
			name: "numbers/large_integer",
			text: "def num = 999999999",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("num"),
					NewOperatorAssign(),
					NewLiteralInt(999999999),
				}
			},
		},
		{
			name: "numbers/float_zero_fractional",
			text: "def val = 1.0",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("val"),
					NewOperatorAssign(),
					NewLiteralFloat(1.0),
				}
			},
		},
		{
			name: "numbers/float_leading_zero",
			text: "def small = 0.123",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("small"),
					NewOperatorAssign(),
					NewLiteralFloat(0.123),
				}
			},
		},
		{
			name: "numbers/high_precision_float",
			text: "def pi = 3.14159265359",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("pi"),
					NewOperatorAssign(),
					NewLiteralFloat(3.14159265359),
				}
			},
		},
		{
			name: "identifiers/with_numbers",
			text: "def var123 = 5",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("var123"),
					NewOperatorAssign(),
					NewLiteralInt(5),
				}
			},
		},
		{
			name: "identifiers/with_underscores",
			text: "def my_variable = 10",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("my_variable"),
					NewOperatorAssign(),
					NewLiteralInt(10),
				}
			},
		},
		{
			name: "identifiers/similar_to_keywords",
			text: "def define = 1 def definition = 2",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("define"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewKeywordDefine(),
					NewIdentifier("definition"),
					NewOperatorAssign(),
					NewLiteralInt(2),
				}
			},
		},
		{
			name: "variables/assign_to_new_variable",
			text: `
				def a = 1
				def b = a`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewKeywordDefine(),
					NewIdentifier("b"),
					NewOperatorAssign(),
					NewIdentifier("a"),
				}
			},
		},
		{
			name: "expressions/sum_integers",
			text: `def a = 1 + 2`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewOperatorPlus(),
					NewLiteralInt(2),
				}
			},
		},
		{
			name: "expressions/sum_literal_and_variable",
			text: `
				def a = 1
				def b = a + 2`,
			expected: func() []Token {
				return []Token{
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
				}
			},
		},
		{
			name: "expressions/chained_additions",
			text: "def result = 1 + 2 + 3",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("result"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewOperatorPlus(),
					NewLiteralInt(2),
					NewOperatorPlus(),
					NewLiteralInt(3),
				}
			},
		},
		{
			name: "expressions/mixed_int_float_addition",
			text: "def result = 1 + 2.5",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("result"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewOperatorPlus(),
					NewLiteralFloat(2.5),
				}
			},
		},
		{
			name: "whitespace/newline_at_eof",
			text: "def a = 1\n\n",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(1),
				}
			},
		},
		{
			name: "whitespace/mixed_tabs_newlines",
			text: "def\ta\t=\n1\n+\n2",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewOperatorPlus(),
					NewLiteralInt(2),
				}
			},
		},
		{
			name: "complex/multiple_definitions_and_operations",
			text: "def x = 1 def y = 2.5 def z = x + y",
			expected: func() []Token {
				return []Token{
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
				}
			},
		},
		{
			name: "expressions/subtract_integers",
			text: `def a = 5 - 2`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(5),
					NewOperatorMinus(),
					NewLiteralInt(2),
				}
			},
		},
		{
			name: "expressions/subtract_literal_and_variable",
			text: `
				def a = 10
				def b = a - 3`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(10),
					NewKeywordDefine(),
					NewIdentifier("b"),
					NewOperatorAssign(),
					NewIdentifier("a"),
					NewOperatorMinus(),
					NewLiteralInt(3),
				}
			},
		},
		{
			name: "expressions/subtract_variable_and_literal",
			text: `
				def a = 5
				def b = 10 - a`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("a"),
					NewOperatorAssign(),
					NewLiteralInt(5),
					NewKeywordDefine(),
					NewIdentifier("b"),
					NewOperatorAssign(),
					NewLiteralInt(10),
					NewOperatorMinus(),
					NewIdentifier("a"),
				}
			},
		},
		{
			name: "expressions/chained_subtractions",
			text: "def result = 10 - 3 - 2",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("result"),
					NewOperatorAssign(),
					NewLiteralInt(10),
					NewOperatorMinus(),
					NewLiteralInt(3),
					NewOperatorMinus(),
					NewLiteralInt(2),
				}
			},
		},
		{
			name: "expressions/mixed_int_float_subtraction",
			text: "def result = 5 - 2.5",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("result"),
					NewOperatorAssign(),
					NewLiteralInt(5),
					NewOperatorMinus(),
					NewLiteralFloat(2.5),
				}
			},
		},
		{
			name: "expressions/float_subtraction",
			text: "def result = 3.14 - 1.5",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("result"),
					NewOperatorAssign(),
					NewLiteralFloat(3.14),
					NewOperatorMinus(),
					NewLiteralFloat(1.5),
				}
			},
		},
		{
			name: "expressions/subtract_from_zero",
			text: "def negative = 0 - 5",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("negative"),
					NewOperatorAssign(),
					NewLiteralInt(0),
					NewOperatorMinus(),
					NewLiteralInt(5),
				}
			},
		},
		{
			name: "expressions/mixed_addition_subtraction",
			text: "def result = 1 + 2 - 3",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("result"),
					NewOperatorAssign(),
					NewLiteralInt(1),
					NewOperatorPlus(),
					NewLiteralInt(2),
					NewOperatorMinus(),
					NewLiteralInt(3),
				}
			},
		},
		{
			name: "expressions/mixed_subtraction_addition",
			text: "def result = 10 - 3 + 2",
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("result"),
					NewOperatorAssign(),
					NewLiteralInt(10),
					NewOperatorMinus(),
					NewLiteralInt(3),
					NewOperatorPlus(),
					NewLiteralInt(2),
				}
			},
		},
		{
			name: "expressions/complex_arithmetic_with_variables",
			text: `
				def x = 5
				def y = 3
				def z = x + y - 2`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("x"),
					NewOperatorAssign(),
					NewLiteralInt(5),
					NewKeywordDefine(),
					NewIdentifier("y"),
					NewOperatorAssign(),
					NewLiteralInt(3),
					NewKeywordDefine(),
					NewIdentifier("z"),
					NewOperatorAssign(),
					NewIdentifier("x"),
					NewOperatorPlus(),
					NewIdentifier("y"),
					NewOperatorMinus(),
					NewLiteralInt(2),
				}
			},
		},
		{
			name: "functions/basic_getter",
			text: `def fn { ret 1 }`,

			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("fn"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(1),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/return_zero",
			text: `def zero { ret 0 }`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("zero"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(0),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/return_negative_int",
			text: `def fn { ret -1 }`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("fn"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(-1),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/return_float",
			text: `def pi { ret 3.14 }`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("pi"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralFloat(3.14),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/return_sum_expression",
			text: `def fn { ret 1 + 2 }`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("fn"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(1),
					NewOperatorPlus(),
					NewLiteralInt(2),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/return_subtraction_expression",
			text: `def fn { ret 10 - 3 }`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("fn"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(10),
					NewOperatorMinus(),
					NewLiteralInt(3),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/return_mixed_expression",
			text: `def fn { ret 1 + 2 - 3 }`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("fn"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(1),
					NewOperatorPlus(),
					NewLiteralInt(2),
					NewOperatorMinus(),
					NewLiteralInt(3),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/multiple_getters",
			text: `
				def one { ret 1 }
				def two { ret 2 }`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("one"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(1),
					NewCurlyBracketClose(),
					NewKeywordDefine(),
					NewIdentifier("two"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(2),
					NewCurlyBracketClose(),
				}
			},
		},
		{
			name: "functions/getter_multiline_body",
			text: `
				def fn {
					ret 42
				}`,
			expected: func() []Token {
				return []Token{
					NewKeywordDefine(),
					NewIdentifier("fn"),
					NewCurlyBracketOpen(),
					NewKeywordReturn(),
					NewLiteralInt(42),
					NewCurlyBracketClose(),
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			text := testutil.Dedent(tc.text)
			result := slices.Collect(Read(strings.NewReader(text)))
			if diff := cmp.Diff(result, tc.expected()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
