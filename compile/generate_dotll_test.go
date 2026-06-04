package compile

import (
	"strings"
	"testing"

	"arrow_lang/ast"
	"arrow_lang/config"
	"arrow_lang/testutil"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var multiline = cmpopts.AcyclicTransformer("multiline", func(s string) []string {
	return strings.Split(s, "\n")
})

func TestGenerateDotLL(t *testing.T) {
	type testCase struct {
		name     string
		program  *ast.Program
		expected string
	}

	const commonLL = `
		; ModuleID = 'test.arr'
		source_filename = "test.arr"
		target datalayout = "e-m:o-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-n32:64-S128-Fn32"
		target triple = "arm64-apple-darwin25.5.0"

		declare ptr @malloc(i64)

		declare void @free(ptr)
	`

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

			expected: commonLL + `
				define i32 @main() {
				entry:
				  %a = call ptr @malloc(i64 8)
				  store i64 1, ptr %a, align 8
				  ret i32 0
				}
				`,
		},

		{
			name: "define variable with int and free",

			program: ast.NewProgram(
				ast.NewStatement(
					ast.NewDefine("a",
						ast.NewExpression(ast.NewLiteralInt(1)),
					),
				),
				ast.NewStatement(ast.NewFree("a")),
			),

			expected: commonLL + `
				define i32 @main() {
				entry:
				  %a = call ptr @malloc(i64 8)
				  store i64 1, ptr %a, align 8
				  call void @free(ptr %a)
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
			expected := testutil.Dedent(tc.expected)

			result, err := generateDotLL(&Compilation{
				program:       tc.program,
				config:        baseCompilation.config,
				targetMachine: baseCompilation.targetMachine,
				targetTriple:  baseCompilation.targetTriple,
				targetData:    baseCompilation.targetData,
			})

			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(expected, result.String(), multiline); diff != "" {
				t.Error(diff)
			}
		})
	}

	t.Cleanup(baseCompilation.Dispose)
}
