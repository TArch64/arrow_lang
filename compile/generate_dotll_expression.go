package compile

import (
	"arrow_lang/ast"
	"arrow_lang/errext"

	"tinygo.org/x/go-llvm"
)

func (g *Generation) generateExpression(expression *ast.Expression) llvm.Value {
	defs := map[string]llvm.Value{}
	acc := g.generateExpressionValue(defs, expression.Content[0])

	for _, node := range expression.Content[1:] {
		acc = g.generateExpressionOperation(acc, defs, node)
	}

	return acc
}

func (g *Generation) generateExpressionValue(defs map[string]llvm.Value, node ast.DataNode) llvm.Value {
	switch node := node.(type) {
	case *ast.LiteralInt:
		return llvm.ConstInt(g.std.i64T, uint64(node.Value), node.Value < 0)

	case *ast.LiteralFloat:
		return llvm.ConstFloat(g.std.doubleT, node.Value)

	case *ast.VariableReference:
		defName := node.Reference.Name
		if cached, ok := defs[defName]; ok {
			return cached
		}

		valueType := g.astToType(node.DataType())
		valueName := g.names.WithPrefix(defName + "_v")
		return g.builder.CreateLoad(valueType, g.defined[defName], valueName)

	default:
		panic(errext.Tag("expression", UnreachableErr))
	}
}

func (g *Generation) generateExpressionOperation(acc llvm.Value, defs map[string]llvm.Value, node ast.DataNode) llvm.Value {
	switch node := node.(type) {
	case *ast.ExpressionPlus:
		value := g.generateExpressionValue(defs, node.Value)
		return g.builder.CreateAdd(acc, value, g.names.Random())

	case *ast.ExpressionMinus:
		value := g.generateExpressionValue(defs, node.Value)
		return g.builder.CreateSub(acc, value, g.names.Random())

	default:
		panic(errext.Tag("expression", UnreachableErr))
	}
}
