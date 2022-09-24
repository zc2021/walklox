package main

import (
	"bufio"
	"fmt"
	"os"

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

	err = run(script)
	if err != nil {
		os.Exit(65)
	}
}

func runPrompt(p string) {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(p)
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

func run(script []byte) error {
	inpt := scanner.NewScanner()
	toks := inpt.scanTokens()

	for tok := range toks {
		fmt.Println(tok)
	}
	return nil
}
