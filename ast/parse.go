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

func parseExpression(ctx *ParsingCtx) (*Expression, error) {
	expected, err := expectAnyToken(ctx, "should be an expression",
		&token.LiteralInt{},
		&token.LiteralFloat{},
	)
	if err != nil {
		return nil, err
	}

	switch expected := expected.(type) {
	case *token.LiteralInt:
		return NewExpression(NewLiteralInt(expected.Value)), nil

	case *token.LiteralFloat:
		return NewExpression(NewLiteralFloat(expected.Value)), nil

	default:
		panic("unreachable")
	}
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
