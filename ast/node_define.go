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

var _ Node = (*Define)(nil)

func (*Define) Type() Type {
	return TypeDefine
}
