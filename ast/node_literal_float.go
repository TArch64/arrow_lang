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

var _ DataNode = (*LiteralFloat)(nil)
var _ json.Marshaler = (*LiteralFloat)(nil)

func (*LiteralFloat) Type() Type {
	return TypeLiteralFloat
}

func (f *LiteralFloat) DataType() DataType {
	return DataFloat
}

func (f *LiteralFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":  f.Type(),
		"Value": f.Value,
	})
}
