package ast

type Type string

const (
	TypeProgram           Type = "Program"
	TypeStatement         Type = "Statement"
	TypeDefine            Type = "Define"
	TypeFree              Type = "Free"
	TypeExpression        Type = "Expression"
	TypeExpressionSum     Type = "ExpressionSum"
	TypeLiteralInt        Type = "LiteralInt"
	TypeLiteralFloat      Type = "LiteralFloat"
	TypeVariableReference Type = "VariableReference"
)

type Node interface {
	Type() Type
}
