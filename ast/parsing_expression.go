package ast

type ParsingExpression struct {
	operationNode *ExpressionSum
}

func NewParsingExpression() *ParsingExpression {
	return &ParsingExpression{}
}

func (e *ParsingExpression) Plus(node DataNode) {
	if e.operationNode == nil {
		e.operationNode = NewExpressionSum(node)
		return
	}
	e.operationNode.AddExpressionItem(node)
}

func (e *ParsingExpression) Build() *Expression {
	return NewExpression(e.operationNode.OptimizeExpression())
}
