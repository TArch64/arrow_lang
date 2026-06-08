package token

type Type uint8

const (
	TypeIdentifier Type = iota
	TypeOperatorAssign
	TypeOperatorPlus
	TypeOperatorMinus
	TypeKeywordDefine
	TypeKeywordFree
	TypeLiteralInt
	TypeLiteralFloat
)

type Token interface {
	Type() Type
	String() string
}
