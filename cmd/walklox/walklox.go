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
	return run(script)
}

func runPrompt(p string) int {
	input := bufio.NewReader(os.Stdin)
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
		run(line)
	}
	return int(reporters.SUCCESS)
}

func run(script []byte) int {
	accum := &reporters.Accumulator{}
	inpt := scanner.New(script, accum)
	toks := inpt.ScanTokens()
	if checkErrs(accum) != nil {
		return int(scanner.CTX)
	}
	prs := parser.New(toks, accum)
	stmts := prs.Parse()
	if checkErrs(accum) != nil {
		return int(parser.CTX)
	}
	intpr := interpreter.New(accum)
	intpr.Interpret(stmts)
	if checkErrs(accum) != nil {
		return int(interpreter.CTX)
	}
	return 0
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
