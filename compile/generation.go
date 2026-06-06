package compile

import (
	"arrow_lang/ast"

	"tinygo.org/x/go-llvm"
)

type Generation struct {
	ctx        llvm.Context
	mod        llvm.Module
	builder    llvm.Builder
	targetData llvm.TargetData
	std        *GenerationStd
	names      *GenerationNames
	defined    map[string]llvm.Value
}

func (c *Compilation) newGeneration() *Generation {
	ctx := llvm.NewContext()

	mod := ctx.NewModule(c.config.OutputFilename())
	mod.SetDataLayout(c.targetMachine.CreateTargetData().String())
	mod.SetTarget(c.targetTriple)

	generation := &Generation{
		ctx:        ctx,
		mod:        mod,
		builder:    ctx.NewBuilder(),
		targetData: c.targetData,
		defined:    make(map[string]llvm.Value),
	}

	generation.newStd()
	generation.newNames()
	return generation
}

func (g *Generation) astToType(astType ast.DataType) llvm.Type {
	switch astType {
	case ast.DataInt:
		return g.std.i64T

	case ast.DataFloat:
		return g.std.doubleT

	default:
		panic("unknown ast type")
	}
}
