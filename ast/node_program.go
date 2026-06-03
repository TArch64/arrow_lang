package ast

type Program struct {
	Statements []*Statement
}

func NewProgram(content ...*Statement) *Program {
	return &Program{Statements: content}
}

var _ Node = (*Program)(nil)

func (*Program) Type() Type {
	return TypeProgram
}

func (p *Program) AddStatement(statement *Statement) {
	p.Statements = append(p.Statements, statement)
}
