package dotll

import (
	"strings"
)

type BuilderDefine struct {
	returnType DataType
	name       string
	blocks     []*BuilderBlock
}

func NewDefine(name string) *BuilderDefine {
	return &BuilderDefine{name: name}
}

func (d *BuilderDefine) Return(returnType DataType) *BuilderDefine {
	d.returnType = returnType
	return d
}

func (d *BuilderDefine) Block(block *BuilderBlock) *BuilderDefine {
	d.blocks = append(d.blocks, block)
	return d
}

func (d *BuilderDefine) Render(builder *strings.Builder) {
	builder.WriteString("define ")
	builder.WriteString(string(d.returnType))
	builder.WriteRune(' ')
	builder.WriteString(d.name)
	builder.WriteString("() {")
	newline(builder)

	for _, block := range d.blocks {
		block.Render(builder)
	}

	builder.WriteRune('}')
}
