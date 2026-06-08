package ast

import (
	"fmt"
	"iter"

	"arrow_lang/token"
)

type ParsingSeq struct {
	next     NextToken
	buffered token.Token
	Stop     func()
}

func NewParsingSeq(tokens iter.Seq[token.Token]) *ParsingSeq {
	next, stop := iter.Pull(tokens)

	return &ParsingSeq{
		next: next,
		Stop: stop,
	}
}

func (s *ParsingSeq) Next() (token.Token, bool) {
	if s.buffered != nil {
		buffered := s.buffered
		s.buffered = nil
		return buffered, true
	}
	return s.next()
}

func (s *ParsingSeq) PeekNext() token.Token {
	if s.buffered != nil {
		return s.buffered
	}

	buffered, ok := s.next()
	if !ok {
		return nil
	}

	s.buffered = buffered
	return buffered
}

func (s *ParsingSeq) ExpectAny(explain string, expectations ...token.Type) (token.Token, error) {
	if len(expectations) == 0 {
		panic("expectAnyToken called with no expectations")
	}

	t, ok := s.Next()
	if !ok {
		return nil, UnexpectedEOFErr
	}

	for _, expectation := range expectations {
		if t.Type() == expectation {
			return t, nil
		}
	}

	return nil, fmt.Errorf("%w: %s, got: %s", UnexpectedTokenErr, explain, t.String())
}

func ExpectToken[T token.Token](ctx *ParsingCtx, explain string) (T, error) {
	var typed T
	expectation, err := ctx.Seq.ExpectAny(explain, typed.Type())
	if err != nil {
		return typed, err
	}

	return expectation.(T), nil
}
