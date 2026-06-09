package token

type KeywordReturn struct{}

func NewKeywordReturn() *KeywordReturn {
	return &KeywordReturn{}
}

var _ Token = (*KeywordReturn)(nil)

func (*KeywordReturn) Type() Type {
	return TypeKeywordReturn
}

func (*KeywordReturn) String() string {
	return "Keyword(ret)"
}
