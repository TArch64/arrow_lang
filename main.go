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
	var input, output string
	var debug bool
	flag.StringVar(&input, "i", "", "input file name")
	flag.StringVar(&output, "o", "", "output file name")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	if input == "" || output == "" {
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.OpenFile(input, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	program, err := ast.Parse(token.Read(file))
	if err != nil {
		log.Fatal(err)
	}

	err = compile.Compile(program, &config.Compiler{
		Output: output,
		Debug:  debug,
		Ctx:    context.Background(),
	})

	if err != nil {
		log.Fatal(err)
	}
}
