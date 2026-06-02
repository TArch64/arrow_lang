package dotll

import (
	"strings"
)

type BuilderAssign struct {
	name  string
	value Builder
}

func NewAssign(name string) *BuilderAssign {
	return &BuilderAssign{name: name}
}

func (a *BuilderAssign) To(value Builder) *BuilderAssign {
	a.value = value
	return a
}

func (a *BuilderAssign) Render(builder *strings.Builder) {
	builder.WriteRune('%')
	builder.WriteString(a.name)
	builder.WriteString(" = ")
	a.value.Render(builder)
}
