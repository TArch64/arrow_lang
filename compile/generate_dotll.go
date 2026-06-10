package compile

import (
	"errors"

	"arrow_lang/ast"
	"arrow_lang/errext"

	"tinygo.org/x/go-llvm"
)

var (
	UnexpectedStatementErr = errext.Tag("dotll", errors.New("unexpected statement"))
	UnreachableErr         = errext.Tag("dotll", errors.New("unreachable"))
	UnknownDataTypeErr     = errext.Tag("dotll", errors.New("unknown data type"))
)

func (c *Compilation) Generate() (llvm.Module, error) {
	generation := c.newGeneration()

	mainFn := llvm.AddFunction(generation.mod, "main",
		llvm.FunctionType(generation.std.i32T, nil, false),
	)

	entryBlock := generation.ctx.AddBasicBlock(mainFn, "entry")
	generation.builder.SetInsertPointAtEnd(entryBlock)

	for _, statement := range c.program.Statements {
		generation.generateStatement(statement)
		generation.builder.SetInsertPointAtEnd(entryBlock)
	}

	generation.builder.CreateRet(
		llvm.ConstInt(generation.std.i32T, 0, false),
	)

	err := llvm.VerifyModule(generation.mod, llvm.PrintMessageAction)
	return generation.mod, err
}

func (g *Generation) generateStatement(statement *ast.Statement) {
	switch statement := statement.Content.(type) {
	case *ast.Variable:
		g.generateVariable(statement)

	case *ast.Free:
		g.generateFree(statement)

	case *ast.Function:
		g.generateFunction(statement)

	case *ast.FunctionReturn:
		g.generateFunctionReturn(statement)

	default:
		panic(UnexpectedStatementErr)
	}
}

func (g *Generation) generateVariable(variable *ast.Variable) {
	defType := g.astToType(variable.DataType())

	def := g.builder.CreateCall(
		g.std.mallocT,
		g.std.malloc,
		[]llvm.Value{g.std.sizeOf(defType)},
		g.names.WithPrefix(variable.Name),
	)

	g.builder.CreateStore(g.generateExpression(variable.Expression), def)
	g.definedVariables[variable.Name] = def
}

func (g *Generation) generateFree(free *ast.Free) {
	g.builder.CreateCall(
		g.std.freeT,
		g.std.free,
		[]llvm.Value{g.definedVariables[free.Reference.Name]},
		"",
	)
}

func (g *Generation) generateFunction(function *ast.Function) {
	funcName := g.names.WithPrefix(function.Name)
	funcRetType := g.astToType(function.ReturnDataType())
	funcType := llvm.FunctionType(funcRetType, []llvm.Type{}, false)
	funcValue := llvm.AddFunction(g.mod, funcName, funcType)

	entryBlock := g.ctx.AddBasicBlock(funcValue, "entry")
	g.builder.SetInsertPointAtEnd(entryBlock)

	for _, statement := range function.Statements {
		g.generateStatement(statement)
		g.builder.SetInsertPointAtEnd(entryBlock)
	}

	g.definedFunctions[function.Name] = &DefinedFunction{
		Type:  funcType,
		Value: funcValue,
	}
}

func (g *Generation) generateFunctionReturn(ret *ast.FunctionReturn) {
	g.builder.CreateRet(g.generateExpression(ret.Expression))
}
