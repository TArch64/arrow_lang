package ast

import (
	"encoding/json"
)

type Expression struct {
	Content []DataNode
}

func NewExpression(content ...DataNode) *Expression {
	return &Expression{Content: content}
}

var _ DataNode = (*Expression)(nil)
var _ json.Marshaler = (*Expression)(nil)

func (*Expression) Type() Type {
	return TypeExpression
}

func (e *Expression) DataType() DataType {
	return e.Content[0].DataType()
}

func (e *Expression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":    "Expression",
		"Content": e.Content,
	})
}
