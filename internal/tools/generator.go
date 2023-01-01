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

type FieldGen func(name, param string, impt, list bool) FieldStr

func bare(name, param, tp string) FieldStr {
	return FieldStr{
		Name:  name,
		Param: param,
		Type:  tp,
	}
}

func sgl(name, param, pkg, tp string) FieldStr {
	tp = fmt.Sprintf("%s%s", pkg, tp)
	return bare(name, param, tp)
}

func list(name, param, pkg, tp string) FieldStr {
	tp = fmt.Sprintf("[]%s%s", pkg, tp)
	return bare(name, param, tp)
}

var Fields = map[string]FieldGen{
	"token": func(name, param string, impt, lst bool) FieldStr {
		pkg := "*"
		if impt {
			pkg = "*tokens."
		}
		tp := "Token"
		if lst {
			return list(name, param, pkg, tp)
		}
		return sgl(name, param, pkg, tp)
	},
	"expression": func(name, param string, impt, lst bool) FieldStr {
		pkg := ""
		if impt {
			pkg = "expressions."
		}
		tp := "Expr"
		if lst {
			return list(name, param, pkg, tp)
		}
		return sgl(name, param, pkg, tp)
	},
	"interface": func(name, param string, imp, lst bool) FieldStr {
		tp := "interface{}"
		if lst {
			return list(name, param, "", tp)
		}
		return bare(name, param, tp)
	},
	"stmt": func(name, param string, imp, lst bool) FieldStr {
		pkg := ""
		if imp {
			pkg = "statements."
		}
		tp := "Stmt"
		if lst {
			return list(name, param, pkg, tp)
		}
		return sgl(name, param, pkg, tp)
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
