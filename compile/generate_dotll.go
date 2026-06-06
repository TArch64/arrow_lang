package compile

import (
	"arrow_lang/ast"

	"tinygo.org/x/go-llvm"
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
	}

	generation.builder.CreateRet(
		llvm.ConstInt(generation.std.i32T, 0, false),
	)

	err := llvm.VerifyModule(generation.mod, llvm.PrintMessageAction)
	return generation.mod, err
}

func (g *Generation) generateStatement(statement *ast.Statement) {
	switch statement := statement.Content.(type) {
	case *ast.Define:
		g.generateDefine(statement)
	case *ast.Free:
		g.generateFree(statement)

	default:
		panic("unknown statement type")
	}
}

func (g *Generation) generateDefine(define *ast.Define) {
	defType := g.astToType(define.DataType())

	def := g.builder.CreateCall(
		g.std.mallocT,
		g.std.malloc,
		[]llvm.Value{g.std.sizeOf(defType)},
		g.names.WithPrefix(define.Name),
	)

	g.builder.CreateStore(g.generateExpression(define.Expression), def)
	g.defined[define.Name] = def
}

func (g *Generation) generateFree(free *ast.Free) {
	g.builder.CreateCall(
		g.std.freeT,
		g.std.free,
		[]llvm.Value{g.defined[free.Reference.Name]},
		"",
	)
}
