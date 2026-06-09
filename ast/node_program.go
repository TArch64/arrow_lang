package ast

import (
	"encoding/json"
)

type Program struct {
	Statements []*Statement
}

func NewProgram(statements []*Statement) *Program {
	return &Program{Statements: statements}
}

var _ Node = (*Program)(nil)
var _ json.Marshaler = (*Program)(nil)

func (*Program) Type() Type {
	return TypeProgram
}

func (p *Program) AddStatement(statement *Statement) {
	p.Statements = append(p.Statements, statement)
}

func (p *Program) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":       p.Type(),
		"Statements": p.Statements,
	})
}
