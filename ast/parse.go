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
	UnexpectedTokenErr = errext.Tag("ast", errors.New("unexpected token"))
	UnexpectedEOFErr   = errext.Tag("ast", errors.New("unexpected EOF"))
)

func Parse(tokens iter.Seq[token.Token]) (Node, error) {
	program := NewProgram()
	next, stop := iter.Pull(tokens)
	defer stop()

	for {
		t, ok := next()
		if !ok {
			break
		}

		switch t.Type() {
		case token.TypeKeywordDefine:
			statement, err := parseDefine(next)
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

func parseDefine(next NextToken) (*Statement, error) {
	nameIdentifier, err := expectToken[*token.Identifier](next, "`def` should be followed by name")
	if err != nil {
		return nil, err
	}

	_, err = expectToken[*token.OperatorAssign](next, "`def` should be followed by assign")
	if err != nil {
		return nil, err
	}

	expression, err := parseExpression(next)
	if err != nil {
		return nil, err
	}

	return NewStatement(
		NewDefine(nameIdentifier.Name, expression),
	), nil
}

func expectToken[T token.Token](next NextToken, explain string) (T, error) {
	t, ok := next()
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

func parseExpression(next NextToken) (*Expression, error) {
	literalInt, err := expectToken[*token.LiteralInt](next, "should be expression")
	if err != nil {
		return nil, err
	}

	return NewExpression(
		NewLiteralInt(literalInt.Value),
	), nil
}
