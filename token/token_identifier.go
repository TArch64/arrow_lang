package token

import (
	"fmt"
)

type Identifier struct {
	Name string
}

func NewIdentifier(name string) *Identifier {
	return &Identifier{Name: name}
}

var _ Token = (*Identifier)(nil)

func (*Identifier) Type() Type {
	return TypeIdentifier
}

func (i *Identifier) String() string {
	return fmt.Sprintf(`Identifier(%s)`, i.Name)
}
