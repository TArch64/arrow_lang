package ast

import (
	"encoding/json"
)

type Function struct {
	Name       string
	Statements []*Statement
}

func NewFunction(name string, statements []*Statement) *Function {
	return &Function{
		Name:       name,
		Statements: statements,
	}
}

var _ Node = (*Function)(nil)
var _ json.Marshaler = (*Function)(nil)

func (*Function) Type() Type {
	return TypeFunction
}

func (f *Function) ReturnDataType() DataType {
	if len(f.Statements) == 0 {
		return DataVoid
	}

	last := f.Statements[len(f.Statements)-1]
	if ret, ok := last.Content.(*Return); ok {
		return ret.DataType()
	}

	return DataVoid
}

func (f *Function) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":       f.Type(),
		"Name":       f.Name,
		"Statements": f.Statements,
	})
}
