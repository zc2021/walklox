package main

import (
	"bufio"
	"fmt"
	"os"

	"devZ/lox/internal/environment"
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
		ec := int(reporters.COMMAND)
		os.Exit(ec)
	}
	if numArgs == 1 {
		ec := runFile(args[0])
		os.Exit(ec)
	}
	ec := runPrompt("> ")
	os.Exit(ec)
}

func runFile(filePath string) int {
	script, err := os.ReadFile(filePath)
	if err != nil {
		return int(reporters.FILE)
	}
	env := environment.NewGlobal()
	_, ex := run(script, env)
	return ex
}

func runPrompt(p string) int {
	input := bufio.NewReader(os.Stdin)
	env := environment.NewGlobal()
	for {
		fmt.Print(p)
		line, err := input.ReadBytes('\n')
		if string(line) == "bye!\n" {
			break
		}
		if err != nil {
			return int(reporters.REPL)
		}
		if line == nil {
			break
		}
		val, _ := run(line, env)
		fmt.Println(val)
	}
	return int(reporters.SUCCESS)
}

func run(script []byte, env *environment.Execution) (interface{}, int) {
	accum := reporters.NewAccumulator()
	inpt := scanner.New(script, accum)
	toks := inpt.ScanTokens()
	if checkErrs(accum) != nil {
		return nil, int(scanner.CTX)
	}
	prs := parser.New(toks, accum)
	stmts := prs.Parse()
	if checkErrs(accum) != nil {
		return nil, int(parser.CTX)
	}
	intpr := interpreter.New(accum, env)
	val := intpr.Interpret(stmts)
	if checkErrs(accum) != nil {
		return nil, int(interpreter.CTX)
	}
	return val, int(reporters.SUCCESS)
}

func checkErrs(a *reporters.Accumulator) error {
	if a.HasErrors() {
		err := a.LastError()
		a.PrintErrors()
		a.ResetErrors()
		return err
	}
	return nil
}
