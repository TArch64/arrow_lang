package ast

import (
	"encoding/json"
)

type LiteralFloat struct {
	Value float64
}

func NewLiteralFloat(value float64) *LiteralFloat {
	return &LiteralFloat{Value: value}
}

var _ DataLiteralNode = (*LiteralFloat)(nil)
var _ json.Marshaler = (*LiteralFloat)(nil)

func (*LiteralFloat) Type() Type {
	return TypeLiteralFloat
}

func (f *LiteralFloat) DataType() DataType {
	return DataFloat
}

func (f *LiteralFloat) LiteralValue() any {
	return f.Value
}

func (f *LiteralFloat) SetLiteralValue(value any) {
	f.Value = value.(float64)
}

func (f *LiteralFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":  f.Type(),
		"Value": f.Value,
	})
}
