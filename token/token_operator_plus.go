package token

type OperatorPlus struct{}

var _ Token = (*OperatorPlus)(nil)

func (*OperatorPlus) Type() Type {
	return TypeOperatorPlus
}

func (p *OperatorPlus) String() string {
	return "Operator(+)"
}
