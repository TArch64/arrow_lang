package compile

import (
	"strings"
	"testing"

	"arrow_lang/ast"
	"arrow_lang/config"

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
				; ModuleID = 'test.arr'
				source_filename = "test.arr"
				target datalayout = "e-m:o-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-n32:64-S128-Fn32"
				target triple = "arm64-apple-darwin25.5.0"

				define i32 @main() {
				entry:
				  %a = alloca i64, align 8
				  store i64 1, ptr %a, align 8
				  ret i32 0
				}
				`,
		},
	}

	baseCompilation := &Compilation{
		config: &config.Compiler{
			Output: "/tmp/test.arr",
		},
	}
	if err := initLLVM(baseCompilation); err != nil {
		t.Error(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := dedent(tc.expected)

			result, err := generateDotLL(&Compilation{
				program:       tc.program,
				config:        baseCompilation.config,
				targetMachine: baseCompilation.targetMachine,
				targetTriple:  baseCompilation.targetTriple,
			})

			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(expected, result.String(), multiline); diff != "" {
				t.Error(diff)
			}
		})
	}
}
