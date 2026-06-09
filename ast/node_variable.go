package ast

import (
	"encoding/json"
)

type Variable struct {
	Name       string
	Expression *Expression
}

func NewVariable(name string, expression *Expression) *Variable {
	return &Variable{
		Name:       name,
		Expression: expression,
	}
}

var _ DataNode = (*Variable)(nil)
var _ json.Marshaler = (*Variable)(nil)

func (*Variable) Type() Type {
	return TypeVariable
}

func (v *Variable) DataType() DataType {
	return v.Expression.DataType()
}

func (v *Variable) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":       v.Type(),
		"Name":       v.Name,
		"Expression": v.Expression,
	})
}
