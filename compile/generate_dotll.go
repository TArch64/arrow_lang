package compile

import (
	"arrow_lang/ast"

	"tinygo.org/x/go-llvm"
)

type Generation struct {
	ctx     llvm.Context
	mod     llvm.Module
	builder llvm.Builder

	i32 llvm.Type
	i64 llvm.Type
}

func (c *Compilation) newGeneration() *Generation {
	ctx := llvm.NewContext()

	mod := ctx.NewModule(c.config.OutputFilename())
	mod.SetDataLayout(c.targetMachine.CreateTargetData().String())
	mod.SetTarget(c.targetTriple)

	return &Generation{
		ctx:     ctx,
		mod:     mod,
		builder: ctx.NewBuilder(),
		i32:     ctx.Int32Type(),
		i64:     ctx.Int64Type(),
	}
}

func (c *Generation) astToType(astType ast.DataType) llvm.Type {
	switch astType {
	case ast.DataInt:
		return c.i64

	default:
		panic("unknown ast type")
	}
}

func generateDotLL(compilation *Compilation) (llvm.Module, error) {
	generation := compilation.newGeneration()

	mainFn := llvm.AddFunction(generation.mod, "main",
		llvm.FunctionType(generation.i32, nil, false),
	)

	entryBlock := generation.ctx.AddBasicBlock(mainFn, "entry")
	generation.builder.SetInsertPointAtEnd(entryBlock)

	for _, statement := range compilation.program.Statements {
		generateStatement(generation, statement)
	}

	generation.builder.CreateRet(
		llvm.ConstInt(generation.i32, 0, false),
	)

	err := llvm.VerifyModule(generation.mod, llvm.PrintMessageAction)
	return generation.mod, err
}

func generateStatement(generation *Generation, statement *ast.Statement) {
	switch statement := statement.Content.(type) {
	case *ast.Define:
		generateDefine(generation, statement)

	default:
		panic("unknown statement type")
	}
}

func generateDefine(generation *Generation, define *ast.Define) {
	defType := generation.astToType(define.DataType())
	alloca := generation.builder.CreateAlloca(defType, define.Name)

	literalInt := define.Expression.Content[0].(*ast.LiteralInt)
	value := llvm.ConstInt(defType, uint64(literalInt.Value), literalInt.Value < 0)

	generation.builder.CreateStore(value, alloca)
}
