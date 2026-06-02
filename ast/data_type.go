package ast

type DataType uint8

const (
	DataInt DataType = iota
)

type DataNode interface {
	Node
	DataType() DataType
}
