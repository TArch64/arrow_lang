package token

type OperatorAssign struct{}

func NewOperatorAssign() *OperatorAssign {
	return &OperatorAssign{}
}

var _ Token = (*OperatorAssign)(nil)

func (*OperatorAssign) Type() Type {
	return TypeOperatorAssign
}

func (*OperatorAssign) String() string {
	return "Operator(=)"
}
