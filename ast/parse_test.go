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
		tokens       []token.Token
		expectedErr  error
		expectedNode Node
	}

	testCases := []testCase{
		{
			name:        "error_cases/define_eof_after_keyword",
			tokens:      []token.Token{token.NewKeywordDefine()},
			expectedErr: UnexpectedEOFErr,
		},
		{
			name: "error_cases/define_invalid_token_after_keyword",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewOperatorAssign(),
			},
			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "error_cases/define_eof_after_identifier",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
			},
			expectedErr: UnexpectedEOFErr,
		},
		{
			name: "error_cases/define_invalid_token_after_identifier",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorPlus(),
			},
			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "error_cases/expression_incomplete_sum",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
				token.NewOperatorAssign(),
			},
			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "error_cases/undefined_variable_reference",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
				token.NewKeywordDefine(),
				token.NewIdentifier("b"),
				token.NewOperatorAssign(),
				token.NewIdentifier("c"),
			},
			expectedErr: UndefinedVariableErr,
		},
		{
			name: "basic/define_int",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("a",
						NewExpression(NewLiteralInt(1)),
					),
				),
			),
		},
		{
			name: "basic/define_negative_int",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(-1),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("a",
						NewExpression(NewLiteralInt(-1)),
					),
				),
			),
		},
		{
			name: "basic/define_float",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralFloat(1.123),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("a",
						NewExpression(NewLiteralFloat(1.123)),
					),
				),
			),
		},
		{
			name: "basic/define_negative_float",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralFloat(-1.123),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("a",
						NewExpression(NewLiteralFloat(-1.123)),
					),
				),
			),
		},
		{
			name: "basic/free_variable",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
				token.NewKeywordFree(),
				token.NewIdentifier("a"),
			},
			expectedNode: func() Node {
				defA := NewDefine("a",
					NewExpression(NewLiteralInt(1)),
				)
				return NewProgram(
					NewStatement(defA),
					NewStatement(NewFree(defA)),
				)
			}(),
		},
		{
			name: "expressions/sum_two_integers",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
				token.NewOperatorPlus(),
				token.NewLiteralInt(2),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("a",
						NewExpression(NewLiteralInt(3)),
					),
				),
			),
		},
		{
			name: "variables/assign_variable_to_variable",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
				token.NewKeywordDefine(),
				token.NewIdentifier("b"),
				token.NewOperatorAssign(),
				token.NewIdentifier("a"),
			},
			expectedNode: func() Node {
				defA := NewDefine("a",
					NewExpression(NewLiteralInt(1)),
				)
				return NewProgram(
					NewStatement(defA),
					NewStatement(
						NewDefine("b",
							NewExpression(NewVariableReference(defA)),
						),
					),
				)
			}(),
		},
		{
			name: "expressions/sum_variable_and_literal",
			tokens: []token.Token{
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
			},
			expectedNode: func() Node {
				defA := NewDefine("a", NewExpression(NewLiteralInt(1)))
				return NewProgram(
					NewStatement(defA),
					NewStatement(
						NewDefine("b",
							NewExpression(
								NewExpressionSum(
									NewVariableReference(defA),
									NewLiteralInt(2),
								),
							),
						),
					),
				)
			}(),
		},
		{
			name: "numbers/zero_integer",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("zero"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(0),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("zero",
						NewExpression(NewLiteralInt(0)),
					),
				),
			),
		},
		{
			name: "numbers/large_integer",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("big"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(999999999),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("big",
						NewExpression(NewLiteralInt(999999999)),
					),
				),
			),
		},
		{
			name: "numbers/float_zero_fractional",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("val"),
				token.NewOperatorAssign(),
				token.NewLiteralFloat(5.0),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("val",
						NewExpression(NewLiteralFloat(5.0)),
					),
				),
			),
		},
		{
			name: "numbers/float_leading_zero",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("small"),
				token.NewOperatorAssign(),
				token.NewLiteralFloat(0.5),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("small",
						NewExpression(NewLiteralFloat(0.5)),
					),
				),
			),
		},
		{
			name: "expressions/sum_two_floats",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("result"),
				token.NewOperatorAssign(),
				token.NewLiteralFloat(1.5),
				token.NewOperatorPlus(),
				token.NewLiteralFloat(2.5),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("result",
						NewExpression(NewLiteralFloat(4.0)),
					),
				),
			),
		},
		{
			name: "expressions/sum_mixed_int_float",
			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("mixed"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
				token.NewOperatorPlus(),
				token.NewLiteralFloat(2.5),
			},
			expectedNode: NewProgram(
				NewStatement(
					NewDefine("mixed",
						NewExpression(NewLiteralFloat(3.5)),
					),
				),
			),
		},
		{
			name: "complex/multiple_statements",
			tokens: []token.Token{
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
			},
			expectedNode: func() Node {
				defX := NewDefine("x", NewExpression(NewLiteralInt(10)))
				defY := NewDefine("y", NewExpression(NewLiteralInt(20)))
				return NewProgram(
					NewStatement(defX),
					NewStatement(defY),
					NewStatement(
						NewDefine("sum",
							NewExpression(
								NewExpressionSum(
									NewVariableReference(defX),
									NewVariableReference(defY),
								),
							),
						),
					),
				)
			}(),
		},
		{
			name: "complex/define_and_free_sequence",
			tokens: []token.Token{
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
			},
			expectedNode: func() Node {
				defTemp := NewDefine("temp", NewExpression(NewLiteralInt(42)))
				return NewProgram(
					NewStatement(defTemp),
					NewStatement(
						NewDefine("copy",
							NewExpression(NewVariableReference(defTemp)),
						),
					),
					NewStatement(NewFree(defTemp)),
				)
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(slices.Values(tc.tokens))
			if tc.expectedErr != nil {
				if errors.Is(err, tc.expectedErr) {
					return
				}
			}

			if err != nil {
				t.Error(err)
				return
			}

			if diff := cmp.Diff(tc.expectedNode, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
