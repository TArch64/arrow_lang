package ast

import (
	"encoding/json"
)

type Return struct {
	Expression *Expression
}

func NewReturn(expression *Expression) *Return {
	return &Return{
		Expression: expression,
	}
}

var _ DataNode = (*Return)(nil)
var _ json.Marshaler = (*Return)(nil)

func (*Return) Type() Type {
	return TypeReturn
}

func (d *Return) DataType() DataType {
	return d.Expression.DataType()
}

func (d *Return) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":       d.Type(),
		"Expression": d.Expression,
	})
}
