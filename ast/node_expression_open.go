package ast

import (
	"encoding/json"
)

type ExpressionOpen struct {
	Value DataNode
}

func NewExpressionOpen(value DataNode) *ExpressionOpen {
	return &ExpressionOpen{Value: value}
}

var _ ExpressionNode = (*ExpressionOpen)(nil)
var _ json.Marshaler = (*ExpressionOpen)(nil)

func (s *ExpressionOpen) Type() Type {
	return TypeExpressionOpen
}

func (s *ExpressionOpen) DataType() DataType {
	return s.Value.DataType()
}

func (s *ExpressionOpen) OperationValue() DataNode {
	return s.Value
}

func (s *ExpressionOpen) SetOperationValue(value DataNode) {
	s.Value = value
}

func (s *ExpressionOpen) OptimizeLiteralOperation(base DataLiteralNode) DataNode {
	return base
}

func (s *ExpressionOpen) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":  s.Type(),
		"Value": s.Value,
	})
}
