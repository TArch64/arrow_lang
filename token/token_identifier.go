package token

import (
	"fmt"
)

type Identifier struct {
	Name string
}

var _ Token = (*Identifier)(nil)

func (*Identifier) Type() Type {
	return TypeIdentifier
}

func (i *Identifier) String() string {
	return fmt.Sprintf(`Identifier(%s)`, i.Name)
}
