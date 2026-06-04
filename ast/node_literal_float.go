package ast

type LiteralFloat struct {
	Value float64
}

func NewLiteralFloat(value float64) *LiteralFloat {
	return &LiteralFloat{Value: value}
}

var _ DataNode = (*LiteralFloat)(nil)

func (*LiteralFloat) Type() Type {
	return TypeLiteralFloat
}

func (i *LiteralFloat) DataType() DataType {
	return DataFloat
}
