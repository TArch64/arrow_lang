package ast

import (
	"fmt"

	"arrow_lang/token"
)

type ParsingScope struct {
	parent           *ParsingScope
	definedVariables map[string]*Variable
	definedFunctions map[string]*Function
}

func NewParsingScope(parent *ParsingScope) *ParsingScope {
	return &ParsingScope{
		parent:           parent,
		definedVariables: make(map[string]*Variable),
		definedFunctions: make(map[string]*Function),
	}
}

func (s *ParsingScope) NewChildScope() *ParsingScope {
	return NewParsingScope(s)
}

func (s *ParsingScope) AddVariable(variable *Variable) {
	s.definedVariables[variable.Name] = variable
}

func (s *ParsingScope) AddFunction(function *Function) {
	s.definedFunctions[function.Name] = function
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

func (s *ParsingScope) ExpectFunctionDefined(identifier *token.Identifier) (*Function, error) {
	function, ok := s.definedFunctions[identifier.Name]
	if !ok {
		if s.parent != nil {
			return s.parent.ExpectFunctionDefined(identifier)
		}

		return nil, fmt.Errorf("%w: %s", UndefinedFunctionErr, identifier.Name)
	}

	return function, nil
}
