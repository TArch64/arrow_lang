package compile

import (
	"arrow_lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *Generation) generateExpression(expression *ast.Expression) llvm.Value {
	defValues := map[string]llvm.Value{}
	return g.generateExpressionNode(defValues, expression.Content)
}

func (g *Generation) generateExpressionNode(defValues map[string]llvm.Value, node ast.DataNode) llvm.Value {
	switch node := node.(type) {
	case *ast.LiteralInt:
		return llvm.ConstInt(g.std.i64T, uint64(node.Value), node.Value < 0)

	case *ast.LiteralFloat:
		return llvm.ConstFloat(g.std.doubleT, node.Value)

	case *ast.VariableReference:
		defName := node.Reference.Name
		if cached, ok := defValues[defName]; ok {
			return cached
		}

		valueType := g.astToType(node.DataType())
		valueName := g.names.WithPrefix(defName + "_v")
		return g.builder.CreateLoad(valueType, g.defined[defName], valueName)

	case *ast.ExpressionSum:
		current := g.generateExpressionNode(defValues, node.Content[0])
		for _, node := range node.Content[1:] {
			value := g.generateExpressionNode(defValues, node)
			current = g.builder.CreateAdd(current, value, g.names.Random())
		}
		return current

	default:
		panic("unreachable")
	}
}
