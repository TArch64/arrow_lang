package ast

import (
	"encoding/json"
)

type Defer struct {
	Statement *Statement
}

func NewDefer(statement *Statement) *Defer {
	return &Defer{Statement: statement}
}

var _ Node = (*Defer)(nil)
var _ json.Marshaler = (*Defer)(nil)

func (*Defer) Type() Type {
	return TypeDefer
}

func (f *Defer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":      f.Type(),
		"Statement": f.Statement,
	})
}
