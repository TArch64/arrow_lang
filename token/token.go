package token

type Type uint8

const (
	TypeIdentifier Type = iota
	TypeOperatorAssign
	TypeOperatorPlus
	TypeOperatorMinus
	TypeKeywordDefine
	TypeKeywordFree
	TypeKeywordDefer
	TypeKeywordReturn
	TypeLiteralInt
	TypeLiteralFloat
	TypeCurlyBracketOpen
	TypeCurlyBracketClose
	TypeParenthesesOpen
	TypeParenthesesClose
)

type Token interface {
	Type() Type
	String() string
}
