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
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 1, ptr %a_1, align 8
				  ret i32 0
				}
				`,
		},

		{
			name: "define variable with int and free",

			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralInt(1)))

				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(ast.NewFree(defA)),
				)
			}(),

			expected: commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 1, ptr %a_1, align 8
				  call void @free(ptr %a_1)
				  ret i32 0
				}
				`,
		},

		{
			name: "define variable with float",

			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralFloat(1.123)))

				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(ast.NewFree(defA)),
				)
			}(),

			expected: commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store double 1.123000e+00, ptr %a_1, align 8
				  call void @free(ptr %a_1)
				  ret i32 0
				}
				`,
		},

		{
			name: "define variable with assign to variable",

			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralInt(1)))

				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(
						ast.NewDefine("b",
							ast.NewExpression(ast.NewVariableReference(defA)),
						),
					),
				)
			}(),

			expected: commonLL + `
			define i32 @main() {
			entry:
			  %a_1 = call ptr @malloc(i64 8)
			  store i64 1, ptr %a_1, align 8
			  %b_2 = call ptr @malloc(i64 8)
			  %a_v_3 = load i64, ptr %a_1, align 8
			  store i64 %a_v_3, ptr %b_2, align 8
			  ret i32 0
			}
			`,
		},

		{
			name: "define variable with sum of variable and literal",

			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralInt(1)))

				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(
						ast.NewDefine("b",
							ast.NewExpression(
								ast.NewExpressionSum(
									ast.NewVariableReference(defA),
									ast.NewLiteralInt(2),
								),
							),
						),
					),
				)
			}(),

			expected: commonLL + `
			define i32 @main() {
			entry:
			  %a_1 = call ptr @malloc(i64 8)
			  store i64 1, ptr %a_1, align 8
			  %b_2 = call ptr @malloc(i64 8)
			  %a_v_3 = load i64, ptr %a_1, align 8
			  %_4 = add i64 %a_v_3, 2
			  store i64 %_4, ptr %b_2, align 8
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
			compilation := &Compilation{
				program:       tc.program,
				config:        baseCompilation.config,
				targetMachine: baseCompilation.targetMachine,
				targetTriple:  baseCompilation.targetTriple,
				targetData:    baseCompilation.targetData,
			}

			result, err := compilation.Generate()
			if err != nil {
				t.Error(err)
			}

			expected := testutil.Dedent(tc.expected)
			if diff := cmp.Diff(expected, result.String(), multiline); diff != "" {
				t.Error(diff)
			}
		})
	}

	t.Cleanup(baseCompilation.Dispose)
}
