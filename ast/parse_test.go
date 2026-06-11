package ast

import (
	"errors"
	"slices"
	"testing"

	"arrow_lang/token"

	"github.com/google/go-cmp/cmp"
)

func assertParse(t *testing.T, tokens []token.Token, expected Node) {
	t.Helper()
	result, err := Parse(slices.Values(tokens))
	if err != nil {
		t.Error(err)
		return
	}
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func assertParseError(t *testing.T, tokens []token.Token, expected error) {
	t.Helper()
	if _, err := Parse(slices.Values(tokens)); !errors.Is(err, expected) {
		t.Errorf("expected error %v, got %v", expected, err)
	}
}

func TestParseEmpty(t *testing.T) {
	assertParse(t, []token.Token{}, NewProgram([]*Statement{}))
}

func TestParseDefineInt(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralInt(1)}),
			),
		),
	}))
}

func TestParseDefineNegativeInt(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(-1),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralInt(-1)}),
			),
		),
	}))
}

func TestParseDefineZeroInt(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("zero"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(0),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("zero",
				NewExpression([]DataNode{NewLiteralInt(0)}),
			),
		),
	}))
}

func TestParseDefineLargeInt(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("big"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(999999999),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("big",
				NewExpression([]DataNode{NewLiteralInt(999999999)}),
			),
		),
	}))
}

func TestParseDefineFloat(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(1.123),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralFloat(1.123)}),
			),
		),
	}))
}

func TestParseDefineNegativeFloat(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(-1.123),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralFloat(-1.123)}),
			),
		),
	}))
}

func TestParseDefineWholeNumberFloat(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("val"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(5.0),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("val",
				NewExpression([]DataNode{NewLiteralFloat(5.0)}),
			),
		),
	}))
}

func TestParseDefineLeadingZeroFloat(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("small"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(0.5),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("small",
				NewExpression([]DataNode{NewLiteralFloat(0.5)}),
			),
		),
	}))
}

func TestParseDefineFromVariable(t *testing.T) {
	defA := NewVariable("a",
		NewExpression([]DataNode{NewLiteralInt(1)}),
	)
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordDefine(),
		token.NewIdentifier("b"),
		token.NewOperatorAssign(),
		token.NewIdentifier("a"),
	}, NewProgram([]*Statement{
		NewStatement(defA),
		NewStatement(
			NewVariable("b",
				NewExpression([]DataNode{NewVariableReference(defA)}),
			),
		),
	}))
}

func TestParseDefineEOFAfterKeyword(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
	}, UnexpectedEOFErr)
}

func TestParseDefineInvalidTokenAfterKeyword(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewOperatorAssign(),
	}, UnexpectedTokenErr)
}

func TestParseDefineEOFAfterIdentifier(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
	}, UnexpectedEOFErr)
}

func TestParseDefineInvalidTokenAfterIdentifier(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorPlus(),
	}, UnexpectedTokenErr)
}

func TestParseStatementInvalidFirstToken(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
	}, UnexpectedTokenErr)
}

func TestParseUndefinedVariableReference(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordDefine(),
		token.NewIdentifier("b"),
		token.NewOperatorAssign(),
		token.NewIdentifier("c"),
	}, UndefinedVariableErr)
}

func TestParseFreeVariable(t *testing.T) {
	defA := NewVariable("a",
		NewExpression([]DataNode{NewLiteralInt(1)}),
	)
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordFree(),
		token.NewIdentifier("a"),
	}, NewProgram([]*Statement{
		NewStatement(defA),
		NewStatement(NewFree(defA)),
	}))
}

func TestParseDefineThenFree(t *testing.T) {
	defTemp := NewVariable("temp", NewExpression([]DataNode{NewLiteralInt(42)}))
	assertParse(t, []token.Token{
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
	}, NewProgram([]*Statement{
		NewStatement(defTemp),
		NewStatement(
			NewVariable("copy",
				NewExpression([]DataNode{NewVariableReference(defTemp)}),
			),
		),
		NewStatement(NewFree(defTemp)),
	}))
}

func TestParseFreeEOFAfterKeyword(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordFree(),
	}, UnexpectedEOFErr)
}

func TestParseFreeInvalidTokenAfterKeyword(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordFree(),
		token.NewLiteralInt(1),
	}, UnexpectedTokenErr)
}

func TestParseFreeUndefinedVariable(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordFree(),
		token.NewIdentifier("missing"),
	}, UndefinedVariableErr)
}

func TestParseFreeFunctionName(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordFree(),
		token.NewIdentifier("fn"),
	}, UndefinedVariableErr)
}

func TestParseAddIntegers(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralInt(3)}),
			),
		),
	}))
}

func TestParseAddFloats(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(1.5),
		token.NewOperatorPlus(),
		token.NewLiteralFloat(2.5),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralFloat(4.0)}),
			),
		),
	}))
}

func TestParseAddMixedIntFloat(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("mixed"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewOperatorPlus(),
		token.NewLiteralFloat(2.5),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("mixed",
				NewExpression([]DataNode{NewLiteralFloat(3.5)}),
			),
		),
	}))
}

func TestParseAddFloatAndInt(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(1.5),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralFloat(3.5)}),
			),
		),
	}))
}

func TestParseAddVariableAndLiteral(t *testing.T) {
	defA := NewVariable("a", NewExpression([]DataNode{NewLiteralInt(1)}))
	assertParse(t, []token.Token{
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
	}, NewProgram([]*Statement{
		NewStatement(defA),
		NewStatement(
			NewVariable("b",
				NewExpression([]DataNode{
					NewVariableReference(defA),
					NewExpressionPlus(NewLiteralInt(2)),
				}),
			),
		),
	}))
}

func TestParseAddChainedLiterals(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
		token.NewOperatorPlus(),
		token.NewLiteralInt(3),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralInt(6)}),
			),
		),
	}))
}

func TestParseAddChainedMixedToFloat(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
		token.NewOperatorPlus(),
		token.NewLiteralFloat(1.5),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralFloat(4.5)}),
			),
		),
	}))
}

func TestParseSubtractIntegers(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(5),
		token.NewOperatorMinus(),
		token.NewLiteralInt(3),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(2)}),
			),
		),
	}))
}

func TestParseSubtractFloats(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(7.5),
		token.NewOperatorMinus(),
		token.NewLiteralFloat(2.25),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralFloat(5.25)}),
			),
		),
	}))
}

func TestParseSubtractIntFromFloat(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(10.5),
		token.NewOperatorMinus(),
		token.NewLiteralInt(5),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralFloat(5.5)}),
			),
		),
	}))
}

func TestParseSubtractFloatFromInt(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(8),
		token.NewOperatorMinus(),
		token.NewLiteralFloat(3.5),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralFloat(4.5)}),
			),
		),
	}))
}

func TestParseSubtractZero(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(42),
		token.NewOperatorMinus(),
		token.NewLiteralInt(0),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(42)}),
			),
		),
	}))
}

func TestParseSubtractEqualValues(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(7),
		token.NewOperatorMinus(),
		token.NewLiteralInt(7),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(0)}),
			),
		),
	}))
}

func TestParseSubtractToNegative(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(3),
		token.NewOperatorMinus(),
		token.NewLiteralInt(8),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(-5)}),
			),
		),
	}))
}

func TestParseSubtractFromNegative(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(-10),
		token.NewOperatorMinus(),
		token.NewLiteralInt(5),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(-15)}),
			),
		),
	}))
}

func TestParseSubtractNegativeOperand(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(10),
		token.NewOperatorMinus(),
		token.NewLiteralInt(-3),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(13)}),
			),
		),
	}))
}

func TestParseSubtractFloatPrecision(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralFloat(3.14159),
		token.NewOperatorMinus(),
		token.NewLiteralFloat(2.71828),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralFloat(0.42330999999999985)}),
			),
		),
	}))
}

func TestParseSubtractChainedLiterals(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(100),
		token.NewOperatorMinus(),
		token.NewLiteralInt(25),
		token.NewOperatorMinus(),
		token.NewLiteralInt(15),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(60)}),
			),
		),
	}))
}

func TestParseSubtractChainedThreeLiterals(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(10),
		token.NewOperatorMinus(),
		token.NewLiteralInt(3),
		token.NewOperatorMinus(),
		token.NewLiteralInt(2),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{NewLiteralInt(5)}),
			),
		),
	}))
}

func TestParseSubtractVariableFromLiteral(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(15)}))
	assertParse(t, []token.Token{
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
	}, NewProgram([]*Statement{
		NewStatement(defX),
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{
					NewLiteralInt(20),
					NewExpressionMinus(NewVariableReference(defX)),
				}),
			),
		),
	}))
}

func TestParseSubtractLiteralFromVariable(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(25)}))
	assertParse(t, []token.Token{
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
	}, NewProgram([]*Statement{
		NewStatement(defX),
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{
					NewVariableReference(defX),
					NewExpressionMinus(NewLiteralInt(12)),
				}),
			),
		),
	}))
}

func TestParseSubtractVariableFromVariable(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(30)}))
	defY := NewVariable("y", NewExpression([]DataNode{NewLiteralInt(18)}))
	assertParse(t, []token.Token{
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
	}, NewProgram([]*Statement{
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
	}))
}

func TestParseMixedPlusMinusLiterals(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("result"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(10),
		token.NewOperatorPlus(),
		token.NewLiteralInt(5),
		token.NewOperatorMinus(),
		token.NewLiteralInt(3),
	}, NewProgram([]*Statement{
		NewStatement(
			NewVariable("result",
				NewExpression([]DataNode{NewLiteralInt(12)}),
			),
		),
	}))
}

func TestParseMixedPlusMinusWithVariables(t *testing.T) {
	defA := NewVariable("a", NewExpression([]DataNode{NewLiteralInt(50)}))
	defB := NewVariable("b", NewExpression([]DataNode{NewLiteralInt(20)}))
	defC := NewVariable("c", NewExpression([]DataNode{NewLiteralInt(8)}))
	assertParse(t, []token.Token{
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
	}, NewProgram([]*Statement{
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
	}))
}

func TestParseFoldLiteralsTrailingVariable(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(1)}))
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("x"),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
		token.NewOperatorPlus(),
		token.NewLiteralInt(3),
	}, NewProgram([]*Statement{
		NewStatement(defX),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewVariableReference(defX),
					NewExpressionPlus(NewLiteralInt(5)),
				}),
			),
		),
	}))
}

func TestParseNoFoldLiteralsAcrossVariable(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(1)}))
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewOperatorPlus(),
		token.NewIdentifier("x"),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
	}, NewProgram([]*Statement{
		NewStatement(defX),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewLiteralInt(1),
					NewExpressionPlus(NewVariableReference(defX)),
					NewExpressionPlus(NewLiteralInt(2)),
				}),
			),
		),
	}))
}

func TestParseFoldMixedOperatorsTrailingVariable(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(1)}))
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("x"),
		token.NewOperatorPlus(),
		token.NewLiteralInt(1),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
		token.NewOperatorMinus(),
		token.NewLiteralInt(3),
	}, NewProgram([]*Statement{
		NewStatement(defX),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewVariableReference(defX),
					NewExpressionPlus(NewLiteralInt(0)),
				}),
			),
		),
	}))
}

func TestParseMultipleStatements(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(10)}))
	defY := NewVariable("y", NewExpression([]DataNode{NewLiteralInt(20)}))
	assertParse(t, []token.Token{
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
	}, NewProgram([]*Statement{
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
	}))
}

func TestParseExpressionEOFAfterAssign(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
	}, UnexpectedEOFErr)
}

func TestParseExpressionInvalidTokenAfterAssign(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewKeywordDefine(),
	}, UnexpectedTokenErr)
}

func TestParseExpressionEOFAfterOperator(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewOperatorPlus(),
	}, UnexpectedEOFErr)
}

func TestParseUnexpectedAssignAfterValue(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewOperatorAssign(),
	}, UnexpectedTokenErr)
}

func TestParseFunction(t *testing.T) {
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(
			NewFunction("fn", []*Statement{
				NewStatement(
					NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
				),
			}),
		),
	}))
}

func TestParseFunctionWithLocalVariable(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(1)}))
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordReturn(),
		token.NewIdentifier("x"),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(
			NewFunction("fn", []*Statement{
				NewStatement(defX),
				NewStatement(
					NewFunctionReturn(
						NewExpression([]DataNode{NewVariableReference(defX)}),
					),
				),
			}),
		),
	}))
}

func TestParseFunctionWithFree(t *testing.T) {
	defA := NewVariable("a", NewExpression([]DataNode{NewLiteralInt(1)}))
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordFree(),
		token.NewIdentifier("a"),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(
			NewFunction("fn", []*Statement{
				NewStatement(defA),
				NewStatement(NewFree(defA)),
				NewStatement(
					NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
				),
			}),
		),
	}))
}

func TestParseNestedFunction(t *testing.T) {
	inner := NewFunction("inner", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("outer"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefine(),
		token.NewIdentifier("inner"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordReturn(),
		token.NewIdentifier("inner"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(
			NewFunction("outer", []*Statement{
				NewStatement(inner),
				NewStatement(
					NewFunctionReturn(NewExpression([]DataNode{NewFunctionCall(inner)})),
				),
			}),
		),
	}))
}

func TestParseFunctionReferencesOuterVariable(t *testing.T) {
	defX := NewVariable("x", NewExpression([]DataNode{NewLiteralInt(5)}))
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(5),
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewIdentifier("x"),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(defX),
		NewStatement(
			NewFunction("fn", []*Statement{
				NewStatement(
					NewFunctionReturn(
						NewExpression([]DataNode{NewVariableReference(defX)}),
					),
				),
			}),
		),
	}))
}

func TestParseReturnEOFAfterKeyword(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordReturn(),
	}, UnexpectedEOFErr)
}

func TestParseFunctionEmptyBody(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewCurlyBracketClose(),
	}, UnexpectedTokenErr)
}

func TestParseFunctionLocalVariableNotVisibleOutside(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordReturn(),
		token.NewIdentifier("x"),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("x"),
	}, UndefinedVariableErr)
}

func TestParseFunctionCall(t *testing.T) {
	function := NewFunction("fn", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
	}, NewProgram([]*Statement{
		NewStatement(function),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewFunctionCall(function),
				}),
			),
		),
	}))
}

func TestParseCallPlusLiteral(t *testing.T) {
	function := NewFunction("fn", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
	}, NewProgram([]*Statement{
		NewStatement(function),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewFunctionCall(function),
					NewExpressionPlus(NewLiteralInt(2)),
				}),
			),
		),
	}))
}

func TestParseLiteralPlusCall(t *testing.T) {
	function := NewFunction("fn", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(2),
		token.NewOperatorPlus(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
	}, NewProgram([]*Statement{
		NewStatement(function),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewLiteralInt(2),
					NewExpressionPlus(NewFunctionCall(function)),
				}),
			),
		),
	}))
}

func TestParseCallMinusCall(t *testing.T) {
	one := NewFunction("one", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	two := NewFunction("two", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(2)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("one"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("two"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(2),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("two"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewOperatorMinus(),
		token.NewIdentifier("one"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
	}, NewProgram([]*Statement{
		NewStatement(one),
		NewStatement(two),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewFunctionCall(two),
					NewExpressionMinus(NewFunctionCall(one)),
				}),
			),
		),
	}))
}

func TestParseCallPlusVariable(t *testing.T) {
	function := NewFunction("fn", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	defB := NewVariable("b", NewExpression([]DataNode{NewLiteralInt(5)}))
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("b"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(5),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewOperatorPlus(),
		token.NewIdentifier("b"),
	}, NewProgram([]*Statement{
		NewStatement(function),
		NewStatement(defB),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewFunctionCall(function),
					NewExpressionPlus(NewVariableReference(defB)),
				}),
			),
		),
	}))
}

func TestParseChainedCallExpression(t *testing.T) {
	function := NewFunction("fn", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewOperatorPlus(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewOperatorMinus(),
		token.NewLiteralInt(2),
	}, NewProgram([]*Statement{
		NewStatement(function),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewFunctionCall(function),
					NewExpressionPlus(NewFunctionCall(function)),
					NewExpressionMinus(NewLiteralInt(2)),
				}),
			),
		),
	}))
}

func TestParseFunctionCallsOuterFunction(t *testing.T) {
	one := NewFunction("one", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("one"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("two"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewIdentifier("one"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(one),
		NewStatement(
			NewFunction("two", []*Statement{
				NewStatement(
					NewFunctionReturn(
						NewExpression([]DataNode{NewFunctionCall(one)}),
					),
				),
			}),
		),
	}))
}

func TestParseFloatReturningCallInExpression(t *testing.T) {
	function := NewFunction("fn", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralFloat(1.5)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralFloat(1.5),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewOperatorPlus(),
		token.NewLiteralInt(2),
	}, NewProgram([]*Statement{
		NewStatement(function),
		NewStatement(
			NewVariable("a",
				NewExpression([]DataNode{
					NewFunctionCall(function),
					NewExpressionPlus(NewLiteralInt(2)),
				}),
			),
		),
	}))
}

func TestParseCallUndefinedFunction(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
	}, UndefinedFunctionErr)
}

func TestParseFunctionUsedAsVariable(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
	}, UndefinedVariableErr)
}

func TestParseCallVoidFunctionInExpression(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
	}, UnexpectedTokenErr)
}

func TestParseDeferFree(t *testing.T) {
	defA := NewVariable("a", NewExpression([]DataNode{
		NewLiteralInt(1),
	}))

	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefine(),
		token.NewIdentifier("a"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordDefer(),
		token.NewKeywordFree(),
		token.NewIdentifier("a"),
		token.NewKeywordReturn(),
		token.NewIdentifier("a"),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(
			NewFunction("fn", []*Statement{
				NewStatement(defA),
				NewStatement(NewDefer(NewStatement(NewFree(defA)))),
				NewStatement(
					NewFunctionReturn(
						NewExpression([]DataNode{NewVariableReference(defA)}),
					),
				),
			}),
		),
	}))
}

func TestParseDeferFunctionDefinition(t *testing.T) {
	inner := NewFunction("inner", []*Statement{
		NewStatement(
			NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
		),
	})
	assertParse(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("outer"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefer(),
		token.NewKeywordDefine(),
		token.NewIdentifier("inner"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
	}, NewProgram([]*Statement{
		NewStatement(
			NewFunction("outer", []*Statement{
				NewStatement(NewDefer(NewStatement(NewFunctionCall(inner)))),
				NewStatement(
					NewFunctionReturn(NewExpression([]DataNode{NewLiteralInt(1)})),
				),
			}),
		),
	}))
}

func TestParseDeferTopLevel(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefer(),
		token.NewKeywordFree(),
		token.NewIdentifier("a"),
	}, UnexpectedTokenErr)
}

func TestParseDeferVariableDefinition(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefer(),
		token.NewKeywordDefine(),
		token.NewIdentifier("x"),
		token.NewOperatorAssign(),
		token.NewLiteralInt(1),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
	}, UnexpectedTokenErr)
}

func TestParseDeferReturn(t *testing.T) {
	assertParseError(t, []token.Token{
		token.NewKeywordDefine(),
		token.NewIdentifier("fn"),
		token.NewParenthesesOpen(),
		token.NewParenthesesClose(),
		token.NewCurlyBracketOpen(),
		token.NewKeywordDefer(),
		token.NewKeywordReturn(),
		token.NewLiteralInt(1),
		token.NewCurlyBracketClose(),
	}, UnexpectedTokenErr)
}
