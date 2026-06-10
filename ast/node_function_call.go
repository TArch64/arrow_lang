package ast

import (
	"encoding/json"
)

type FunctionCall struct {
	Function *Function
}

func NewFunctionCall(function *Function) *FunctionCall {
	return &FunctionCall{Function: function}
}

var _ DataNode = (*FunctionCall)(nil)
var _ json.Marshaler = (*FunctionCall)(nil)

func (*FunctionCall) Type() Type {
	return TypeFunctionCall
}

func (d *FunctionCall) DataType() DataType {
	return d.Function.ReturnDataType()
}

func (d *FunctionCall) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":     d.Type(),
		"Function": d.Function.Name,
	})
}
