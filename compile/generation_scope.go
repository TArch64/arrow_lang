package compile

import (
	"fmt"

	"arrow_lang/ast"

	"tinygo.org/x/go-llvm"
)

type DefinedFunction struct {
	Type  llvm.Type
	Value llvm.Value
}

type GenerationScope struct {
	definedVariables map[string]llvm.Value
	definedFunctions map[string]*DefinedFunction
	Deferred         []*ast.Statement
}

func newGenerationScope() *GenerationScope {
	return &GenerationScope{
		definedVariables: make(map[string]llvm.Value),
		definedFunctions: make(map[string]*DefinedFunction),
	}
}

func (g *Generation) newScope() {
	g.scope = newGenerationScope()
	g.scopePath = []*GenerationScope{g.scope}
}

func (g *Generation) diveScope() {
	g.scope = newGenerationScope()
	g.scopePath = append(g.scopePath, g.scope)
}

func (g *Generation) ascendScope() {
	if len(g.scopePath) == 1 {
		panic(fmt.Errorf("%w: trying to ascend top level scope", UnreachableErr))
	}

	g.scopePath = g.scopePath[:len(g.scopePath)-1]
	g.scope = g.scopePath[len(g.scopePath)-1]
}

func (c *GenerationScope) AddVariable(name string, value llvm.Value) {
	c.definedVariables[name] = value
}

func (c *GenerationScope) Variable(name string) llvm.Value {
	return c.definedVariables[name]
}

func (c *GenerationScope) AddFunction(name string, funcType llvm.Type, funcValue llvm.Value) {
	c.definedFunctions[name] = &DefinedFunction{
		Type:  funcType,
		Value: funcValue,
	}
}

func (c *GenerationScope) Function(name string) (llvm.Type, llvm.Value) {
	def := c.definedFunctions[name]
	return def.Type, def.Value
}

func (c *GenerationScope) AddDeferred(statement *ast.Statement) {
	c.Deferred = append(c.Deferred, statement)
}
