package testutil

import (
	"strings"
)

func Dedent(text string) string {
	text = strings.TrimLeft(text, "\n")
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, "\t")
	}
	return strings.Join(lines, "\n")
}
