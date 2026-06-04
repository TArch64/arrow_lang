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
	defer parsingCtx.stop()

	for {
		t, ok := parsingCtx.next()
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
	nameIdentifier, err := expectToken[*token.Identifier](ctx, "`def` should be followed by name")
	if err != nil {
		return nil, err
	}

	_, err = expectToken[*token.OperatorAssign](ctx, "`def` should be followed by assign")
	if err != nil {
		return nil, err
	}

	expression, err := parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	define := NewDefine(nameIdentifier.Name, expression)
	ctx.addDefine(define)

	return NewStatement(define), nil
}

func expectToken[T token.Token](ctx *ParsingCtx, explain string) (T, error) {
	t, ok := ctx.next()
	var typed T

	if !ok {
		return typed, UnexpectedEOFErr
	}

	if t.Type() != typed.Type() {
		return typed, fmt.Errorf("%w: %s, got: %s", UnexpectedTokenErr, explain, t.String())
	}

	typed = t.(T)
	return typed, nil
}

func parseExpression(ctx *ParsingCtx) (*Expression, error) {
	literalInt, err := expectToken[*token.LiteralInt](ctx, "should be expression")
	if err != nil {
		return nil, err
	}

	return NewExpression(
		NewLiteralInt(literalInt.Value),
	), nil
}

func parseFree(ctx *ParsingCtx) (*Statement, error) {
	nameIdentifier, err := expectToken[*token.Identifier](ctx, "`free` should be followed by variable name")
	if err != nil {
		return nil, err
	}

	if !ctx.isDefined(nameIdentifier.Name) {
		return nil, fmt.Errorf("%w: %s", UndefinedVariableErr, nameIdentifier.Name)
	}

	return NewStatement(NewFree(nameIdentifier.Name)), nil
}
