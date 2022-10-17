//go:build generate

package main

import (
	"os"

	"devZ/lox/internal/tools"
)

func main() {
	binary := tools.StructStr{
		Name:  "Binary",
		Param: "bi",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "left",
				Param: "lf",
				Type:  "Expr"},
			tools.FieldStr{
				Name:  "operator",
				Param: "op",
				Type:  "*tokens.Token"},
			tools.FieldStr{
				Name:  "right",
				Param: "rt",
				Type:  "Expr"}}}

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

	unary := tools.StructStr{
		Name:  "Unary",
		Param: "un",
		Fields: []tools.FieldStr{
			tools.FieldStr{
				Name:  "operator",
				Param: "op",
				Type:  "*tokens.Token"},
			tools.FieldStr{
				Name:  "right",
				Param: "rt",
				Type:  "Expr"}}}

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

	visitor := tools.InterfaceStr{
		Name: "Visitor",
		Sigs: []tools.FuncStr{
			tools.VisitSig(&binary),
			tools.VisitSig(&group),
			tools.VisitSig(&literal),
			tools.VisitSig(&unary),
		},
	}

	exprs := []tools.StructStr{binary, group, literal, unary}
	interfaces := []tools.InterfaceStr{expression, visitor}

	var methods []tools.FuncStr
	var funcs []tools.FuncStr

	for _, e := range exprs {
		methods = append(methods, tools.AcceptMethod(&e))
		methods = append(methods, tools.Getters(&e)...)
		methods = append(methods, tools.Setters(&e)...)
		funcs = append(funcs, tools.ConstructorFunc(&e))
	}

	pkgInfo := tools.PkgTemplateData{
		Package:    "expressions",
		Structs:    exprs,
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
