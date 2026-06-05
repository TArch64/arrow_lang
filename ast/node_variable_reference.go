package ast

import (
	"encoding/json"
)

type VariableReference struct {
	Reference *Define
}

func NewVariableReference(reference *Define) *VariableReference {
	return &VariableReference{Reference: reference}
}

var _ DataNode = (*VariableReference)(nil)
var _ json.Marshaler = (*VariableReference)(nil)

func (*VariableReference) Type() Type {
	return TypeVariableReference
}

func (r *VariableReference) DataType() DataType {
	return r.Reference.DataType()
}

func (r *VariableReference) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":      "VariableReference",
		"Reference": r.Reference.Name,
	})
}
