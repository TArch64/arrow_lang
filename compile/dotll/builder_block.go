package dotll

import (
	"strings"
)

type BuilderBlock struct {
	name       string
	statements []Builder
}

func NewBlock(name string) *BuilderBlock {
	return &BuilderBlock{name: name}
}

func (b *BuilderBlock) Statement(statement Builder) *BuilderBlock {
	b.statements = append(b.statements, statement)
	return b
}

func (b *BuilderBlock) Render(builder *strings.Builder) {
	builder.WriteString(b.name)
	builder.WriteRune(':')
	newline(builder)

	for _, statement := range b.statements {
		builder.WriteString("  ")
		statement.Render(builder)
		newline(builder)
	}
}
