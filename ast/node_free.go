package ast

import (
	"encoding/json"
)

type Free struct {
	Reference *Variable
}

func NewFree(reference *Variable) *Free {
	return &Free{Reference: reference}
}

var _ Node = (*Free)(nil)
var _ json.Marshaler = (*Free)(nil)

func (*Free) Type() Type {
	return TypeFree
}

func (f *Free) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":      f.Type(),
		"Reference": f.Reference.Name,
	})
}
