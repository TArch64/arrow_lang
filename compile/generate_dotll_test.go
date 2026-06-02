package compile

import (
	"strings"
	"testing"

	"arrow_lang/ast"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var multiline = cmpopts.AcyclicTransformer("multiline", func(s string) []string {
	return strings.Split(s, "\n")
})

func dedent(text string) string {
	text = strings.TrimLeft(text, "\n")
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, "\t")
	}
	return strings.Join(lines, "\n")
}

func TestGenerateDotLL(t *testing.T) {
	type testCase struct {
		name     string
		program  *ast.Program
		expected string
	}

	testCases := []testCase{
		{
			name: "define variable with literal int",

			program: ast.NewProgram(
				ast.NewStatement(
					ast.NewDefine("a",
						ast.NewExpression(ast.NewLiteralInt(1)),
					),
				),
			),

			expected: `
				define i32 @main() {
				entry:
				  %a = alloca i64
				  store i64 1, ptr %a
				  ret i32 0
				}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := dedent(tc.expected)
			result := generateDotLL(tc.program)
			if diff := cmp.Diff(expected, result, multiline); diff != "" {
				t.Error(diff)
			}
		})
	}
}
