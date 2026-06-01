package token

import (
	"fmt"
)

type LiteralInt struct {
	Value int
}

var _ Token = (*LiteralInt)(nil)

func (*LiteralInt) Type() Type {
	return TypeLiteralInt
}

func (i *LiteralInt) String() string {
	return fmt.Sprintf("Int(%d)", i.Value)
}
