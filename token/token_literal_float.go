package token

import (
	"fmt"
)

type LiteralFloat struct {
	Value float64
}

func NewLiteralFloat(value float64) *LiteralFloat {
	return &LiteralFloat{Value: value}
}

var _ Token = (*LiteralFloat)(nil)

func (*LiteralFloat) Type() Type {
	return TypeLiteralFloat
}

func (i *LiteralFloat) String() string {
	return fmt.Sprintf("Float(%f)", i.Value)
}
