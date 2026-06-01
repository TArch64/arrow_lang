package ast

type Program struct {
	Content []*Statement
}

func NewProgram(content ...*Statement) *Program {
	return &Program{Content: content}
}

var _ Node = (*Program)(nil)

func (*Program) Type() Type {
	return TypeProgram
}

func (p *Program) AddStatement(statement *Statement) {
	p.Content = append(p.Content, statement)
}
