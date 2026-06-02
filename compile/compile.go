package compile

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"arrow_lang/ast"
	"arrow_lang/config"
)

func Compile(program *ast.Program, compilerConfig *config.Compiler) error {
	var tempDir, llFile string
	if compilerConfig.Debug {
		tempDir = path.Dir(compilerConfig.Output)
		llFile = fmt.Sprintf("%s.*.ll", path.Base(compilerConfig.Output))
	} else {
		tempDir = path.Join(os.TempDir(), "arrow_lang")
		llFile = "program.*.ll"
	}
	tempFile, err := os.CreateTemp(tempDir, llFile)
	if err != nil {
		return err
	}

	defer func() {
		if err = tempFile.Close(); err != nil {
			log.Println(err.Error())
		}
		if !compilerConfig.Debug {
			if err = os.Remove(tempFile.Name()); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	dotll := generateDotLL(program)
	_, err = tempFile.WriteString(dotll)
	if err != nil {
		return err
	}

	return llvmCompile(tempFile.Name(), compilerConfig)
}

func llvmCompile(input string, compilerConfig *config.Compiler) error {
	err := os.MkdirAll(path.Dir(compilerConfig.Output), 0644)
	if err != nil {
		return err
	}

	var clangArgs []string
	if compilerConfig.Debug {
		clangArgs = append(clangArgs, "-g", "-O0")
	}

	clangArgs = append(clangArgs,
		input,
		"-o", compilerConfig.Output,
	)

	cmd := exec.CommandContext(compilerConfig.Ctx, "clang", clangArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	log.Println(string(stdout))
	return nil
}
