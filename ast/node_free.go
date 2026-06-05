package ast

import (
	"encoding/json"
)

type Free struct {
	Name string
}

func NewFree(name string) *Free {
	return &Free{Name: name}
}

var _ Node = (*Free)(nil)
var _ json.Marshaler = (*Free)(nil)

func (*Free) Type() Type {
	return TypeFree
}

func (f *Free) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type": "Free",
		"Name": f.Name,
	})
}
