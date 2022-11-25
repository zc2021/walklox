// Package tools provides internal tools for development of GoWalkLox.
package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

//go:generate stringer -type=TokenType ../tokens
//go:generate go run gen_expressions.go
//go:generate gofmt -w ../expressions/expr_structs_ints.go
//go:generate go run gen_statements.go
//go:generate gofmt -w ../statements/stmt_structs_ints.go

type FieldGen func(name, param string, imp bool) FieldStr

var Fields = map[string]FieldGen{
	"token": func(name, param string, imp bool) FieldStr {
		pkg := ""
		if imp {
			pkg = "tokens."
		}
		return FieldStr{
			Name:  name,
			Param: param,
			Type:  fmt.Sprintf("*%sToken", pkg),
		}
	},
	"expression": func(name, param string, imp bool) FieldStr {
		pkg := ""
		if imp {
			pkg = "expressions."
		}
		return FieldStr{
			Name:  name,
			Param: param,
			Type:  fmt.Sprintf("%sExpr", pkg),
		}
	},
	"interface": func(name, param string, imp bool) FieldStr {
		return FieldStr{
			Name:  name,
			Param: param,
			Type:  "interface{}",
		}
	},
	"stmt": func(name, param string, imp bool) FieldStr {
		pkg := ""
		if imp {
			pkg = "statements."
		}
		return FieldStr{
			Name:  name,
			Param: param,
			Type:  fmt.Sprintf("%sStmt", pkg),
		}
	},
	"stmt_list": func(name, param string, imp bool) FieldStr {
		pkg := ""
		if imp {
			pkg = "statements."
		}
		return FieldStr{
			Name:  name,
			Param: param,
			Type:  fmt.Sprintf("[]%sStmt", pkg),
		}
	},
}

func GenerateVisitorPkgFile(nm string, ft *PkgTemplateData) {
	var methods, funcs, visitSigs []FuncStr

	for _, s := range ft.Structs {
		methods = append(methods, AcceptMethod(&s))
		methods = append(methods, Getters(&s)...)
		methods = append(methods, Setters(&s)...)
		funcs = append(funcs, ConstructorFunc(&s))
		visitSigs = append(visitSigs, VisitSig(&s))
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
