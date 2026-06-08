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
		program  func() *ast.Program
		expected func() string
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
			name: "basic/define_literal_int",
			program: func() *ast.Program {
				return ast.NewProgram(
					ast.NewStatement(
						ast.NewDefine("a",
							ast.NewExpression(ast.NewLiteralInt(1)),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 1, ptr %a_1, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "basic/define_negative_int",
			program: func() *ast.Program {
				return ast.NewProgram(
					ast.NewStatement(
						ast.NewDefine("a",
							ast.NewExpression(ast.NewLiteralInt(-42)),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 -42, ptr %a_1, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "basic/define_zero_int",
			program: func() *ast.Program {
				return ast.NewProgram(
					ast.NewStatement(
						ast.NewDefine("zero",
							ast.NewExpression(ast.NewLiteralInt(0)),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %zero_1 = call ptr @malloc(i64 8)
				  store i64 0, ptr %zero_1, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "basic/define_float",
			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralFloat(1.123)))
				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(ast.NewFree(defA)),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store double 1.123000e+00, ptr %a_1, align 8
				  call void @free(ptr %a_1)
				  ret i32 0
				}
				`
			},
		},
		{
			name: "basic/define_negative_float",
			program: func() *ast.Program {
				return ast.NewProgram(
					ast.NewStatement(
						ast.NewDefine("neg",
							ast.NewExpression(ast.NewLiteralFloat(-3.14)),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %neg_1 = call ptr @malloc(i64 8)
				  store double -3.140000e+00, ptr %neg_1, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "basic/define_zero_float",
			program: func() *ast.Program {
				return ast.NewProgram(
					ast.NewStatement(
						ast.NewDefine("zero",
							ast.NewExpression(ast.NewLiteralFloat(0.0)),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %zero_1 = call ptr @malloc(i64 8)
				  store double 0.000000e+00, ptr %zero_1, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "memory/define_and_free_int",
			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralInt(1)))
				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(ast.NewFree(defA)),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 1, ptr %a_1, align 8
				  call void @free(ptr %a_1)
				  ret i32 0
				}
				`
			},
		},
		{
			name: "memory/define_and_free_float",
			program: func() *ast.Program {
				defA := ast.NewDefine("pi", ast.NewExpression(ast.NewLiteralFloat(3.14159)))
				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(ast.NewFree(defA)),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %pi_1 = call ptr @malloc(i64 8)
				  store double 3.141590e+00, ptr %pi_1, align 8
				  call void @free(ptr %pi_1)
				  ret i32 0
				}
				`
			},
		},
		{
			name: "variables/assign_int_to_variable",
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
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 1, ptr %a_1, align 8
				  %b_2 = call ptr @malloc(i64 8)
				  %a_v_3 = load i64, ptr %a_1, align 8
				  store i64 %a_v_3, ptr %b_2, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "variables/assign_float_to_variable",
			program: func() *ast.Program {
				defA := ast.NewDefine("original", ast.NewExpression(ast.NewLiteralFloat(2.718)))
				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(
						ast.NewDefine("copy",
							ast.NewExpression(ast.NewVariableReference(defA)),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %original_1 = call ptr @malloc(i64 8)
				  store double 2.718000e+00, ptr %original_1, align 8
				  %copy_2 = call ptr @malloc(i64 8)
				  %original_v_3 = load double, ptr %original_1, align 8
				  store double %original_v_3, ptr %copy_2, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "expressions/sum_variable_and_literal",
			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralInt(1)))
				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(
						ast.NewDefine("b",
							ast.NewExpression(
								ast.NewVariableReference(defA),
								ast.NewExpressionPlus(ast.NewLiteralInt(2)),
							),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
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
				`
			},
		},
		{
			name: "expressions/sum_two_variables",
			program: func() *ast.Program {
				defX := ast.NewDefine("x", ast.NewExpression(ast.NewLiteralInt(5)))
				defY := ast.NewDefine("y", ast.NewExpression(ast.NewLiteralInt(10)))
				return ast.NewProgram(
					ast.NewStatement(defX),
					ast.NewStatement(defY),
					ast.NewStatement(
						ast.NewDefine("sum",
							ast.NewExpression(
								ast.NewVariableReference(defX),
								ast.NewExpressionPlus(ast.NewVariableReference(defY)),
							),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %x_1 = call ptr @malloc(i64 8)
				  store i64 5, ptr %x_1, align 8
				  %y_2 = call ptr @malloc(i64 8)
				  store i64 10, ptr %y_2, align 8
				  %sum_3 = call ptr @malloc(i64 8)
				  %x_v_4 = load i64, ptr %x_1, align 8
				  %y_v_5 = load i64, ptr %y_2, align 8
				  %_6 = add i64 %x_v_4, %y_v_5
				  store i64 %_6, ptr %sum_3, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "expressions/subtract_literal_from_variable",
			program: func() *ast.Program {
				defX := ast.NewDefine("x", ast.NewExpression(ast.NewLiteralInt(25)))
				return ast.NewProgram(
					ast.NewStatement(defX),
					ast.NewStatement(
						ast.NewDefine("result",
							ast.NewExpression(
								ast.NewVariableReference(defX),
								ast.NewExpressionMinus(ast.NewLiteralInt(12)),
							),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %x_1 = call ptr @malloc(i64 8)
				  store i64 25, ptr %x_1, align 8
				  %result_2 = call ptr @malloc(i64 8)
				  %x_v_3 = load i64, ptr %x_1, align 8
				  %_4 = sub i64 %x_v_3, 12
				  store i64 %_4, ptr %result_2, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "expressions/subtract_variable_from_literal",
			program: func() *ast.Program {
				defX := ast.NewDefine("x", ast.NewExpression(ast.NewLiteralInt(15)))
				return ast.NewProgram(
					ast.NewStatement(defX),
					ast.NewStatement(
						ast.NewDefine("result",
							ast.NewExpression(
								ast.NewLiteralInt(20),
								ast.NewExpressionMinus(ast.NewVariableReference(defX)),
							),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %x_1 = call ptr @malloc(i64 8)
				  store i64 15, ptr %x_1, align 8
				  %result_2 = call ptr @malloc(i64 8)
				  %x_v_3 = load i64, ptr %x_1, align 8
				  %_4 = sub i64 20, %x_v_3
				  store i64 %_4, ptr %result_2, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "expressions/subtract_two_variables",
			program: func() *ast.Program {
				defX := ast.NewDefine("x", ast.NewExpression(ast.NewLiteralInt(30)))
				defY := ast.NewDefine("y", ast.NewExpression(ast.NewLiteralInt(18)))
				return ast.NewProgram(
					ast.NewStatement(defX),
					ast.NewStatement(defY),
					ast.NewStatement(
						ast.NewDefine("result",
							ast.NewExpression(
								ast.NewVariableReference(defX),
								ast.NewExpressionMinus(ast.NewVariableReference(defY)),
							),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %x_1 = call ptr @malloc(i64 8)
				  store i64 30, ptr %x_1, align 8
				  %y_2 = call ptr @malloc(i64 8)
				  store i64 18, ptr %y_2, align 8
				  %result_3 = call ptr @malloc(i64 8)
				  %x_v_4 = load i64, ptr %x_1, align 8
				  %y_v_5 = load i64, ptr %y_2, align 8
				  %_6 = sub i64 %x_v_4, %y_v_5
				  store i64 %_6, ptr %result_3, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "expressions/mixed_plus_minus_with_variables",
			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralInt(50)))
				defB := ast.NewDefine("b", ast.NewExpression(ast.NewLiteralInt(20)))
				defC := ast.NewDefine("c", ast.NewExpression(ast.NewLiteralInt(8)))
				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(defB),
					ast.NewStatement(defC),
					ast.NewStatement(
						ast.NewDefine("result",
							ast.NewExpression(
								ast.NewVariableReference(defA),
								ast.NewExpressionMinus(ast.NewVariableReference(defB)),
								ast.NewExpressionPlus(ast.NewVariableReference(defC)),
							),
						),
					),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 50, ptr %a_1, align 8
				  %b_2 = call ptr @malloc(i64 8)
				  store i64 20, ptr %b_2, align 8
				  %c_3 = call ptr @malloc(i64 8)
				  store i64 8, ptr %c_3, align 8
				  %result_4 = call ptr @malloc(i64 8)
				  %a_v_5 = load i64, ptr %a_1, align 8
				  %b_v_6 = load i64, ptr %b_2, align 8
				  %_7 = sub i64 %a_v_5, %b_v_6
				  %c_v_8 = load i64, ptr %c_3, align 8
				  %_9 = add i64 %_7, %c_v_8
				  store i64 %_9, ptr %result_4, align 8
				  ret i32 0
				}
				`
			},
		},
		{
			name: "complex/multiple_operations",
			program: func() *ast.Program {
				defA := ast.NewDefine("a", ast.NewExpression(ast.NewLiteralInt(100)))
				defB := ast.NewDefine("b", ast.NewExpression(ast.NewVariableReference(defA)))

				defC := ast.NewDefine("c", ast.NewExpression(
					ast.NewVariableReference(defB),
					ast.NewExpressionPlus(ast.NewLiteralInt(50)),
				))

				return ast.NewProgram(
					ast.NewStatement(defA),
					ast.NewStatement(defB),
					ast.NewStatement(defC),
					ast.NewStatement(ast.NewFree(defA)),
					ast.NewStatement(ast.NewFree(defB)),
				)
			},
			expected: func() string {
				return commonLL + `
				define i32 @main() {
				entry:
				  %a_1 = call ptr @malloc(i64 8)
				  store i64 100, ptr %a_1, align 8
				  %b_2 = call ptr @malloc(i64 8)
				  %a_v_3 = load i64, ptr %a_1, align 8
				  store i64 %a_v_3, ptr %b_2, align 8
				  %c_4 = call ptr @malloc(i64 8)
				  %b_v_5 = load i64, ptr %b_2, align 8
				  %_6 = add i64 %b_v_5, 50
				  store i64 %_6, ptr %c_4, align 8
				  call void @free(ptr %a_1)
				  call void @free(ptr %b_2)
				  ret i32 0
				}
				`
			},
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
				program:       tc.program(),
				config:        baseCompilation.config,
				targetMachine: baseCompilation.targetMachine,
				targetTriple:  baseCompilation.targetTriple,
				targetData:    baseCompilation.targetData,
			}

			result, err := compilation.Generate()
			if err != nil {
				t.Error(err)
			}

			expected := testutil.Dedent(tc.expected())
			if diff := cmp.Diff(expected, result.String(), multiline); diff != "" {
				t.Error(diff)
			}
		})
	}

	t.Cleanup(baseCompilation.Dispose)
}
