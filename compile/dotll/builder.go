package dotll

import (
	"strings"
)

func newline(builder *strings.Builder) {
	builder.WriteRune('\n')
}

type Builder interface {
	Render(builder *strings.Builder)
}

func Render(builder Builder) string {
	var strBuilder strings.Builder
	builder.Render(&strBuilder)
	return strBuilder.String()
}
