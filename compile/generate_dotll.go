package compile

import (
	"arrow_lang/ast"
	"arrow_lang/compile/dotll"
)

// 				define i32 @main() {
//				entry:
//				  %x = alloca i64
//				  store i64 42, ptr %x
//				  ret i32 0
//				}

func generateDotLL(program *ast.Program) string {
	main := dotll.NewDefine("@main").Return(dotll.DataInt32)
	entry := dotll.NewBlock("entry")

	for _, statement := range program.Content {
		generated := generateStatement(statement)
		for _, statement := range generated {
			entry.Statement(statement)
		}
	}

	entry.Statement(
		dotll.NewCall(dotll.CallRet).ArgInt32(0),
	)

	main.Block(entry)
	return dotll.Render(main)
}

func generateStatement(statement *ast.Statement) []dotll.Builder {
	switch statement := statement.Content.(type) {
	case *ast.Define:
		return generateDefine(statement)

	default:
		panic("unknown statement type")
	}
}

func generateDefine(define *ast.Define) []dotll.Builder {
	statements := []dotll.Builder{
		dotll.NewAssign(define.Name).To(
			dotll.NewCall(dotll.CallAlloca).
				ArgType(astTypeToDotLL(define.DataType())),
		),
	}

	for _, node := range define.Expression.Content {
		switch node := node.(type) {
		case *ast.LiteralInt:
			statements = append(statements,
				dotll.NewCall(dotll.CallStore).
					ArgInt64(int64(node.Value)).
					ArgPtr(define.Name),
			)
		default:
			panic("unknown expression type")
		}
	}

	return statements
}

func astTypeToDotLL(astType ast.DataType) dotll.DataType {
	switch astType {
	case ast.DataInt:
		return dotll.DataInt64
	default:
		panic("unknown ast type")
	}
}
