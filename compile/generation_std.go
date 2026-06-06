package compile

import (
	"tinygo.org/x/go-llvm"
)

type GenerationStd struct {
	targetData llvm.TargetData

	i8T     llvm.Type
	i32T    llvm.Type
	i64T    llvm.Type
	doubleT llvm.Type
	ptrT    llvm.Type
	voidT   llvm.Type

	malloc  llvm.Value
	mallocT llvm.Type

	free  llvm.Value
	freeT llvm.Type
}

func (g *Generation) newStd() {
	std := GenerationStd{
		targetData: g.targetData,

		i8T:     g.ctx.Int8Type(),
		i32T:    g.ctx.Int32Type(),
		i64T:    g.ctx.Int64Type(),
		doubleT: g.ctx.DoubleType(),
		voidT:   g.ctx.VoidType(),
	}

	std.ptrT = llvm.PointerType(std.i8T, 0)

	std.mallocT = llvm.FunctionType(std.ptrT, []llvm.Type{std.i64T}, false)
	std.malloc = llvm.AddFunction(g.mod, "malloc", std.mallocT)

	std.freeT = llvm.FunctionType(std.voidT, []llvm.Type{std.ptrT}, false)
	std.free = llvm.AddFunction(g.mod, "free", std.freeT)

	g.std = &std
}

func (s *GenerationStd) sizeOf(target llvm.Type) llvm.Value {
	size := s.targetData.TypeAllocSize(target)
	return llvm.ConstInt(s.i64T, size, false)
}
