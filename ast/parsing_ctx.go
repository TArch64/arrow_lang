package ast

import (
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
