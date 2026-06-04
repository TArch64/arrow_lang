package compile

import (
	"tinygo.org/x/go-llvm"
)

type GenerationStd struct {
	targetData llvm.TargetData

	i8T   llvm.Type
	i32T  llvm.Type
	i64T  llvm.Type
	ptrT  llvm.Type
	voidT llvm.Type

	malloc  llvm.Value
	mallocT llvm.Type

	free  llvm.Value
	freeT llvm.Type
}

func generateDotLLStd(generation *Generation) *GenerationStd {
	std := GenerationStd{
		targetData: generation.targetData,

		i8T:   generation.ctx.Int8Type(),
		i32T:  generation.ctx.Int32Type(),
		i64T:  generation.ctx.Int64Type(),
		voidT: generation.ctx.VoidType(),
	}

	std.ptrT = llvm.PointerType(std.i8T, 0)

	std.mallocT = llvm.FunctionType(std.ptrT, []llvm.Type{std.i64T}, false)
	std.malloc = llvm.AddFunction(generation.mod, "malloc", std.mallocT)

	std.freeT = llvm.FunctionType(std.voidT, []llvm.Type{std.ptrT}, false)
	std.free = llvm.AddFunction(generation.mod, "free", std.freeT)

	return &std
}

func (s *GenerationStd) sizeOf(target llvm.Type) llvm.Value {
	size := s.targetData.TypeAllocSize(target)
	return llvm.ConstInt(s.i64T, size, false)
}
