package ast

import (
	"encoding/json"
)

type ExpressionSum struct {
	Content []DataNode
}

func NewExpressionSum(content ...DataNode) *ExpressionSum {
	return &ExpressionSum{Content: content}
}

var _ ExpressionNode = (*ExpressionSum)(nil)
var _ json.Marshaler = (*ExpressionSum)(nil)

func (s *ExpressionSum) Type() Type {
	return TypeExpressionSum
}

func (s *ExpressionSum) DataType() DataType {
	for _, node := range s.Content {
		if node.DataType() == DataFloat {
			return DataFloat
		}
	}

	return DataInt
}

func (s *ExpressionSum) AddExpressionItem(node DataNode) {
	s.Content = append(s.Content, node)
}

func (s *ExpressionSum) OptimizeExpression() DataNode {
	if len(s.Content) == 1 {
		return s.Content[0]
	}
	var result []DataNode
	var constInt *LiteralInt
	var constFloat *LiteralFloat

	for _, node := range s.Content {
		if expressionNode, ok := node.(ExpressionNode); ok {
			node = expressionNode.OptimizeExpression()
		}
		if intNode, ok := node.(*LiteralInt); ok {
			if constInt == nil {
				constInt = intNode
			} else {
				constInt.Value += intNode.Value
			}
		} else if floatNode, ok := node.(*LiteralFloat); ok {
			if constFloat == nil {
				constFloat = floatNode
			} else {
				constFloat.Value += floatNode.Value
			}
		} else {
			result = append(result, node)
		}
	}

	if constFloat != nil && constInt != nil {
		constFloat.Value += float64(constInt.Value)
		result = append(result, constFloat)
	} else if constInt != nil {
		result = append(result, constInt)
	} else if constFloat != nil {
		result = append(result, constFloat)
	}

	if len(result) == 1 {
		return result[0]
	}

	s.Content = result
	return s
}

func (s *ExpressionSum) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":    s.Type(),
		"Content": s.Content,
	})
}
