//go:build generate

package main

import (
	"os"

	"devZ/lox/internal/tools"
)

func main() {
	exprStmt := tools.StructStr{
		Name:  "Expression",
		Param: "expst",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "expression",
				Param: "ex",
				Type:  "expressions.Expr",
			},
		},
	}

	prnStmt := tools.StructStr{
		Name:  "Print",
		Param: "prnst",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "expression",
				Param: "ex",
				Type:  "expressions.Expr",
			},
		},
	}

	statement := tools.InterfaceStr{
		Name: "Stmt",
		Sigs: []tools.FuncStr{
			tools.FuncStr{
				Name:   "Accept",
				Params: []string{"v Visitor"},
			},
		},
	}

	strcs := []tools.StructStr{exprStmt, prnStmt}

	var methods, funcs, visitSigs []tools.FuncStr

	av_void := true

	for _, s := range strcs {
		methods = append(methods, tools.AcceptMethod(&s, av_void))
		methods = append(methods, tools.Getters(&s)...)
		methods = append(methods, tools.Setters(&s)...)
		funcs = append(funcs, tools.ConstructorFunc(&s))
		visitSigs = append(visitSigs, tools.VisitSig(&s, av_void))
	}

	visitor := tools.InterfaceStr{
		Name: "Visitor",
		Sigs: visitSigs,
	}

	interfaces := []tools.InterfaceStr{statement, visitor}

	pkgInfo := tools.PkgTemplateData{
		Package:    "statements",
		Imports:    []string{"devZ/lox/internal/expressions"},
		Structs:    strcs,
		Interfaces: interfaces,
		Methods:    methods,
		Functions:  funcs,
	}

	f, err := os.Create("../statements/stmt_structs_ints.go")
	if err != nil {
		panic(err)
	}
	tools.GeneratePkgFile(f, &pkgInfo)
}
