package dotll

import (
	"fmt"
	"strings"
)

const (
	CallAlloca = "alloca"
	CallStore  = "store"
	CallRet    = "ret"
)

type CallArg struct {
	Type  DataType
	Value string
}

type BuilderCall struct {
	name string
	args []*CallArg
}

func NewCall(name string) *BuilderCall {
	return &BuilderCall{name: name}
}

func (c *BuilderCall) ArgInt32(value int32) *BuilderCall {
	c.args = append(c.args, &CallArg{
		Type:  DataInt32,
		Value: fmt.Sprintf("%d", value),
	})

	return c
}

func (c *BuilderCall) ArgInt64(value int64) *BuilderCall {
	c.args = append(c.args, &CallArg{
		Type:  DataInt64,
		Value: fmt.Sprintf("%d", value),
	})

	return c
}

func (c *BuilderCall) ArgPtr(value string) *BuilderCall {
	c.args = append(c.args, &CallArg{
		Type:  DataPtr,
		Value: fmt.Sprintf("%%%s", value),
	})

	return c
}

func (c *BuilderCall) ArgType(value DataType) *BuilderCall {
	c.args = append(c.args, &CallArg{Type: value})
	return c
}

func (c *BuilderCall) Render(builder *strings.Builder) {
	builder.WriteString(c.name)

	firstArg := true
	for _, arg := range c.args {
		if !firstArg {
			builder.WriteRune(',')
		}
		builder.WriteRune(' ')

		builder.WriteString(string(arg.Type))

		if arg.Value != "" {
			builder.WriteRune(' ')
			builder.WriteString(arg.Value)
		}

		firstArg = false
	}
}
