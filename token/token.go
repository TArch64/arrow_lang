package token

type Type uint8

const (
	TypeIdentifier Type = iota
	TypeOperatorPlus
	TypeOperatorAssign
	TypeKeywordDefine
	TypeKeywordFree
	TypeLiteralInt
	TypeLiteralFloat
)

type Token interface {
	Type() Type
	String() string
}
