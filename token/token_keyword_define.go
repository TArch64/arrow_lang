package token

type KeywordDefine struct{}

func NewKeywordDefine() *KeywordDefine {
	return &KeywordDefine{}
}

var _ Token = (*KeywordDefine)(nil)

func (*KeywordDefine) Type() Type {
	return TypeKeywordDefine
}

func (*KeywordDefine) String() string {
	return "Keyword(def)"
}
