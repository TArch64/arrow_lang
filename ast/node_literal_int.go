package ast

import (
	"encoding/json"
)

type LiteralInt struct {
	Value int
}

func NewLiteralInt(value int) *LiteralInt {
	return &LiteralInt{Value: value}
}

var _ DataLiteralNode = (*LiteralInt)(nil)
var _ json.Marshaler = (*LiteralInt)(nil)

func (*LiteralInt) Type() Type {
	return TypeLiteralInt
}

func (i *LiteralInt) DataType() DataType {
	return DataInt
}

func (f *LiteralInt) LiteralValue() any {
	return f.Value
}

func (f *LiteralInt) SetLiteralValue(value any) {
	f.Value = value.(int)
}

func (i *LiteralInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":  i.Type(),
		"Value": i.Value,
	})
}
