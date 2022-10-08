//go:build generate

package main

import (
	"fmt"
	"os"

	"devZ/lox/internal/tools"
)

func visitSig(param, acceptorStr string) tools.FuncStr {
	casedAcceptor := tools.UpperString(acceptorStr)
	name := fmt.Sprintf("Visit%s", casedAcceptor)
	return tools.FuncStr{
		Name: name,
		Params: []string{
			fmt.Sprintf("%s *%s", param, acceptorStr),
		},
	}
}

func acceptMethod(param, recStr string) tools.FuncStr {
	casedAcceptor := tools.UpperString(recStr)
	visitName := fmt.Sprintf("Visit%s", casedAcceptor)
	return tools.FuncStr{
		Name:     "Accept",
		Receiver: fmt.Sprintf("%s *%s", param, recStr),
		Params:   []string{"v Visitor"},
		Body: []string{
			fmt.Sprintf("v.%s(%s)", visitName, param),
		},
	}
}

func main() {
	f, err := os.Create("../expressions/expr_structs_ints.go")
	if err != nil {
		panic(err)
	}

	exprs := []tools.StructStr{
		tools.StructStr{"Expr", []string{}},
		tools.StructStr{"Binary", []string{"Left Expr", "Operator tokens.TokID", "Right Expr"}},
		tools.StructStr{"Grouping", []string{"Expression Expr"}},
		tools.StructStr{"Literal", []string{"Value interface{}"}},
		tools.StructStr{"Unary", []string{"Operator tokens.TokID", "Right Expr"}},
	}

	interfaces := []tools.InterfaceStr{
		tools.InterfaceStr{
			Name: "Acceptor",
			Sigs: []tools.FuncStr{
				tools.FuncStr{
					Name:   "Accept",
					Params: []string{"v *Visitor"}},
			},
		},
		tools.InterfaceStr{
			Name: "Visitor",
			Sigs: []tools.FuncStr{
				visitSig("ex", "Expr"),
				visitSig("bi", "Binary"),
				visitSig("gr", "Grouping"),
				visitSig("li", "Literal"),
				visitSig("un", "Unary"),
			},
		},
	}

	methods := []tools.FuncStr{
		acceptMethod("ex", "Expr"),
		acceptMethod("bi", "Binary"),
		acceptMethod("gr", "Grouping"),
		acceptMethod("li", "Literal"),
		acceptMethod("un", "Unary"),
	}

	pkgInfo := tools.PkgTemplateData{
		Package:    "expressions",
		Structs:    exprs,
		Interfaces: interfaces,
		Methods:    methods,
	}

	tools.GeneratePkgFile(f, &pkgInfo)
}
