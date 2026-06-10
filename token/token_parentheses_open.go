package token

type ParenthesesOpen struct{}

func NewParenthesesOpen() *ParenthesesOpen {
	return &ParenthesesOpen{}
}

var _ Token = (*ParenthesesOpen)(nil)

func (*ParenthesesOpen) Type() Type {
	return TypeParenthesesOpen
}

func (*ParenthesesOpen) String() string {
	return "BracketOpen(Parentheses)"
}
