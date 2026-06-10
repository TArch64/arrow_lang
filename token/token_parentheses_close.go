package token

type ParenthesesClose struct{}

func NewParenthesesClose() *ParenthesesClose {
	return &ParenthesesClose{}
}

var _ Token = (*ParenthesesClose)(nil)

func (*ParenthesesClose) Type() Type {
	return TypeParenthesesClose
}

func (*ParenthesesClose) String() string {
	return "BracketClose(Parentheses)"
}
