//go:build generate

package main

import (
	"os"

	"devZ/lox/internal/tools"
)

func main() {
	bibod := "op.SetBiFunc()"
	binary := tools.StructStr{
		Name:  "Binary",
		Param: "bi",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "left",
				Param: "lf",
				Type:  "Expr",
			},
			tools.FieldStr{
				Name:  "operator",
				Param: "op",
				Type:  "*tokens.Token",
				SetBd: []string{bibod},
			},
			tools.FieldStr{
				Name:  "right",
				Param: "rt",
				Type:  "Expr",
			},
		},
		CnstBd: []string{bibod},
	}

	group := tools.StructStr{
		Name:  "Grouping",
		Param: "gr",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "expression",
				Param: "ex",
				Type:  "Expr"}}}

	literal := tools.StructStr{
		Name:  "Literal",
		Param: "li",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "value",
				Param: "val",
				Type:  "interface{}"}}}

	unbod := "op.SetUnFunc()"
	unary := tools.StructStr{
		Name:  "Unary",
		Param: "un",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "operator",
				Param: "op",
				Type:  "*tokens.Token",
				SetBd: []string{unbod}},
			tools.FieldStr{
				Name:  "right",
				Param: "rt",
				Type:  "Expr"}},
		CnstBd: []string{unbod}}

	expression := tools.InterfaceStr{
		Name: "Expr",
		Sigs: []tools.FuncStr{
			tools.FuncStr{
				Name:   "Accept",
				Params: []string{"v Visitor"},
				Return: []string{"interface{}"},
			},
		},
	}

	strcs := []tools.StructStr{binary, group, literal, unary}

	var methods, funcs, visitSigs []tools.FuncStr

	av_void := false

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

	interfaces := []tools.InterfaceStr{expression, visitor}

	pkgInfo := tools.PkgTemplateData{
		Package:    "expressions",
		Imports:    []string{"devZ/lox/internal/tokens"},
		Structs:    strcs,
		Interfaces: interfaces,
		Methods:    methods,
		Functions:  funcs,
	}

	f, err := os.Create("../expressions/expr_structs_ints.go")
	if err != nil {
		panic(err)
	}
	tools.GeneratePkgFile(f, &pkgInfo)
}
