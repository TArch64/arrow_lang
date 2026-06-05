package compile

import (
	"arrow_lang/ast"

	"tinygo.org/x/go-llvm"
)

func generateDotLL(compilation *Compilation) (llvm.Module, error) {
	generation := compilation.newGeneration()

	mainFn := llvm.AddFunction(generation.mod, "main",
		llvm.FunctionType(generation.std.i32T, nil, false),
	)

	entryBlock := generation.ctx.AddBasicBlock(mainFn, "entry")
	generation.builder.SetInsertPointAtEnd(entryBlock)

	for _, statement := range compilation.program.Statements {
		generateStatement(generation, statement)
	}

	generation.builder.CreateRet(
		llvm.ConstInt(generation.std.i32T, 0, false),
	)

	err := llvm.VerifyModule(generation.mod, llvm.PrintMessageAction)
	return generation.mod, err
}

func generateStatement(generation *Generation, statement *ast.Statement) {
	switch statement := statement.Content.(type) {
	case *ast.Define:
		generateDefine(generation, statement)
	case *ast.Free:
		generateFree(generation, statement)

	default:
		panic("unknown statement type")
	}
}

func generateDefine(generation *Generation, define *ast.Define) {
	defType := generation.astToType(define.DataType())

	def := generation.builder.CreateCall(
		generation.std.mallocT,
		generation.std.malloc,
		[]llvm.Value{generation.std.sizeOf(defType)},
		define.Name,
	)

	value := generateDefineValue(generation, defType, define.Expression)
	generation.builder.CreateStore(value, def)
	generation.defined[define.Name] = def
}

func generateDefineValue(generation *Generation, def llvm.Type, expression *ast.Expression) llvm.Value {
	switch expression := expression.Content[0].(type) {
	case *ast.LiteralInt:
		return llvm.ConstInt(def, uint64(expression.Value), expression.Value < 0)

	case *ast.LiteralFloat:
		return llvm.ConstFloat(def, expression.Value)

	case *ast.VariableReference:
		return generation.defined[expression.Reference.Name]

	default:
		panic("unknown expression type")
	}
}

func generateFree(generation *Generation, free *ast.Free) {
	generation.builder.CreateCall(
		generation.std.freeT,
		generation.std.free,
		[]llvm.Value{generation.defined[free.Reference.Name]},
		"",
	)
}
