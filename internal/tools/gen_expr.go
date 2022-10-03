// The following build directive is necessary to make the package coherent:

//go:build generate

// This program generates Expr structs. It can be invoked using go generate.
// As of 10/2/22, this file will always show as in error in your IDE. No other package
// imports text/template, meaning the package isn't recognized as a valid import until
// go generate is used.
package main

import (
	"os"
	"text/template"
)

func main() {
	f, err := os.Create("../expressions/expr_structs.go")
	if err != nil {
		panic(err)
	}
	type expression struct {
		Name   string
		Fields []string
	}

	var exprs = []expression{
		expression{"Binary", []string{"Left Expr", "Operator tokens.TokID", "Right Expr"}},
		expression{"Grouping", []string{"Expression Expr"}},
		expression{"Literal", []string{"Value interface{}"}},
		expression{"Unary", []string{"Operator tokens.TokID", "Right Expr"}},
	}
	exprsTemplate.Execute(f, exprs)
}

var tmpString = `// Code generated by walklox/internal/tools/gen_expr.go. DO NOT EDIT.

package expressions

import "devZ/lox/internal/tokens"

type Expr struct {}

{{range $expr := .}}
type {{$expr.Name}} struct { {{range $field := $expr.Fields}}
	{{$field}} {{end}}
}
{{end}}
`

var exprsTemplate = template.Must(template.New("expressionStructs").Parse(tmpString))
