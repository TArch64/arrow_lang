package token

type KeywordFree struct{}

func NewKeywordFree() *KeywordFree {
	return &KeywordFree{}
}

var _ Token = (*KeywordFree)(nil)

func (*KeywordFree) Type() Type {
	return TypeKeywordFree
}

func (*KeywordFree) String() string {
	return "Keyword(free)"
}
