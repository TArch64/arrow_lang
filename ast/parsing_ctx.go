package ast

import (
	"iter"

	"arrow_lang/token"
)

type ParsingCtx struct {
	Seq     *ParsingSeq
	defined map[string]*Define
}

func NewParsingCtx(tokens iter.Seq[token.Token]) *ParsingCtx {
	return &ParsingCtx{
		Seq:     NewParsingSeq(tokens),
		defined: make(map[string]*Define),
	}
}

func (c *ParsingCtx) AddDefine(define *Define) {
	c.defined[define.Name] = define
}

func (c *ParsingCtx) IsDefined(name string) bool {
	_, ok := c.defined[name]
	return ok
}
