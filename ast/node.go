package ast

type Type uint8

const (
	TypeProgram Type = iota
	TypeStatement
	TypeDefine
	TypeFree
	TypeExpression
	TypeLiteralInt
	TypeLiteralFloat
)

type Node interface {
	Type() Type
}
