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
	program := NewProgram([]*Statement{})
	ctx := NewParsingCtx(tokens)
	defer ctx.Seq.Stop()

	for {
		if ctx.Seq.PeekNext() == nil {
			break
		}

		statement, err := parseStatement(ctx)
		if err != nil {
			return nil, err
		}

		program.AddStatement(statement)
	}

	return program, nil
}

func parseStatement(ctx *ParsingCtx) (*Statement, error) {
	expected, err := ctx.Seq.ExpectAny("Invalid statement",
		token.TypeKeywordDefine,
		token.TypeKeywordFree,
		token.TypeKeywordReturn,
	)

	if err != nil {
		return nil, err
	}

	switch expected := expected.(type) {
	case *token.KeywordDefine:
		return parseDefine(ctx)

	case *token.KeywordFree:
		return parseFree(ctx)

	case *token.KeywordReturn:
		return parseReturn(ctx)

	default:
		return nil, fmt.Errorf("%w: %s", UnexpectedTokenErr, expected.String())
	}
}

func parseDefine(ctx *ParsingCtx) (*Statement, error) {
	nameIdentifier, err := ExpectToken[*token.Identifier](ctx, "`def` should be followed by name")
	if err != nil {
		return nil, err
	}

	expected, err := ctx.Seq.ExpectAny("should be followed by assign or function curly brace",
		token.TypeOperatorAssign,
		token.TypeCurlyBracketOpen,
	)
	if err != nil {
		return nil, err
	}

	switch expected.(type) {
	case *token.OperatorAssign:
		return parseVariable(ctx, nameIdentifier)

	case *token.CurlyBracketOpen:
		return parseFunction(ctx, nameIdentifier)

	default:
		panic(errext.Tag("define", UnreachableErr))
	}
}

func parseVariable(ctx *ParsingCtx, nameIdentifier *token.Identifier) (*Statement, error) {
	expression, err := parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	variable := NewVariable(nameIdentifier.Name, expression)
	ctx.Scope.AddVariable(variable)

	return NewStatement(variable), nil
}

func parseFunction(ctx *ParsingCtx, nameIdentifier *token.Identifier) (*Statement, error) {
	ctx.DiveScope()
	defer ctx.AscendScope()

	var statements []*Statement

	for {
		statement, err := parseStatement(ctx)
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)

		if next := ctx.Seq.PeekNext(); next != nil && next.Type() == token.TypeCurlyBracketClose {
			ctx.Seq.Next()
			break
		}
	}

	return NewStatement(
		NewFunction(nameIdentifier.Name, statements),
	), nil
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
			variable, err := ctx.Scope.ExpectVariableDefined(expected)
			if err != nil {
				return nil, err
			}
			if variable.DataType() != DataFloat && variable.DataType() != DataInt {
				return nil, fmt.Errorf("%w: %s cannot be used in math expressions", UnexpectedTokenErr, variable.DataType())
			}

			addNode(NewVariableReference(variable))

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

	variable, err := ctx.Scope.ExpectVariableDefined(nameIdentifier)
	if err != nil {
		return nil, err
	}

	return NewStatement(NewFree(variable)), nil
}

func parseReturn(ctx *ParsingCtx) (*Statement, error) {
	expression, err := parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return NewStatement(NewReturn(expression)), nil
}
