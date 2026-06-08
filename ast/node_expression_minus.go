package ast

import (
	"encoding/json"
)

type ExpressionMinus struct {
	Value DataNode
}

func NewExpressionMinus(value DataNode) *ExpressionMinus {
	return &ExpressionMinus{Value: value}
}

var _ ExpressionNode = (*ExpressionMinus)(nil)
var _ json.Marshaler = (*ExpressionMinus)(nil)

func (s *ExpressionMinus) Type() Type {
	return TypeExpressionMinus
}

func (s *ExpressionMinus) DataType() DataType {
	return s.Value.DataType()
}

func (s *ExpressionMinus) OperationValue() DataNode {
	return s.Value
}

func (s *ExpressionMinus) SetOperationValue(value DataNode) {
	s.Value = value
}

func (s *ExpressionMinus) OptimizeLiteralOperation(base DataLiteralNode) DataNode {
	var right float64
	if node, ok := s.Value.(DataLiteralNode); ok {
		right = normalizeNumberLiteral(node.LiteralValue())
	} else {
		return s.Value
	}

	left := normalizeNumberLiteral(base.LiteralValue())

	switch {
	case base.DataType() != s.Value.DataType():
		return NewLiteralFloat(left - right)

	case base.DataType() == DataInt:
		base.SetLiteralValue(int(left - right))
		return base

	case base.DataType() == DataFloat:
		base.SetLiteralValue(left - right)
		return base

	default:
		panic(UnreachableErr)
	}
}

func (s *ExpressionMinus) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":  s.Type(),
		"Value": s.Value,
	})
}
