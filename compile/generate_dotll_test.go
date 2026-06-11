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

const commonLL = `
	; ModuleID = 'test.arr'
	source_filename = "test.arr"
	target datalayout = "e-m:o-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-n32:64-S128-Fn32"
	target triple = "arm64-apple-darwin25.5.0"

	declare ptr @malloc(i64)

	declare void @free(ptr)
`

func newTestCompilation(t *testing.T, program *ast.Program) *Compilation {
	t.Helper()

	compilation := &Compilation{
		program: program,
		config: &config.Compiler{
			Output: "/tmp/test.arr",
		},
	}

	if err := initLLVM(compilation); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(compilation.Dispose)

	return compilation
}

func assertGenerate(t *testing.T, compilation *Compilation, expected string) {
	t.Helper()

	result, err := compilation.Generate()
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(testutil.Dedent(expected), result.String(), multiline); diff != "" {
		t.Error(diff)
	}
}

func TestGenerateDefineLiteralInt(t *testing.T) {
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_1 = call ptr @malloc(i64 8)
	  store i64 1, ptr %a_1, align 8
	  ret i32 0
	}
	`)
}

func TestGenerateDefineNegativeInt(t *testing.T) {
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(-42)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_1 = call ptr @malloc(i64 8)
	  store i64 -42, ptr %a_1, align 8
	  ret i32 0
	}
	`)
}

func TestGenerateDefineZeroInt(t *testing.T) {
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewVariable("zero",
				ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(0)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %zero_1 = call ptr @malloc(i64 8)
	  store i64 0, ptr %zero_1, align 8
	  ret i32 0
	}
	`)
}

func TestGenerateDefineFloat(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(1.123)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(ast.NewFree(defA)),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_1 = call ptr @malloc(i64 8)
	  store double 1.123000e+00, ptr %a_1, align 8
	  call void @free(ptr %a_1)
	  ret i32 0
	}
	`)
}

func TestGenerateDefineNegativeFloat(t *testing.T) {
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewVariable("neg",
				ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(-3.14)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %neg_1 = call ptr @malloc(i64 8)
	  store double -3.140000e+00, ptr %neg_1, align 8
	  ret i32 0
	}
	`)
}

func TestGenerateDefineZeroFloat(t *testing.T) {
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewVariable("zero",
				ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(0.0)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %zero_1 = call ptr @malloc(i64 8)
	  store double 0.000000e+00, ptr %zero_1, align 8
	  ret i32 0
	}
	`)
}

func TestGenerateDefineAndFreeInt(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(ast.NewFree(defA)),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_1 = call ptr @malloc(i64 8)
	  store i64 1, ptr %a_1, align 8
	  call void @free(ptr %a_1)
	  ret i32 0
	}
	`)
}

func TestGenerateDefineAndFreeFloat(t *testing.T) {
	defA := ast.NewVariable("pi", ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(3.14159)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(ast.NewFree(defA)),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %pi_1 = call ptr @malloc(i64 8)
	  store double 3.141590e+00, ptr %pi_1, align 8
	  call void @free(ptr %pi_1)
	  ret i32 0
	}
	`)
}

func TestGenerateAssignIntToVariable(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(
			ast.NewVariable("b",
				ast.NewExpression([]ast.DataNode{ast.NewVariableReference(defA)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_1 = call ptr @malloc(i64 8)
	  store i64 1, ptr %a_1, align 8
	  %b_2 = call ptr @malloc(i64 8)
	  %a_v_3 = load i64, ptr %a_1, align 8
	  store i64 %a_v_3, ptr %b_2, align 8
	  ret i32 0
	}
	`)
}

func TestGenerateAssignFloatToVariable(t *testing.T) {
	defA := ast.NewVariable("original", ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(2.718)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(
			ast.NewVariable("copy",
				ast.NewExpression([]ast.DataNode{ast.NewVariableReference(defA)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %original_1 = call ptr @malloc(i64 8)
	  store double 2.718000e+00, ptr %original_1, align 8
	  %copy_2 = call ptr @malloc(i64 8)
	  %original_v_3 = load double, ptr %original_1, align 8
	  store double %original_v_3, ptr %copy_2, align 8
	  ret i32 0
	}
	`)
}

func TestGenerateSumVariableAndLiteral(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(
			ast.NewVariable("b",
				ast.NewExpression([]ast.DataNode{
					ast.NewVariableReference(defA),
					ast.NewExpressionPlus(ast.NewLiteralInt(2)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
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
	`)
}

func TestGenerateSumTwoVariables(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(5)}))
	defY := ast.NewVariable("y", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(10)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defX),
		ast.NewStatement(defY),
		ast.NewStatement(
			ast.NewVariable("sum",
				ast.NewExpression([]ast.DataNode{
					ast.NewVariableReference(defX),
					ast.NewExpressionPlus(ast.NewVariableReference(defY)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
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
	`)
}

func TestGenerateSubtractLiteralFromVariable(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(25)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defX),
		ast.NewStatement(
			ast.NewVariable("result",
				ast.NewExpression([]ast.DataNode{
					ast.NewVariableReference(defX),
					ast.NewExpressionMinus(ast.NewLiteralInt(12)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
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
	`)
}

func TestGenerateSubtractVariableFromLiteral(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(15)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defX),
		ast.NewStatement(
			ast.NewVariable("result",
				ast.NewExpression([]ast.DataNode{
					ast.NewLiteralInt(20),
					ast.NewExpressionMinus(ast.NewVariableReference(defX)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
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
	`)
}

func TestGenerateSubtractTwoVariables(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(30)}))
	defY := ast.NewVariable("y", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(18)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defX),
		ast.NewStatement(defY),
		ast.NewStatement(
			ast.NewVariable("result",
				ast.NewExpression([]ast.DataNode{
					ast.NewVariableReference(defX),
					ast.NewExpressionMinus(ast.NewVariableReference(defY)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
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
	`)
}

func TestGenerateMixedPlusMinusWithVariables(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(50)}))
	defB := ast.NewVariable("b", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(20)}))
	defC := ast.NewVariable("c", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(8)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(defB),
		ast.NewStatement(defC),
		ast.NewStatement(
			ast.NewVariable("result",
				ast.NewExpression([]ast.DataNode{
					ast.NewVariableReference(defA),
					ast.NewExpressionMinus(ast.NewVariableReference(defB)),
					ast.NewExpressionPlus(ast.NewVariableReference(defC)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
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
	`)
}

func TestGenerateMultipleOperations(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(100)}))
	defB := ast.NewVariable("b", ast.NewExpression([]ast.DataNode{ast.NewVariableReference(defA)}))
	defC := ast.NewVariable("c", ast.NewExpression([]ast.DataNode{
		ast.NewVariableReference(defB),
		ast.NewExpressionPlus(ast.NewLiteralInt(50)),
	}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(defA),
		ast.NewStatement(defB),
		ast.NewStatement(defC),
		ast.NewStatement(ast.NewFree(defA)),
		ast.NewStatement(ast.NewFree(defB)),
	}))

	assertGenerate(t, compilation, commonLL+`
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
	`)
}

func TestGenerateDefineBasicGetter(t *testing.T) {
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  ret i64 1
	}
	`)
}

func TestGenerateGetterWithLocalVariable(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(defX),
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{ast.NewVariableReference(defX)}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  %x_2 = call ptr @malloc(i64 8)
	  store i64 1, ptr %x_2, align 8
	  %x_v_3 = load i64, ptr %x_2, align 8
	  ret i64 %x_v_3
	}
	`)
}

func TestGenerateCallBasicGetter(t *testing.T) {
	function := ast.NewFunction("fn", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)})),
		),
	})
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(function),
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{
					ast.NewFunctionCall(function),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_2 = call ptr @malloc(i64 8)
	  %_3 = call i64 @fn_1()
	  store i64 %_3, ptr %a_2, align 8
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  ret i64 1
	}
	`)
}

func TestGenerateCallPlusLiteral(t *testing.T) {
	function := ast.NewFunction("fn", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)})),
		),
	})
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(function),
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{
					ast.NewFunctionCall(function),
					ast.NewExpressionPlus(ast.NewLiteralInt(2)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_2 = call ptr @malloc(i64 8)
	  %_3 = call i64 @fn_1()
	  %_4 = add i64 %_3, 2
	  store i64 %_4, ptr %a_2, align 8
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  ret i64 1
	}
	`)
}

func TestGenerateLiteralPlusCall(t *testing.T) {
	function := ast.NewFunction("fn", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)})),
		),
	})
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(function),
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{
					ast.NewLiteralInt(2),
					ast.NewExpressionPlus(ast.NewFunctionCall(function)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_2 = call ptr @malloc(i64 8)
	  %_3 = call i64 @fn_1()
	  %_4 = add i64 2, %_3
	  store i64 %_4, ptr %a_2, align 8
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  ret i64 1
	}
	`)
}

func TestGenerateCallMinusCall(t *testing.T) {
	one := ast.NewFunction("one", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)})),
		),
	})
	two := ast.NewFunction("two", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(2)})),
		),
	})
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(one),
		ast.NewStatement(two),
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{
					ast.NewFunctionCall(two),
					ast.NewExpressionMinus(ast.NewFunctionCall(one)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_3 = call ptr @malloc(i64 8)
	  %_4 = call i64 @two_2()
	  %_5 = call i64 @one_1()
	  %_6 = sub i64 %_4, %_5
	  store i64 %_6, ptr %a_3, align 8
	  ret i32 0
	}

	define i64 @one_1() {
	entry:
	  ret i64 1
	}

	define i64 @two_2() {
	entry:
	  ret i64 2
	}
	`)
}

func TestGenerateCallPlusVariable(t *testing.T) {
	function := ast.NewFunction("fn", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)})),
		),
	})
	defB := ast.NewVariable("b", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(5)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(function),
		ast.NewStatement(defB),
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{
					ast.NewFunctionCall(function),
					ast.NewExpressionPlus(ast.NewVariableReference(defB)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %b_2 = call ptr @malloc(i64 8)
	  store i64 5, ptr %b_2, align 8
	  %a_3 = call ptr @malloc(i64 8)
	  %_4 = call i64 @fn_1()
	  %b_v_5 = load i64, ptr %b_2, align 8
	  %_6 = add i64 %_4, %b_v_5
	  store i64 %_6, ptr %a_3, align 8
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  ret i64 1
	}
	`)
}

func TestGenerateChainedCallExpression(t *testing.T) {
	function := ast.NewFunction("fn", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)})),
		),
	})
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(function),
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{
					ast.NewFunctionCall(function),
					ast.NewExpressionPlus(ast.NewFunctionCall(function)),
					ast.NewExpressionMinus(ast.NewLiteralInt(2)),
				}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_2 = call ptr @malloc(i64 8)
	  %_3 = call i64 @fn_1()
	  %_4 = call i64 @fn_1()
	  %_5 = add i64 %_3, %_4
	  %_6 = sub i64 %_5, 2
	  store i64 %_6, ptr %a_2, align 8
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  ret i64 1
	}
	`)
}

func TestGenerateDeferFree(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(defX),
				ast.NewStatement(
					ast.NewDefer(
						ast.NewStatement(ast.NewFree(defX)),
					),
				),
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{
							ast.NewVariableReference(defX),
						}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  %x_2 = call ptr @malloc(i64 8)
	  store i64 1, ptr %x_2, align 8
	  %x_v_3 = load i64, ptr %x_2, align 8
	  call void @free(ptr %x_2)
	  ret i64 %x_v_3
	}
	`)
}

func TestGenerateMultipleDefers(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	defB := ast.NewVariable("b", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(2)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(defA),
				ast.NewStatement(defB),
				ast.NewStatement(
					ast.NewDefer(ast.NewStatement(ast.NewFree(defA))),
				),
				ast.NewStatement(
					ast.NewDefer(ast.NewStatement(ast.NewFree(defB))),
				),
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{ast.NewVariableReference(defA)}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  %a_2 = call ptr @malloc(i64 8)
	  store i64 1, ptr %a_2, align 8
	  %b_3 = call ptr @malloc(i64 8)
	  store i64 2, ptr %b_3, align 8
	  %a_v_4 = load i64, ptr %a_2, align 8
	  call void @free(ptr %a_2)
	  call void @free(ptr %b_3)
	  ret i64 %a_v_4
	}
	`)
}

func TestGenerateDeferFreeWithReturnExpression(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(defX),
				ast.NewStatement(
					ast.NewDefer(ast.NewStatement(ast.NewFree(defX))),
				),
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{
							ast.NewVariableReference(defX),
							ast.NewExpressionPlus(ast.NewLiteralInt(2)),
						}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  %x_2 = call ptr @malloc(i64 8)
	  store i64 1, ptr %x_2, align 8
	  %x_v_3 = load i64, ptr %x_2, align 8
	  %_4 = add i64 %x_v_3, 2
	  call void @free(ptr %x_2)
	  ret i64 %_4
	}
	`)
}

func TestGenerateFloatReturningFunction(t *testing.T) {
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(1.5)}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define double @fn_1() {
	entry:
	  ret double 1.500000e+00
	}
	`)
}

func TestGenerateFloatGetterWithLocalVariable(t *testing.T) {
	defX := ast.NewVariable("x", ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(1.5)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(defX),
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{ast.NewVariableReference(defX)}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define double @fn_1() {
	entry:
	  %x_2 = call ptr @malloc(i64 8)
	  store double 1.500000e+00, ptr %x_2, align 8
	  %x_v_3 = load double, ptr %x_2, align 8
	  ret double %x_v_3
	}
	`)
}

func TestGenerateFloatReturningCall(t *testing.T) {
	function := ast.NewFunction("fn", []*ast.Statement{
		ast.NewStatement(
			ast.NewFunctionReturn(ast.NewExpression([]ast.DataNode{ast.NewLiteralFloat(1.5)})),
		),
	})
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(function),
		ast.NewStatement(
			ast.NewVariable("a",
				ast.NewExpression([]ast.DataNode{ast.NewFunctionCall(function)}),
			),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  %a_2 = call ptr @malloc(i64 8)
	  %_3 = call double @fn_1()
	  store double %_3, ptr %a_2, align 8
	  ret i32 0
	}

	define double @fn_1() {
	entry:
	  ret double 1.500000e+00
	}
	`)
}

func TestGenerateFunctionWithFree(t *testing.T) {
	defA := ast.NewVariable("a", ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}))
	compilation := newTestCompilation(t, ast.NewProgram([]*ast.Statement{
		ast.NewStatement(
			ast.NewFunction("fn", []*ast.Statement{
				ast.NewStatement(defA),
				ast.NewStatement(ast.NewFree(defA)),
				ast.NewStatement(
					ast.NewFunctionReturn(
						ast.NewExpression([]ast.DataNode{ast.NewLiteralInt(1)}),
					),
				),
			}),
		),
	}))

	assertGenerate(t, compilation, commonLL+`
	define i32 @main() {
	entry:
	  ret i32 0
	}

	define i64 @fn_1() {
	entry:
	  %a_2 = call ptr @malloc(i64 8)
	  store i64 1, ptr %a_2, align 8
	  call void @free(ptr %a_2)
	  ret i64 1
	}
	`)
}
