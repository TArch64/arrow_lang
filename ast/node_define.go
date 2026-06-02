package ast

type Define struct {
	Name       string
	Expression *Expression
}

func NewDefine(name string, expression *Expression) *Define {
	return &Define{
		Name:       name,
		Expression: expression,
	}
}

var _ DataNode = (*Define)(nil)

func (*Define) Type() Type {
	return TypeDefine
}

func (d *Define) DataType() DataType {
	return d.Expression.DataType()
}
