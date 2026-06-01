package ast

type Expression struct {
	Content []Node
}

func NewExpression(content ...Node) *Expression {
	return &Expression{Content: content}
}

var _ Node = (*Expression)(nil)

func (*Expression) Type() Type {
	return TypeExpression
}
