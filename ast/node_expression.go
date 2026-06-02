package ast

type Expression struct {
	Content []DataNode
}

func NewExpression(content ...DataNode) *Expression {
	return &Expression{Content: content}
}

var _ DataNode = (*Expression)(nil)

func (*Expression) Type() Type {
	return TypeExpression
}

func (d *Expression) DataType() DataType {
	return d.Content[0].DataType()
}
