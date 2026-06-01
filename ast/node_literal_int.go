package ast

type LiteralInt struct {
	Value int
}

func NewLiteralInt(value int) *LiteralInt {
	return &LiteralInt{Value: value}
}

var _ Node = (*LiteralInt)(nil)

func (*LiteralInt) Type() Type {
	return TypeLiteralInt
}
