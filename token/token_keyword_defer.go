package token

type KeywordDefer struct{}

func NewKeywordDefer() *KeywordDefer {
	return &KeywordDefer{}
}

var _ Token = (*KeywordDefer)(nil)

func (*KeywordDefer) Type() Type {
	return TypeKeywordDefer
}

func (*KeywordDefer) String() string {
	return "Keyword(defer)"
}
