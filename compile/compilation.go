package compile

import (
	"arrow_lang/ast"
	"arrow_lang/config"

	"tinygo.org/x/go-llvm"
)

type Debug struct {
	Filename string
	Dir      string
}

type Compilation struct {
	config        *config.Compiler
	program       *ast.Program
	targetMachine llvm.TargetMachine
	targetTriple  string
	targetData    llvm.TargetData
}

func (c *Compilation) Dispose() {
	c.targetMachine.Dispose()
	c.targetData.Dispose()
}
