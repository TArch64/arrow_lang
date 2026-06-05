package ast

import (
	"encoding/json"
)

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
var _ json.Marshaler = (*Define)(nil)

func (*Define) Type() Type {
	return TypeDefine
}

func (d *Define) DataType() DataType {
	return d.Expression.DataType()
}

func (d *Define) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":       "Define",
		"Name":       d.Name,
		"Expression": d.Expression,
	})
}
