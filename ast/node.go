package ast

type Type uint8

const (
	TypeProgram Type = iota
	TypeStatement
	TypeDefine
	TypeExpression
	TypeLiteralInt
)

type Node interface {
	Type() Type
}
