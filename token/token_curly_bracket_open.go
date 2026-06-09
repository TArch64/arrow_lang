package token

type CurlyBracketOpen struct{}

func NewCurlyBracketOpen() *CurlyBracketOpen {
	return &CurlyBracketOpen{}
}

var _ Token = (*CurlyBracketOpen)(nil)

func (*CurlyBracketOpen) Type() Type {
	return TypeCurlyBracketOpen
}

func (*CurlyBracketOpen) String() string {
	return "BracketOpen(Curly)"
}
