package ast

import (
	"errors"
	"fmt"
	"iter"

	"arrow_lang/errext"
	"arrow_lang/token"
)

type NextToken func() (token.Token, bool)

var (
	UnexpectedTokenErr   = errext.Tag("ast", errors.New("unexpected token"))
	UnexpectedEOFErr     = errext.Tag("ast", errors.New("unexpected EOF"))
	UndefinedVariableErr = errext.Tag("ast", errors.New("undefined variable"))
	UndefinedFunctionErr = errext.Tag("ast", errors.New("undefined function"))
	UnreachableErr       = errext.Tag("ast", errext.UnreachableErr)
)

func Parse(tokens iter.Seq[token.Token]) (*Program, error) {
	program := NewProgram([]*Statement{})
	ctx := NewParsingCtx(tokens)
	defer ctx.Seq.Stop()

	for {
		if ctx.Seq.PeekNext() == nil {
			break
		}

		statement, err := parseStatement(ctx)
		if err != nil {
			return nil, err
		}

		program.AddStatement(statement)
	}

	return program, nil
}

func parseStatement(ctx *ParsingCtx) (*Statement, error) {
	expected, err := ctx.Seq.ExpectAny("Invalid statement",
		token.TypeKeywordDefine,
		token.TypeKeywordFree,
		token.TypeKeywordReturn,
	)

	if err != nil {
		return nil, err
	}

	switch expected := expected.(type) {
	case *token.KeywordDefine:
		return parseDefine(ctx)

	case *token.KeywordFree:
		return parseFree(ctx)

	case *token.KeywordReturn:
		return parseFunctionReturn(ctx)

	default:
		return nil, fmt.Errorf("%w: %s", UnexpectedTokenErr, expected.String())
	}
}

func parseDefine(ctx *ParsingCtx) (*Statement, error) {
	nameIdentifier, err := ExpectToken[*token.Identifier](ctx, "`def` should be followed by name")
	if err != nil {
		return nil, err
	}

	expected, err := ctx.Seq.ExpectAny("should be followed by assign or function definition",
		token.TypeOperatorAssign,
		token.TypeParenthesesOpen,
	)
	if err != nil {
		return nil, err
	}

	switch expected.(type) {
	case *token.OperatorAssign:
		return parseVariable(ctx, nameIdentifier)

	case *token.ParenthesesOpen:
		return parseFunction(ctx, nameIdentifier)

	default:
		panic(errext.Tag("define", UnreachableErr))
	}
}

func parseVariable(ctx *ParsingCtx, nameIdentifier *token.Identifier) (*Statement, error) {
	expression, err := parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	variable := NewVariable(nameIdentifier.Name, expression)
	ctx.Scope.AddVariable(variable)

	return NewStatement(variable), nil
}

func parseFunction(ctx *ParsingCtx, nameIdentifier *token.Identifier) (*Statement, error) {
	ctx.Seq.Next() // skip parentheses
	ctx.Seq.Next() // skip curly bracket open

	var statements []*Statement
	ctx.DiveScope()

	for {
		statement, err := parseStatement(ctx)
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)

		if next := ctx.Seq.PeekNext(); next != nil && next.Type() == token.TypeCurlyBracketClose {
			ctx.Seq.Next()
			break
		}
	}

	ctx.AscendScope()

	function := NewFunction(nameIdentifier.Name, statements)
	ctx.Scope.AddFunction(function)

	return NewStatement(function), nil
}

func parseExpression(ctx *ParsingCtx) (*Expression, error) {
	expression := NewParsingExpression()
	var addNode ParsingExpressionAdd = expression.Open

	for {
		expected, err := ctx.Seq.ExpectAny("should be an expression",
			token.TypeLiteralInt,
			token.TypeLiteralFloat,
			token.TypeIdentifier,
		)
		if err != nil {
			return nil, err
		}

		switch expected := expected.(type) {
		case *token.LiteralInt:
			addNode(NewLiteralInt(expected.Value))

		case *token.LiteralFloat:
			addNode(NewLiteralFloat(expected.Value))

		case *token.Identifier:
			if next := ctx.Seq.PeekNext(); next != nil && next.Type() == token.TypeParenthesesOpen {
				function, err := ctx.Scope.ExpectFunctionDefined(expected)
				if err != nil {
					return nil, err
				}

				dataType := function.ReturnDataType()
				if dataType != DataFloat && dataType != DataInt {
					return nil, fmt.Errorf("%w: %s cannot be used in math expressions", UnexpectedTokenErr, dataType)
				}

				call, err := parseFunctionCall(ctx, function)
				if err != nil {
					return nil, err
				}

				addNode(call)
			} else {
				variable, err := ctx.Scope.ExpectVariableDefined(expected)
				if err != nil {
					return nil, err
				}

				dataType := variable.DataType()
				if dataType != DataFloat && dataType != DataInt {
					return nil, fmt.Errorf("%w: %s cannot be used in math expressions", UnexpectedTokenErr, dataType)
				}

				addNode(NewVariableReference(variable))
			}

		default:
			panic(errext.Tag("expression", UnreachableErr))
		}

		switch ctx.Seq.PeekNext().(type) {
		case *token.OperatorPlus:
			ctx.Seq.Next()
			addNode = expression.Plus

		case *token.OperatorMinus:
			ctx.Seq.Next()
			addNode = expression.Minus

		default:
			return expression.Build(), nil
		}
	}

	return expression.Build(), nil
}

func parseFree(ctx *ParsingCtx) (*Statement, error) {
	nameIdentifier, err := ExpectToken[*token.Identifier](ctx, "`free` should be followed by variable name")
	if err != nil {
		return nil, err
	}

	variable, err := ctx.Scope.ExpectVariableDefined(nameIdentifier)
	if err != nil {
		return nil, err
	}

	return NewStatement(NewFree(variable)), nil
}

func parseFunctionReturn(ctx *ParsingCtx) (*Statement, error) {
	expression, err := parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return NewStatement(NewFunctionReturn(expression)), nil
}

func parseFunctionCall(ctx *ParsingCtx, function *Function) (*FunctionCall, error) {
	ctx.Seq.Next() // skip parentheses open
	ctx.Seq.Next() // skip parentheses close
	return NewFunctionCall(function), nil
}
