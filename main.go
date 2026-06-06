package main

import (
	"context"
	"flag"
	"log"
	"os"

	"arrow_lang/ast"
	"arrow_lang/compile"
	"arrow_lang/config"
	"arrow_lang/token"
)

func main() {
	compilerConfig := &config.Compiler{
		Ctx: context.Background(),
	}

	flag.StringVar(&compilerConfig.Input, "i", "", "input file name")
	flag.StringVar(&compilerConfig.Output, "o", "", "output file name")
	flag.BoolVar(&compilerConfig.Debug, "debug", false, "debug mode")
	flag.Parse()

	if compilerConfig.Input == "" || compilerConfig.Output == "" {
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.OpenFile(compilerConfig.Input, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	program, err := ast.Parse(token.Read(file))
	if err != nil {
		log.Fatal(err)
	}

	if err = compile.Compile(program, compilerConfig); err != nil {
		log.Fatal(err)
	}
}
