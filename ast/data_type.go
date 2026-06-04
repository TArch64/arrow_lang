package ast

type DataType uint8

const (
	DataInt DataType = iota
	DataFloat
)

type DataNode interface {
	Node
	DataType() DataType
}
