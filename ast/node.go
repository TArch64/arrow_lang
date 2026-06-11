package ast

type Type string

const (
	TypeProgram           Type = "Program"
	TypeStatement         Type = "Statement"
	TypeVariable          Type = "Variable"
	TypeFree              Type = "Free"
	TypeDefer             Type = "Defer"
	TypeFunction          Type = "Function"
	TypeFunctionReturn    Type = "FunctionReturn"
	TypeFunctionCall      Type = "FunctionCall"
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
