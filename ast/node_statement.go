package ast

import (
	"encoding/json"
)

type Statement struct {
	Content Node
}

func NewStatement(content Node) *Statement {
	return &Statement{Content: content}
}

var _ Node = (*Statement)(nil)
var _ json.Marshaler = (*Statement)(nil)

func (*Statement) Type() Type {
	return TypeStatement
}

func (s *Statement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"Type":    "Statement",
		"Content": s.Content,
	})
}
