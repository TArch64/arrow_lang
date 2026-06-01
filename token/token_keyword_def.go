package token

type KeywordDef struct{}

var _ Token = (*KeywordDef)(nil)

func (*KeywordDef) Type() Type {
	return TypeKeywordDef
}

func (d *KeywordDef) String() string {
	return "Keyword(def)"
}
