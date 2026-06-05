package ast

import (
	"fmt"
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

func (c *ParsingCtx) ExpectDefined(identifier *token.Identifier) (*Define, error) {
	define, ok := c.defined[identifier.Name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", UndefinedVariableErr, identifier.Name)
	}
	return define, nil
}
