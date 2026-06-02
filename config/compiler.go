package config

import (
	"context"
)

type Compiler struct {
	Output string
	Debug  bool
	Ctx    context.Context
}
