package ast

import (
	"encoding/json"
)

type ExpressionPlus struct {
	Value DataNode
}

func NewExpressionPlus(value DataNode) *ExpressionPlus {
	return &ExpressionPlus{Value: value}
}

var _ ExpressionNode = (*ExpressionPlus)(nil)
var _ json.Marshaler = (*ExpressionPlus)(nil)

func (s *ExpressionPlus) Type() Type {
	return TypeExpressionPlus
}

func (s *ExpressionPlus) DataType() DataType {
	return s.Value.DataType()
}

func (s *ExpressionPlus) OperationValue() DataNode {
	return s.Value
}

func (s *ExpressionPlus) SetOperationValue(value DataNode) {
	s.Value = value
}

func (s *ExpressionPlus) OptimizeLiteralOperation(base DataLiteralNode) DataNode {
	var right float64
	if node, ok := s.Value.(DataLiteralNode); ok {
		right = normalizeNumberLiteral(node.LiteralValue())
	} else {
		return s.Value
	}

	left := normalizeNumberLiteral(base.LiteralValue())

	switch {
	case base.DataType() != s.Value.DataType():
		return NewLiteralFloat(left + right)

	case base.DataType() == DataInt:
		base.SetLiteralValue(int(left + right))
		return base

	case base.DataType() == DataFloat:
		base.SetLiteralValue(left + right)
		return base

	default:
		panic(UnreachableErr)
	}
}

func (s *ExpressionPlus) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":  s.Type(),
		"Value": s.Value,
	})
}
