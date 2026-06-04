package ast

type Free struct {
	Name string
}

func NewFree(name string) *Free {
	return &Free{Name: name}
}

var _ Node = (*Free)(nil)

func (*Free) Type() Type {
	return TypeFree
}
