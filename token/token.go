package token

type Type uint8

const (
	TypeIdentifier Type = iota
	TypeOperatorPlus
	TypeOperatorAssign
	TypeKeywordDefine
	TypeKeywordFree
	TypeLiteralInt
)

type Token interface {
	Type() Type
	String() string
}
