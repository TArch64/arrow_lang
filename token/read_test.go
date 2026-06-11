package token

import (
	"slices"
	"strings"
	"testing"

	"arrow_lang/testutil"

	"github.com/google/go-cmp/cmp"
)

func assertRead(t *testing.T, text string, expected []Token) {
	t.Helper()
	result := slices.Collect(Read(strings.NewReader(testutil.Dedent(text))))
	if diff := cmp.Diff(result, expected); diff != "" {
		t.Error(diff)
	}
}

func TestReadEmptyInput(t *testing.T) {
	assertRead(t, "", nil)
}

func TestReadWhitespaceOnly(t *testing.T) {
	assertRead(t, "   \t\n  ", nil)
}

func TestReadDefineInt(t *testing.T) {
	assertRead(t, "def a = 1", []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
	})
}

func TestReadDefineNegativeInt(t *testing.T) {
	assertRead(t, "def a = -1", []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(-1),
	})
}

func TestReadDefineZeroInt(t *testing.T) {
	assertRead(t, "def a = 0", []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(0),
	})
}

func TestReadDefineLargeInt(t *testing.T) {
	assertRead(t, "def num = 999999999", []Token{
		NewKeywordDefine(),
		NewIdentifier("num"),
		NewOperatorAssign(),
		NewLiteralInt(999999999),
	})
}

func TestReadDefineFloat(t *testing.T) {
	assertRead(t, "def a = 1.123", []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralFloat(1.123),
	})
}

func TestReadDefineNegativeFloat(t *testing.T) {
	assertRead(t, "def a = -1.123", []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralFloat(-1.123),
	})
}

func TestReadDefineWholeNumberFloat(t *testing.T) {
	assertRead(t, "def val = 1.0", []Token{
		NewKeywordDefine(),
		NewIdentifier("val"),
		NewOperatorAssign(),
		NewLiteralFloat(1.0),
	})
}

func TestReadDefineLeadingZeroFloat(t *testing.T) {
	assertRead(t, "def small = 0.123", []Token{
		NewKeywordDefine(),
		NewIdentifier("small"),
		NewOperatorAssign(),
		NewLiteralFloat(0.123),
	})
}

func TestReadDefineHighPrecisionFloat(t *testing.T) {
	assertRead(t, "def pi = 3.14159265359", []Token{
		NewKeywordDefine(),
		NewIdentifier("pi"),
		NewOperatorAssign(),
		NewLiteralFloat(3.14159265359),
	})
}

func TestReadSingleIdentifier(t *testing.T) {
	assertRead(t, "variable", []Token{
		NewIdentifier("variable"),
	})
}

func TestReadIdentifierWithDigits(t *testing.T) {
	assertRead(t, "def var123 = 5", []Token{
		NewKeywordDefine(),
		NewIdentifier("var123"),
		NewOperatorAssign(),
		NewLiteralInt(5),
	})
}

func TestReadIdentifierWithUnderscores(t *testing.T) {
	assertRead(t, "def my_variable = 10", []Token{
		NewKeywordDefine(),
		NewIdentifier("my_variable"),
		NewOperatorAssign(),
		NewLiteralInt(10),
	})
}

func TestReadKeywordPrefixedIdentifiers(t *testing.T) {
	assertRead(t, "def define = 1 def definition = 2", []Token{
		NewKeywordDefine(),
		NewIdentifier("define"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordDefine(),
		NewIdentifier("definition"),
		NewOperatorAssign(),
		NewLiteralInt(2),
	})
}

func TestReadDefineFromVariable(t *testing.T) {
	assertRead(t, `
		def a = 1
		def b = a`, []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordDefine(),
		NewIdentifier("b"),
		NewOperatorAssign(),
		NewIdentifier("a"),
	})
}

func TestReadFreeVariable(t *testing.T) {
	assertRead(t, `
		def a = 1
		free a`, []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordFree(),
		NewIdentifier("a"),
	})
}

func TestReadAddition(t *testing.T) {
	assertRead(t, `def a = 1 + 2`, []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewOperatorPlus(),
		NewLiteralInt(2),
	})
}

func TestReadAdditionWithVariable(t *testing.T) {
	assertRead(t, `
		def a = 1
		def b = a + 2`, []Token{
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
	})
}

func TestReadChainedAddition(t *testing.T) {
	assertRead(t, "def result = 1 + 2 + 3", []Token{
		NewKeywordDefine(),
		NewIdentifier("result"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewOperatorPlus(),
		NewLiteralInt(2),
		NewOperatorPlus(),
		NewLiteralInt(3),
	})
}

func TestReadAdditionWithFloat(t *testing.T) {
	assertRead(t, "def result = 1 + 2.5", []Token{
		NewKeywordDefine(),
		NewIdentifier("result"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewOperatorPlus(),
		NewLiteralFloat(2.5),
	})
}

func TestReadSubtraction(t *testing.T) {
	assertRead(t, `def a = 5 - 2`, []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(5),
		NewOperatorMinus(),
		NewLiteralInt(2),
	})
}

func TestReadSubtractionWithVariableLeft(t *testing.T) {
	assertRead(t, `
		def a = 10
		def b = a - 3`, []Token{
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
	})
}

func TestReadSubtractionWithVariableRight(t *testing.T) {
	assertRead(t, `
		def a = 5
		def b = 10 - a`, []Token{
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
	})
}

func TestReadChainedSubtraction(t *testing.T) {
	assertRead(t, "def result = 10 - 3 - 2", []Token{
		NewKeywordDefine(),
		NewIdentifier("result"),
		NewOperatorAssign(),
		NewLiteralInt(10),
		NewOperatorMinus(),
		NewLiteralInt(3),
		NewOperatorMinus(),
		NewLiteralInt(2),
	})
}

func TestReadSubtractionFromZero(t *testing.T) {
	assertRead(t, "def negative = 0 - 5", []Token{
		NewKeywordDefine(),
		NewIdentifier("negative"),
		NewOperatorAssign(),
		NewLiteralInt(0),
		NewOperatorMinus(),
		NewLiteralInt(5),
	})
}

func TestReadSubtractionWithFloat(t *testing.T) {
	assertRead(t, "def result = 5 - 2.5", []Token{
		NewKeywordDefine(),
		NewIdentifier("result"),
		NewOperatorAssign(),
		NewLiteralInt(5),
		NewOperatorMinus(),
		NewLiteralFloat(2.5),
	})
}

func TestReadSubtractionBetweenFloats(t *testing.T) {
	assertRead(t, "def result = 3.14 - 1.5", []Token{
		NewKeywordDefine(),
		NewIdentifier("result"),
		NewOperatorAssign(),
		NewLiteralFloat(3.14),
		NewOperatorMinus(),
		NewLiteralFloat(1.5),
	})
}

func TestReadAdditionThenSubtraction(t *testing.T) {
	assertRead(t, "def result = 1 + 2 - 3", []Token{
		NewKeywordDefine(),
		NewIdentifier("result"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewOperatorPlus(),
		NewLiteralInt(2),
		NewOperatorMinus(),
		NewLiteralInt(3),
	})
}

func TestReadSubtractionThenAddition(t *testing.T) {
	assertRead(t, "def result = 10 - 3 + 2", []Token{
		NewKeywordDefine(),
		NewIdentifier("result"),
		NewOperatorAssign(),
		NewLiteralInt(10),
		NewOperatorMinus(),
		NewLiteralInt(3),
		NewOperatorPlus(),
		NewLiteralInt(2),
	})
}

func TestReadMultipleDefinitions(t *testing.T) {
	assertRead(t, "def x = 1 def y = 2.5 def z = x + y", []Token{
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
	})
}

func TestReadExpressionWithVariables(t *testing.T) {
	assertRead(t, `
		def x = 5
		def y = 3
		def z = x + y - 2`, []Token{
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
	})
}

func TestReadTrailingNewlines(t *testing.T) {
	assertRead(t, "def a = 1\n\n", []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
	})
}

func TestReadTabsAndNewlinesAsSeparators(t *testing.T) {
	assertRead(t, "def\ta\t=\n1\n+\n2", []Token{
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewOperatorPlus(),
		NewLiteralInt(2),
	})
}

func TestReadFunction(t *testing.T) {
	assertRead(t, `def fn() { ret 1 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionReturnZero(t *testing.T) {
	assertRead(t, `def zero() { ret 0 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("zero"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(0),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionReturnNegativeInt(t *testing.T) {
	assertRead(t, `def fn() { ret -1 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(-1),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionReturnFloat(t *testing.T) {
	assertRead(t, `def pi() { ret 3.14 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("pi"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralFloat(3.14),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionReturnAddition(t *testing.T) {
	assertRead(t, `def fn() { ret 1 + 2 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewOperatorPlus(),
		NewLiteralInt(2),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionReturnSubtraction(t *testing.T) {
	assertRead(t, `def fn() { ret 10 - 3 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(10),
		NewOperatorMinus(),
		NewLiteralInt(3),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionReturnMixedExpression(t *testing.T) {
	assertRead(t, `def fn() { ret 1 + 2 - 3 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewOperatorPlus(),
		NewLiteralInt(2),
		NewOperatorMinus(),
		NewLiteralInt(3),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionMultilineBody(t *testing.T) {
	assertRead(t, `
		def fn() {
			ret 42
		}`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(42),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionWithLocalVariable(t *testing.T) {
	assertRead(t, `def fn() { def x = 1 ret x }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordDefine(),
		NewIdentifier("x"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordReturn(),
		NewIdentifier("x"),
		NewCurlyBracketClose(),
	})
}

func TestReadMultipleFunctions(t *testing.T) {
	assertRead(t, `
		def one() { ret 1 }
		def two() { ret 2 }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("one"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("two"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(2),
		NewCurlyBracketClose(),
	})
}

func TestReadFunctionCall(t *testing.T) {
	assertRead(t, `
		def fn() { ret 1 }
		def a = fn()`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
	})
}

func TestReadCallPlusLiteral(t *testing.T) {
	assertRead(t, `
		def fn() { ret 1 }
		def a = fn() + 2`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewOperatorPlus(),
		NewLiteralInt(2),
	})
}

func TestReadLiteralPlusCall(t *testing.T) {
	assertRead(t, `
		def fn() { ret 1 }
		def a = 2 + fn()`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(2),
		NewOperatorPlus(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
	})
}

func TestReadCallMinusCall(t *testing.T) {
	assertRead(t, `
		def one() { ret 1 }
		def two() { ret 2 }
		def a = two() - one()`, []Token{
		NewKeywordDefine(),
		NewIdentifier("one"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("two"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(2),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewIdentifier("two"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewOperatorMinus(),
		NewIdentifier("one"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
	})
}

func TestReadCallPlusVariable(t *testing.T) {
	assertRead(t, `
		def fn() { ret 1 }
		def b = 5
		def a = fn() + b`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("b"),
		NewOperatorAssign(),
		NewLiteralInt(5),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewOperatorPlus(),
		NewIdentifier("b"),
	})
}

func TestReadChainedCallExpression(t *testing.T) {
	assertRead(t, `
		def fn() { ret 1 }
		def a = fn() + fn() - 2`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordReturn(),
		NewLiteralInt(1),
		NewCurlyBracketClose(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewOperatorPlus(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewOperatorMinus(),
		NewLiteralInt(2),
	})
}

func TestReadDeferFree(t *testing.T) {
	assertRead(t, `
		def fn() {
			def a = 1
			defer free a
			ret a
		}`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordDefer(),
		NewKeywordFree(),
		NewIdentifier("a"),
		NewKeywordReturn(),
		NewIdentifier("a"),
		NewCurlyBracketClose(),
	})
}

func TestReadDeferInlineFunctionBody(t *testing.T) {
	assertRead(t, `def fn() { def a = 1 defer free a ret a }`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordDefer(),
		NewKeywordFree(),
		NewIdentifier("a"),
		NewKeywordReturn(),
		NewIdentifier("a"),
		NewCurlyBracketClose(),
	})
}

func TestReadMultipleDefers(t *testing.T) {
	assertRead(t, `
		def fn() {
			def a = 1
			def b = 2
			defer free a
			defer free b
			ret a
		}`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordDefine(),
		NewIdentifier("b"),
		NewOperatorAssign(),
		NewLiteralInt(2),
		NewKeywordDefer(),
		NewKeywordFree(),
		NewIdentifier("a"),
		NewKeywordDefer(),
		NewKeywordFree(),
		NewIdentifier("b"),
		NewKeywordReturn(),
		NewIdentifier("a"),
		NewCurlyBracketClose(),
	})
}

func TestReadDeferTabSeparated(t *testing.T) {
	assertRead(t, "def fn() {\n\tdef a = 1\n\tdefer\tfree\ta\n\tret a\n}", []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordDefer(),
		NewKeywordFree(),
		NewIdentifier("a"),
		NewKeywordReturn(),
		NewIdentifier("a"),
		NewCurlyBracketClose(),
	})
}

func TestReadDeferMultilineSeparated(t *testing.T) {
	assertRead(t, `
		def fn() {
			def a = 1
			defer
			free a
			ret a
		}`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordDefine(),
		NewIdentifier("a"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordDefer(),
		NewKeywordFree(),
		NewIdentifier("a"),
		NewKeywordReturn(),
		NewIdentifier("a"),
		NewCurlyBracketClose(),
	})
}

func TestReadDeferKeywordPrefixedIdentifier(t *testing.T) {
	assertRead(t, `
		def fn() {
			def deferred = 1
			ret deferred
		}`, []Token{
		NewKeywordDefine(),
		NewIdentifier("fn"),
		NewParenthesesOpen(),
		NewParenthesesClose(),
		NewCurlyBracketOpen(),
		NewKeywordDefine(),
		NewIdentifier("deferred"),
		NewOperatorAssign(),
		NewLiteralInt(1),
		NewKeywordReturn(),
		NewIdentifier("deferred"),
		NewCurlyBracketClose(),
	})
}
