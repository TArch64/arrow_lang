package ast

type Type uint8

const (
	TypeProgram Type = iota
	TypeStatement
	TypeDefine
	TypeFree
	TypeExpression
	TypeLiteralInt
)

type Node interface {
	Type() Type
}
