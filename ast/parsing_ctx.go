package ast

import (
	"fmt"
	"iter"

	"arrow_lang/token"
)

type ParsingCtx struct {
	next    NextToken
	stop    func()
	defined map[string]*Define
}

func NewParsingCtx(tokens iter.Seq[token.Token]) *ParsingCtx {
	next, stop := iter.Pull(tokens)

	return &ParsingCtx{
		next:    next,
		stop:    stop,
		defined: make(map[string]*Define),
	}
}

func (c *ParsingCtx) addDefine(define *Define) {
	c.defined[define.Name] = define
}

func (c *ParsingCtx) isDefined(name string) bool {
	_, ok := c.defined[name]
	return ok
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

func expectAnyToken(ctx *ParsingCtx, explain string, expectations ...token.Token) (token.Token, error) {
	if len(expectations) == 0 {
		panic("expectAnyToken called with no expectations")
	}

	t, ok := ctx.next()
	if !ok {
		return nil, UnexpectedEOFErr
	}

	for _, expectation := range expectations {
		if t.Type() == expectation.Type() {
			return t, nil
		}
	}

	return nil, fmt.Errorf("%w: %s, got: %s", UnexpectedTokenErr, explain, t.String())
}
