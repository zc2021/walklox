package main

import (
	"bufio"
	"fmt"
	"os"

	"devZ/lox/internal/interpreter"
	"devZ/lox/internal/parser"
	"devZ/lox/internal/reporters"
	"devZ/lox/internal/scanner"
)

func main() {
	args := os.Args[1:]
	numArgs := len(args)
	if numArgs > 1 {
		fmt.Println("Usage: walklox [script]")
		os.Exit(64)
	}
	if numArgs == 1 {
		runFile(args[0])
		os.Exit(0)
	}
	runPrompt("> ")
	os.Exit(0)
}

func runFile(filePath string) {
	script, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	run(script)
}

func runPrompt(p string) {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(p)
		line, err := input.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		if line == nil {
			break
		}
		run(line)
	}
}

func run(script []byte) {
	accum := &reporters.Accumulator{}
	inpt := scanner.New(script, accum)
	toks := inpt.ScanTokens()
	if checkErrs(accum) {
		return
	}
	prs := parser.New(toks, accum)
	expr := prs.Parse()
	if checkErrs(accum) {
		return
	}
	intpr := interpreter.New(accum)
	intpr.Interpret(expr)
	checkErrs(accum)
}

func checkErrs(a *reporters.Accumulator) bool {
	if a.HasErrors() {
		a.PrintErrors()
		a.ResetErrors()
		return true
	}
	return false
}
