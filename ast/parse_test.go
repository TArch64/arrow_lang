package ast

import (
	"errors"
	"slices"
	"testing"

	"arrow_lang/token"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	type testCase struct {
		name         string
		tokens       func() []token.Token
		expectedErr  error
		expectedNode func() Node
	}

	testCases := []testCase{
		{
			name: "error_cases/define_variable_eof_after_keyword",
			tokens: func() []token.Token {
				return []token.Token{token.NewKeywordDefine()}
			},
			expectedErr: UnexpectedEOFErr,
		},
		{
			name: "error_cases/define_variable_invalid_token_after_keyword",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewOperatorAssign(),
				}
			},
			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "error_cases/define_variable_eof_after_identifier",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
				}
			},
			expectedErr: UnexpectedEOFErr,
		},
		{
			name: "error_cases/define_variable_invalid_token_after_identifier",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorPlus(),
				}
			},
			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "error_cases/expression_incomplete_sum",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
					token.NewOperatorAssign(),
				}
			},
			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "error_cases/undefined_variable_reference",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
					token.NewKeywordDefine(),
					token.NewIdentifier("b"),
					token.NewOperatorAssign(),
					token.NewIdentifier("c"),
				}
			},
			expectedErr: UndefinedVariableErr,
		},
		{
			name: "basic/define_variable_int",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("a",
							NewExpression([]DataNode{NewLiteralInt(1)}),
						),
					),
				})
			},
		},
		{
			name: "basic/define_variable_negative_int",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(-1),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("a",
							NewExpression([]DataNode{NewLiteralInt(-1)}),
						),
					),
				})
			},
		},
		{
			name: "basic/define_variable_float",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(1.123),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("a",
							NewExpression([]DataNode{NewLiteralFloat(1.123)}),
						),
					),
				})
			},
		},
		{
			name: "basic/define_variable_negative_float",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(-1.123),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("a",
							NewExpression([]DataNode{NewLiteralFloat(-1.123)}),
						),
					),
				})
			},
		},
		{
			name: "basic/free_variable",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
					token.NewKeywordFree(),
					token.NewIdentifier("a"),
				}
			},
			expectedNode: func() Node {
				defA := NewVariable("a",
					NewExpression([]DataNode{NewLiteralInt(1)}),
				)
				return NewProgram([]*Statement{
					NewStatement(defA),
					NewStatement(NewFree(defA)),
				})
			},
		},
		{
			name: "expressions/sum_two_integers",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
					token.NewOperatorPlus(),
					token.NewLiteralInt(2),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("a",
							NewExpression([]DataNode{NewLiteralInt(3)}),
						),
					),
				})
			},
		},
		{
			name: "variables/assign_variable_to_variable",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
					token.NewKeywordDefine(),
					token.NewIdentifier("b"),
					token.NewOperatorAssign(),
					token.NewIdentifier("a"),
				}
			},
			expectedNode: func() Node {
				defA := NewVariable("a",
					NewExpression([]DataNode{NewLiteralInt(1)}),
				)
				return NewProgram([]*Statement{
					NewStatement(defA),
					NewStatement(
						NewVariable("b",
							NewExpression([]DataNode{NewVariableReference(defA)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/sum_variable_and_literal",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
					token.NewKeywordDefine(),
					token.NewIdentifier("b"),
					token.NewOperatorAssign(),
					token.NewIdentifier("a"),
					token.NewOperatorPlus(),
					token.NewLiteralInt(2),
				}
			},
			expectedNode: func() Node {
				defA := NewVariable("a", NewExpression([]DataNode{NewLiteralInt(1)}))

				return NewProgram([]*Statement{
					NewStatement(defA),
					NewStatement(
						NewVariable("b",
							NewExpression([]DataNode{
								NewVariableReference(defA),
								NewExpressionPlus(NewLiteralInt(2)),
							}),
						),
					),
				})
			},
		},
		{
			name: "numbers/zero_integer",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("zero"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(0),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("zero",
							NewExpression([]DataNode{NewLiteralInt(0)}),
						),
					),
				})
			},
		},
		{
			name: "numbers/large_integer",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("big"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(999999999),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("big",
							NewExpression([]DataNode{NewLiteralInt(999999999)}),
						),
					),
				})
			},
		},
		{
			name: "numbers/float_zero_fractional",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("val"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(5.0),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("val",
							NewExpression([]DataNode{NewLiteralFloat(5.0)}),
						),
					),
				})
			},
		},
		{
			name: "numbers/float_leading_zero",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("small"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(0.5),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("small",
							NewExpression([]DataNode{NewLiteralFloat(0.5)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/sum_two_floats",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(1.5),
					token.NewOperatorPlus(),
					token.NewLiteralFloat(2.5),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralFloat(4.0)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/sum_mixed_int_float",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("mixed"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(1),
					token.NewOperatorPlus(),
					token.NewLiteralFloat(2.5),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("mixed",
							NewExpression([]DataNode{NewLiteralFloat(3.5)}),
						),
					),
				})
			},
		},
		{
			name: "complex/multiple_statements",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("x"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(10),
					token.NewKeywordDefine(),
					token.NewIdentifier("y"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(20),
					token.NewKeywordDefine(),
					token.NewIdentifier("sum"),
					token.NewOperatorAssign(),
					token.NewIdentifier("x"),
					token.NewOperatorPlus(),
					token.NewIdentifier("y"),
				}
			},
			expectedNode: func() Node {
				defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(10)}))
				defY := NewVariable("y", NewExpression([]DataNode{NewLiteralInt(20)}))
				return NewProgram([]*Statement{
					NewStatement(defX),
					NewStatement(defY),
					NewStatement(
						NewVariable("sum",
							NewExpression([]DataNode{
								NewVariableReference(defX),
								NewExpressionPlus(NewVariableReference(defY)),
							}),
						),
					),
				})
			},
		},
		{
			name: "complex/define_variable_and_free_sequence",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("temp"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(42),
					token.NewKeywordDefine(),
					token.NewIdentifier("copy"),
					token.NewOperatorAssign(),
					token.NewIdentifier("temp"),
					token.NewKeywordFree(),
					token.NewIdentifier("temp"),
				}
			},
			expectedNode: func() Node {
				defTemp := NewVariable("temp", NewExpression([]DataNode{NewLiteralInt(42)}))
				return NewProgram([]*Statement{
					NewStatement(defTemp),
					NewStatement(
						NewVariable("copy",
							NewExpression([]DataNode{NewVariableReference(defTemp)}),
						),
					),
					NewStatement(NewFree(defTemp)),
				})
			},
		},

		{
			name: "expressions/subtract_two_integers_basic",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(5),
					token.NewOperatorMinus(),
					token.NewLiteralInt(3),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(2)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_two_floats",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(7.5),
					token.NewOperatorMinus(),
					token.NewLiteralFloat(2.25),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralFloat(5.25)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_mixed_int_from_float",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(10.5),
					token.NewOperatorMinus(),
					token.NewLiteralInt(5),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralFloat(5.5)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_mixed_float_from_int",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(8),
					token.NewOperatorMinus(),
					token.NewLiteralFloat(3.5),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralFloat(4.5)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_zero_from_integer",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(42),
					token.NewOperatorMinus(),
					token.NewLiteralInt(0),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(42)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_same_number",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(7),
					token.NewOperatorMinus(),
					token.NewLiteralInt(7),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(0)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_resulting_negative",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(3),
					token.NewOperatorMinus(),
					token.NewLiteralInt(8),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(-5)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_from_negative_number",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(-10),
					token.NewOperatorMinus(),
					token.NewLiteralInt(5),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(-15)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_negative_number",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(10),
					token.NewOperatorMinus(),
					token.NewLiteralInt(-3),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(13)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_variable_from_literal",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("x"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(15),
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(20),
					token.NewOperatorMinus(),
					token.NewIdentifier("x"),
				}
			},
			expectedNode: func() Node {
				defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(15)}))
				return NewProgram([]*Statement{
					NewStatement(defX),
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{
								NewLiteralInt(20),
								NewExpressionMinus(NewVariableReference(defX)),
							}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_literal_from_variable",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("x"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(25),
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewIdentifier("x"),
					token.NewOperatorMinus(),
					token.NewLiteralInt(12),
				}
			},
			expectedNode: func() Node {
				defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(25)}))
				return NewProgram([]*Statement{
					NewStatement(defX),
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{
								NewVariableReference(defX),
								NewExpressionMinus(NewLiteralInt(12)),
							}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_variable_from_variable",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("x"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(30),
					token.NewKeywordDefine(),
					token.NewIdentifier("y"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(18),
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewIdentifier("x"),
					token.NewOperatorMinus(),
					token.NewIdentifier("y"),
				}
			},
			expectedNode: func() Node {
				defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(30)}))
				defY := NewVariable("y", NewExpression([]DataNode{NewLiteralInt(18)}))
				return NewProgram([]*Statement{
					NewStatement(defX),
					NewStatement(defY),
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{
								NewVariableReference(defX),
								NewExpressionMinus(NewVariableReference(defY)),
							}),
						),
					),
				})
			},
		},
		{
			name: "expressions/chained_minus_operations",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(100),
					token.NewOperatorMinus(),
					token.NewLiteralInt(25),
					token.NewOperatorMinus(),
					token.NewLiteralInt(15),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(60)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/mixed_plus_minus_operations",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(10),
					token.NewOperatorPlus(),
					token.NewLiteralInt(5),
					token.NewOperatorMinus(),
					token.NewLiteralInt(3),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralInt(12)}),
						),
					),
				})
			},
		},
		{
			name: "expressions/complex_mixed_operations_with_variables",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("a"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(50),
					token.NewKeywordDefine(),
					token.NewIdentifier("b"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(20),
					token.NewKeywordDefine(),
					token.NewIdentifier("c"),
					token.NewOperatorAssign(),
					token.NewLiteralInt(8),
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewIdentifier("a"),
					token.NewOperatorMinus(),
					token.NewIdentifier("b"),
					token.NewOperatorPlus(),
					token.NewIdentifier("c"),
				}
			},
			expectedNode: func() Node {
				defA := NewVariable("a", NewExpression([]DataNode{NewLiteralInt(50)}))
				defB := NewVariable("b", NewExpression([]DataNode{NewLiteralInt(20)}))
				defC := NewVariable("c", NewExpression([]DataNode{NewLiteralInt(8)}))
				return NewProgram([]*Statement{
					NewStatement(defA),
					NewStatement(defB),
					NewStatement(defC),
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{
								NewVariableReference(defA),
								NewExpressionMinus(NewVariableReference(defB)),
								NewExpressionPlus(NewVariableReference(defC)),
							}),
						),
					),
				})
			},
		},
		{
			name: "expressions/subtract_float_precision",
			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("result"),
					token.NewOperatorAssign(),
					token.NewLiteralFloat(3.14159),
					token.NewOperatorMinus(),
					token.NewLiteralFloat(2.71828),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewVariable("result",
							NewExpression([]DataNode{NewLiteralFloat(0.42330999999999985)}),
						),
					),
				})
			},
		},
		{
			name: "functions/basic_getter",

			tokens: func() []token.Token {
				return []token.Token{
					token.NewKeywordDefine(),
					token.NewIdentifier("fn"),
					token.NewCurlyBracketOpen(),
					token.NewKeywordReturn(),
					token.NewLiteralInt(1),
					token.NewCurlyBracketClose(),
				}
			},
			expectedNode: func() Node {
				return NewProgram([]*Statement{
					NewStatement(
						NewFunction("fn", []*Statement{
							NewStatement(
								NewReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
							),
						}),
					),
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(slices.Values(tc.tokens()))
			if tc.expectedErr != nil {
				if errors.Is(err, tc.expectedErr) {
					return
				}
			}

			if err != nil {
				t.Error(err)
				return
			}

			if diff := cmp.Diff(tc.expectedNode(), result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
