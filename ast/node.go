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
	TypeVariableReference
)

type Node interface {
	Type() Type
}
