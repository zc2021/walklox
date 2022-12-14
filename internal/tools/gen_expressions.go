//go:build generate

package main

import (
	"devZ/lox/internal/tools"
)

func main() {
	binary := tools.StructStr{
		Name:  "Binary",
		Param: "bi",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("left", "lf", false, false),
			tools.Fields["token"]("operator", "op", true, false),
			tools.Fields["expression"]("right", "rt", false, false),
		},
	}

	group := tools.StructStr{
		Name:  "Grouping",
		Param: "gr",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("expression", "ex", false, false),
		},
	}

	literal := tools.StructStr{
		Name:  "Literal",
		Param: "li",
		Fields: []tools.FieldStr{
			tools.Fields["interface"]("value", "val", false, false),
		},
	}

	unary := tools.StructStr{
		Name:  "Unary",
		Param: "un",
		Fields: []tools.FieldStr{
			tools.Fields["token"]("operator", "op", true, false),
			tools.Fields["expression"]("right", "rt", false, false),
		},
	}

	variable := tools.StructStr{
		Name:  "VarExpr",
		Param: "vr",
		Fields: []tools.FieldStr{
			tools.Fields["token"]("name", "nm", true, false),
		},
	}

	assignment := tools.StructStr{
		Name:  "Assignment",
		Param: "as",
		Fields: []tools.FieldStr{
			tools.Fields["token"]("name", "nm", true, false),
			tools.Fields["expression"]("value", "vl", false, false),
		},
	}

	logical := tools.StructStr{
		Name:  "Logical",
		Param: "lg",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("left", "lf", false, false),
			tools.Fields["token"]("operator", "op", true, false),
			tools.Fields["expression"]("right", "rt", false, false),
		},
	}

	call := tools.StructStr{
		Name:  "Call",
		Param: "cl",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("callee", "cle", false, false),
			tools.Fields["token"]("paren", "pn", true, false),
			tools.Fields["expression"]("arguments", "args", false, true),
		},
	}

	expr_types := []tools.StructStr{
		binary,
		group,
		literal,
		unary,
		variable,
		assignment,
		logical,
		call,
	}

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

	pkgInfo := tools.PkgTemplateData{
		Package:    "expressions",
		Imports:    []string{"devZ/lox/internal/tokens"},
		Structs:    expr_types,
		Interfaces: []tools.InterfaceStr{expression},
	}

	tools.GenerateVisitorPkgFile("expr_structs_ints.go", &pkgInfo)
}
