package ast

type Type string

const (
	TypeProgram           Type = "Program"
	TypeStatement         Type = "Statement"
	TypeVariable          Type = "Variable"
	TypeFunction          Type = "Function"
	TypeFree              Type = "Free"
	TypeReturn            Type = "Return"
	TypeExpression        Type = "Expression"
	TypeExpressionOpen    Type = "ExpressionOpen"
	TypeExpressionPlus    Type = "ExpressionPlus"
	TypeExpressionMinus   Type = "ExpressionMinus"
	TypeLiteralInt        Type = "LiteralInt"
	TypeLiteralFloat      Type = "LiteralFloat"
	TypeVariableReference Type = "VariableReference"
)

type Node interface {
	Type() Type
}
