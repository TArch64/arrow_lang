package ast

type Statement struct {
	Content Node
}

func NewStatement(content Node) *Statement {
	return &Statement{Content: content}
}

var _ Node = (*Statement)(nil)

func (*Statement) Type() Type {
	return TypeStatement
}
