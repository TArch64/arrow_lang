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
	for {
		expected, err := ctx.Seq.ExpectAny("should be an expression",
			token.TypeLiteralInt,
			token.TypeLiteralFloat,
		)
		if err != nil {
			return nil, err
		}

		switch expected := expected.(type) {
		case *token.LiteralInt:
			expression.PlusInt(expected.Value)

		case *token.LiteralFloat:
			expression.PlusFloat(expected.Value)

		default:
			panic("unreachable")
		}

		if ctx.Seq.HasNextAny(token.TypeOperatorPlus) {
			ctx.Seq.Next()
		} else {
			break
		}
	}

	return expression.Build(), nil
}

func parseFree(ctx *ParsingCtx) (*Statement, error) {
	nameIdentifier, err := ExpectToken[*token.Identifier](ctx, "`free` should be followed by variable name")
	if err != nil {
		return nil, err
	}

	if !ctx.IsDefined(nameIdentifier.Name) {
		return nil, fmt.Errorf("%w: %s", UndefinedVariableErr, nameIdentifier.Name)
	}

	return NewStatement(NewFree(nameIdentifier.Name)), nil
}
