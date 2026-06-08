package ast

import (
	"errors"
	"fmt"
	"iter"

	"arrow_lang/errext"
	"arrow_lang/token"
)

type NextToken func() (token.Token, bool)

var (
	UnexpectedTokenErr   = errext.Tag("ast", errors.New("unexpected token"))
	UnexpectedEOFErr     = errext.Tag("ast", errors.New("unexpected EOF"))
	UndefinedVariableErr = errext.Tag("ast", errors.New("undefined variable"))
	UnreachableErr       = errext.Tag("ast", errext.UnreachableErr)
)

func Parse(tokens iter.Seq[token.Token]) (*Program, error) {
	program := NewProgram()
	parsingCtx := NewParsingCtx(tokens)
	defer parsingCtx.Seq.Stop()

	for {
		t, ok := parsingCtx.Seq.Next()
		if !ok {
			break
		}

		switch t.Type() {
		case token.TypeKeywordDefine:
			statement, err := parseDefine(parsingCtx)
			if err != nil {
				return nil, err
			}

			program.AddStatement(statement)

		case token.TypeKeywordFree:
			statement, err := parseFree(parsingCtx)
			if err != nil {
				return nil, err
			}

			program.AddStatement(statement)
		default:
			return nil, fmt.Errorf("%w: %s", UnexpectedTokenErr, t.String())
		}
	}

	return program, nil
}

func parseDefine(ctx *ParsingCtx) (*Statement, error) {
	nameIdentifier, err := ExpectToken[*token.Identifier](ctx, "`def` should be followed by name")
	if err != nil {
		return nil, err
	}

	_, err = ExpectToken[*token.OperatorAssign](ctx, "`def` should be followed by assign")
	if err != nil {
		return nil, err
	}

	expression, err := parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	define := NewDefine(nameIdentifier.Name, expression)
	ctx.AddDefine(define)

	return NewStatement(define), nil
}

func parseExpression(ctx *ParsingCtx) (*Expression, error) {
	expression := NewParsingExpression()
	var addNode ParsingExpressionAdd = expression.Open

	for {
		expected, err := ctx.Seq.ExpectAny("should be an expression",
			token.TypeLiteralInt,
			token.TypeLiteralFloat,
			token.TypeIdentifier,
		)
		if err != nil {
			return nil, err
		}

		switch expected := expected.(type) {
		case *token.LiteralInt:
			addNode(NewLiteralInt(expected.Value))

		case *token.LiteralFloat:
			addNode(NewLiteralFloat(expected.Value))

		case *token.Identifier:
			define, err := ctx.ExpectDefined(expected)
			if err != nil {
				return nil, err
			}
			if define.DataType() != DataFloat && define.DataType() != DataInt {
				return nil, fmt.Errorf("%w: %s cannot be used in math expressions", UnexpectedTokenErr, define.DataType())
			}

			addNode(NewVariableReference(define))

		default:
			panic(errext.Tag("expression", UnreachableErr))
		}

		switch ctx.Seq.PeekNext().(type) {
		case *token.OperatorPlus:
			ctx.Seq.Next()
			addNode = expression.Plus

		case *token.OperatorMinus:
			ctx.Seq.Next()
			addNode = expression.Minus

		default:
			return expression.Build(), nil
		}
	}

	return expression.Build(), nil
}

func parseFree(ctx *ParsingCtx) (*Statement, error) {
	nameIdentifier, err := ExpectToken[*token.Identifier](ctx, "`free` should be followed by variable name")
	if err != nil {
		return nil, err
	}

	define, err := ctx.ExpectDefined(nameIdentifier)
	if err != nil {
		return nil, err
	}

	return NewStatement(NewFree(define)), nil
}
