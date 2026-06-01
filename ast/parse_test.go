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
			name: "define: variable with literal int",

			tokens: []token.Token{
				token.NewKeywordDefine(),
				token.NewIdentifier("a"),
				token.NewOperatorAssign(),
				token.NewLiteralInt(1),
			},

			expectedNode: NewProgram(
				NewStatement(
					NewDefine("a", NewExpression(NewLiteralInt(1))),
				),
			),
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
