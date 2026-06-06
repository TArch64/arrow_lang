package config

import (
	"context"
	"path"
)

type Compiler struct {
	Input  string
	Output string
	Debug  bool
	Ctx    context.Context
}

func (c *Compiler) OutputDir() string {
	return path.Dir(c.Output)
}

func (c *Compiler) OutputFilename() string {
	return path.Base(c.Output)
}
