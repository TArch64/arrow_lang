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
			name:        "define: eof after def",
			tokens:      []token.Token{token.NewKeywordDefine()},
			expectedErr: UnexpectedEOFErr,
		},
		{
			name: "define: invalid token after def",

			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewOperatorAssign(),
			},

			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "define: eof after name",

			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
			},

			expectedErr: UnexpectedEOFErr,
		},
		{
			name: "define: invalid token after name",

			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorPlus(),
			},

			expectedErr: UnexpectedTokenErr,
		},
		{
			name: "define: variable with int",

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
			name: "define: variable with negative int",

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
			name: "define: variable with float",

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
			name: "define: variable with negative float",

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
			name: "free: valid syntax",

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
			name: "define: variable with incomplete sum",

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
			name: "define: variable with sum 2 ints",

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
			name: "define: assign undefined variable",

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
			name: "define: assign another variable",

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
			name: "define: variable with sum of another variable and literal",

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
