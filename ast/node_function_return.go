package ast

import (
	"encoding/json"
)

type FunctionReturn struct {
	Expression *Expression
}

func NewFunctionReturn(expression *Expression) *FunctionReturn {
	return &FunctionReturn{
		Expression: expression,
	}
}

var _ DataNode = (*FunctionReturn)(nil)
var _ json.Marshaler = (*FunctionReturn)(nil)

func (*FunctionReturn) Type() Type {
	return TypeFunctionReturn
}

func (d *FunctionReturn) DataType() DataType {
	return d.Expression.DataType()
}

func (d *FunctionReturn) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":       d.Type(),
		"Expression": d.Expression,
	})
}
