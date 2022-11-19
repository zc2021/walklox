// Package tools provides internal tools for development of GoWalkLox.
package tools

import (
	"os"
	"path/filepath"
)

//go:generate stringer -type=TokenType ../tokens
//go:generate go run gen_expressions.go
//go:generate gofmt -w ../expressions/expr_structs_ints.go
//go:generate go run gen_statements.go
//go:generate gofmt -w ../statements/stmt_structs_ints.go

func GenerateVisitorPkgFile(nm string, ft *PkgTemplateData, av_void bool) {
	var methods, funcs, visitSigs []FuncStr

	for _, s := range ft.Structs {
		methods = append(methods, AcceptMethod(&s, av_void))
		methods = append(methods, Getters(&s)...)
		methods = append(methods, Setters(&s)...)
		funcs = append(funcs, ConstructorFunc(&s))
		visitSigs = append(visitSigs, VisitSig(&s, av_void))
	}

	visitor := InterfaceStr{
		Name: "Visitor",
		Sigs: visitSigs,
	}

	ft.Interfaces = append(ft.Interfaces, visitor)
	ft.Methods = methods
	ft.Functions = funcs

	base_path := filepath.Join("..", ft.Package, nm)
	f_path, err := filepath.Abs(base_path)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(f_path)
	if err != nil {
		panic(err)
	}
	visitTemplate.Execute(f, ft)
}
