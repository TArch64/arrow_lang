package ast

import (
	"fmt"
	"iter"

	"arrow_lang/token"
)

type ParsingCtx struct {
	Seq       *ParsingSeq
	Scope     *ParsingScope
	scopePath []*ParsingScope
}

func NewParsingCtx(tokens iter.Seq[token.Token]) *ParsingCtx {
	scope := NewParsingScope(nil)

	return &ParsingCtx{
		Seq:       NewParsingSeq(tokens),
		Scope:     scope,
		scopePath: []*ParsingScope{scope},
	}
}

func (c *ParsingCtx) DiveScope() {
	c.Scope = c.Scope.NewChildScope()
	c.scopePath = append(c.scopePath, c.Scope)
}

func (c *ParsingCtx) AscendScope() {
	if len(c.scopePath) == 1 {
		panic(fmt.Errorf("%w: trying to ascend top level scope", UnreachableErr))
	}

	c.scopePath = c.scopePath[:len(c.scopePath)-1]
	c.Scope = c.scopePath[len(c.scopePath)-1]
}
