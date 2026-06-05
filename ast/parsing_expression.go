package ast

type ParsingExpression struct {
	nodes    []DataNode
	lastNode DataNode
}

func NewParsingExpression() *ParsingExpression {
	return &ParsingExpression{}
}

func (e *ParsingExpression) PlusInt(value int) {
	if e.tryPrecomputePlus(value) {
		return
	}
	e.lastNode = NewLiteralInt(value)
	e.nodes = append(e.nodes, e.lastNode)
}

func (e *ParsingExpression) PlusFloat(value float64) {
	if e.tryPrecomputePlus(value) {
		return
	}
	e.lastNode = NewLiteralFloat(value)
	e.nodes = append(e.nodes, e.lastNode)
}

func (e *ParsingExpression) PlusVariableReference(define *Define) {
	e.lastNode = NewVariableReference(define)
	e.nodes = append(e.nodes, e.lastNode)
}

func (e *ParsingExpression) tryPrecomputePlus(value any) bool {
	if e.lastNode == nil {
		return false
	}

	switch lastNode := e.lastNode.(type) {
	case *LiteralInt:
		if iValue, ok := value.(int); ok {
			lastNode.Value += iValue
		} else {
			e.lastNode = NewLiteralFloat(float64(lastNode.Value) + value.(float64))
			e.nodes[len(e.nodes)-1] = e.lastNode
		}
		return true

	case *LiteralFloat:
		lastNode.Value += value.(float64)
		return true

	default:
		return false
	}
}

func (e *ParsingExpression) Build() *Expression {
	return NewExpression(e.nodes...)
}
