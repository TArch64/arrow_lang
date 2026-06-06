package ast

type DataType string

const (
	DataInt   DataType = "Int"
	DataFloat DataType = "Float"
)

type DataNode interface {
	Node
	DataType() DataType
}
