package token

type OperatorMinus struct{}

func NewOperatorMinus() *OperatorMinus {
	return &OperatorMinus{}
}

var _ Token = (*OperatorMinus)(nil)

func (*OperatorMinus) Type() Type {
	return TypeOperatorMinus
}

func (*OperatorMinus) String() string {
	return "Operator(+)"
}
