package token

type OperatorAssign struct{}

var _ Token = (*OperatorAssign)(nil)

func (*OperatorAssign) Type() Type {
	return TypeOperatorAssign
}

func (a *OperatorAssign) String() string {
	return "Operator(=)"
}
