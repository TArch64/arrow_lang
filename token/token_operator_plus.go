package token

type OperatorPlus struct{}

func NewOperatorPlus() *OperatorPlus {
	return &OperatorPlus{}
}

var _ Token = (*OperatorPlus)(nil)

func (*OperatorPlus) Type() Type {
	return TypeOperatorPlus
}

func (*OperatorPlus) String() string {
	return "Operator(+)"
}
