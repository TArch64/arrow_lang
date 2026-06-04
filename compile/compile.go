package compile

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"

	"arrow_lang/ast"
	"arrow_lang/config"
	"arrow_lang/errext"

	"tinygo.org/x/go-llvm"
)

func Compile(program *ast.Program, compilerConfig *config.Compiler) (err error) {
	compilation := &Compilation{
		config:  compilerConfig,
		program: program,
	}

	if err = initLLVM(compilation); err != nil {
		return errext.Tag("llvm init", err)
	}

	defer compilation.Dispose()

	dotll, err := generateDotLL(compilation)
	if err != nil {
		return errext.Tag("llvm dotll", err)
	}

	err = os.MkdirAll(path.Dir(compilerConfig.Output), 0644)
	if err != nil {
		return errext.Tag("llvm create output dir", err)
	}

	if compilerConfig.Debug {
		if err = writeDebug(compilation, dotll); err != nil {
			return errext.Tag("llvm write debug", err)
		}
	}

	return llvmCompile(compilation, dotll)
}

func writeDebug(
	compilation *Compilation,
	dotll llvm.Module,
) error {
	astInfo, err := json.MarshalIndent(compilation.program, "", "  ")
	if err != nil {
		return errext.Tag("llvm json indent", err)
	}

	dest := compilation.config.Output
	err = os.WriteFile(dest+".ast.json", astInfo, 0644)
	if err != nil {
		return errext.Tag("llvm write ast", err)
	}

	err = os.WriteFile(dest+".ll", []byte(dotll.String()), 0644)
	return errext.Tag("llvm dotll", err)
}

func initLLVM(compilation *Compilation) (err error) {
	if err = llvm.InitializeNativeTarget(); err != nil {
		err = errext.Tag("llvm init target", err)
		return
	}

	if err = llvm.InitializeNativeAsmPrinter(); err != nil {
		err = errext.Tag("llvm init asm printer", err)
		return
	}

	compilation.targetTriple = llvm.DefaultTargetTriple()
	target, err := llvm.GetTargetFromTriple(compilation.targetTriple)
	if err != nil {
		err = errext.Tag("llvm init target from triple", err)
		return
	}

	compilation.targetMachine = target.CreateTargetMachine(
		compilation.targetTriple,
		"generic",
		"",
		llvm.CodeGenLevelDefault,
		llvm.RelocPIC,
		llvm.CodeModelDefault,
	)

	compilation.targetData = compilation.targetMachine.CreateTargetData()
	return
}

func llvmCompile(compilation *Compilation, dotll llvm.Module) error {
	objFile, err := os.CreateTemp("", compilation.config.OutputFilename()+".*.o")
	if err != nil {
		return errext.Tag("llvm create temp object file", err)
	}

	defer func() {
		if err = objFile.Close(); err != nil {
			err = errext.Tag("llvm close temp object file", err)
			log.Println(err)
			return
		}

		if err = os.Remove(objFile.Name()); err != nil {
			err = errext.Tag("llvm remove temp object file", err)
			log.Println(err)
		}
	}()

	buf, err := compilation.targetMachine.EmitToMemoryBuffer(dotll, llvm.ObjectFile)
	if err != nil {
		return errext.Tag("llvm emit to memory buffer", err)
	}

	_, err = objFile.Write(buf.Bytes())
	if err != nil {
		return errext.Tag("llvm write temp object file", err)
	}

	cmd := exec.CommandContext(compilation.config.Ctx,
		"cc", objFile.Name(), "-o", compilation.config.Output,
	)

	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	return errext.Tag("llvm compile object file", err)
}
