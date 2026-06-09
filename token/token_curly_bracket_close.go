package token

type CurlyBracketClose struct{}

func NewCurlyBracketClose() *CurlyBracketClose {
	return &CurlyBracketClose{}
}

var _ Token = (*CurlyBracketClose)(nil)

func (*CurlyBracketClose) Type() Type {
	return TypeCurlyBracketClose
}

func (*CurlyBracketClose) String() string {
	return "BracketClose(Curly)"
}
