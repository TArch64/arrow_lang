package ast

type DataType string

const (
	DataInt   DataType = "Int"
	DataFloat DataType = "Float"
	DataVoid  DataType = "Void"
)

type DataNode interface {
	Node
	DataType() DataType
}

type DataLiteralNode interface {
	DataNode
	LiteralValue() any
	SetLiteralValue(value any)
}

func normalizeNumberLiteral(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	default:
		panic("not a number")
	}
}
