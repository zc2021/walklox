//go:build generate

package main

import (
	"devZ/lox/internal/tools"
)

func main() {
	expr := tools.StructStr{
		Name:  "Expression",
		Param: "expst",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("expression", "ex", true),
		},
	}

	prn := tools.StructStr{
		Name:  "Print",
		Param: "prnst",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("expression", "ex", true),
		},
	}

	varStmt := tools.StructStr{
		Name:  "VarStmt",
		Param: "varst",
		Fields: []tools.FieldStr{
			tools.Fields["token"]("name", "nm", true),
			tools.Fields["expression"]("initializer", "init", true),
		},
	}

	block := tools.StructStr{
		Name:  "Block",
		Param: "blk",
		Fields: []tools.FieldStr{
			tools.Fields["stmt_list"]("statements", "stmts", false),
		},
	}

	conditional := tools.StructStr{
		Name:  "If",
		Param: "ifst",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("condition", "cnd", true),
			tools.Fields["stmt"]("thenBranch", "thbr", false),
			tools.Fields["stmt"]("elseBranch", "elbr", false),
		},
	}

	statement := tools.InterfaceStr{
		Name: "Stmt",
		Sigs: []tools.FuncStr{
			tools.FuncStr{
				Name:   "Accept",
				Params: []string{"v Visitor"},
				Return: []string{"interface{}"},
			},
		},
	}

	statement_types := []tools.StructStr{expr, prn, varStmt, block, conditional}
	imps := []string{"devZ/lox/internal/expressions", "devZ/lox/internal/tokens"}

	pkgInfo := tools.PkgTemplateData{
		Package:    "statements",
		Imports:    imps,
		Structs:    statement_types,
		Interfaces: []tools.InterfaceStr{statement},
	}

	tools.GenerateVisitorPkgFile("stmt_structs_ints.go", &pkgInfo)
}
