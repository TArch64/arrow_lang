package ast

import (
	"slices"
)

type ParsingExpressionAdd func(node DataNode)

type ParsingExpression struct {
	operations []ExpressionNode
}

func NewParsingExpression() *ParsingExpression {
	return &ParsingExpression{}
}

func (e *ParsingExpression) Open(node DataNode) {
	e.operations = []ExpressionNode{NewExpressionOpen(node)}
}

func (e *ParsingExpression) Plus(node DataNode) {
	e.operations = append(e.operations, NewExpressionPlus(node))
}

func (e *ParsingExpression) Minus(node DataNode) {
	e.operations = append(e.operations, NewExpressionMinus(node))
}

func (e *ParsingExpression) Build() *Expression {
	return NewExpression(e.optimize())
}

func (e *ParsingExpression) optimize() []DataNode {
	if len(e.operations) == 1 {
		return []DataNode{e.operations[0].OperationValue()}
	}

	result := make([]DataNode, 0, len(e.operations))
	var last ExpressionNode

	for _, operation := range e.operations {
		if last == nil {
			result = append(result, operation)
			last = operation
			continue
		}

		if _, ok := operation.OperationValue().(DataLiteralNode); !ok {
			result = append(result, operation)
			last = operation
			continue
		}

		lastLiteral, ok := last.OperationValue().(DataLiteralNode)
		if !ok {
			result = append(result, operation)
			last = operation
			continue
		}

		lastValue := last.OperationValue()
		lastValue = operation.OptimizeLiteralOperation(lastLiteral)
		last.SetOperationValue(lastValue)
	}

	result[0] = result[0].(ExpressionNode).OperationValue()

	if len(result) == 1 {
		return []DataNode{result[0]}
	}

	return slices.Clip(result)
}
