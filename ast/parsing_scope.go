package ast

import (
	"fmt"

	"arrow_lang/token"
)

type ParsingScope struct {
	parent           *ParsingScope
	definedVariables map[string]*Variable
}

func NewParsingScope() *ParsingScope {
	return &ParsingScope{
		definedVariables: make(map[string]*Variable),
	}
}

func (s *ParsingScope) NewChildScope() *ParsingScope {
	return &ParsingScope{
		parent:           s,
		definedVariables: make(map[string]*Variable),
	}
}

func (s *ParsingScope) AddVariable(variable *Variable) {
	s.definedVariables[variable.Name] = variable
}

func (s *ParsingScope) ExpectVariableDefined(identifier *token.Identifier) (*Variable, error) {
	variable, ok := s.definedVariables[identifier.Name]
	if !ok {
		if s.parent != nil {
			return s.parent.ExpectVariableDefined(identifier)
		}

		return nil, fmt.Errorf("%w: %s", UndefinedVariableErr, identifier.Name)
	}

	return variable, nil
}
