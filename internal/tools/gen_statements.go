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
			tools.Fields["expression"]("expression", "ex", true, false),
		},
	}

	prn := tools.StructStr{
		Name:  "Print",
		Param: "prnst",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("expression", "ex", true, false),
		},
	}

	varStmt := tools.StructStr{
		Name:  "VarStmt",
		Param: "varst",
		Fields: []tools.FieldStr{
			tools.Fields["token"]("name", "nm", true, false),
			tools.Fields["expression"]("initializer", "init", true, false),
		},
	}

	block := tools.StructStr{
		Name:  "Block",
		Param: "blk",
		Fields: []tools.FieldStr{
			tools.Fields["stmt"]("statements", "stmts", false, true),
		},
	}

	conditional := tools.StructStr{
		Name:  "If",
		Param: "ifst",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("condition", "cnd", true, false),
			tools.Fields["stmt"]("thenBranch", "thbr", false, false),
			tools.Fields["stmt"]("elseBranch", "elbr", false, false),
		},
	}

	whileStmt := tools.StructStr{
		Name:  "While",
		Param: "whst",
		Fields: []tools.FieldStr{
			tools.Fields["expression"]("condition", "cnd", true, false),
			tools.Fields["stmt"]("body", "bd", false, false),
		},
	}

	brk := tools.StructStr{
		Name:  "Break",
		Param: "brkst",
		Fields: []tools.FieldStr{
			tools.Fields["token"]("tok", "tk", true, false),
		},
	}

	statement_types := []tools.StructStr{
		expr,
		prn,
		varStmt,
		block,
		conditional,
		whileStmt,
		brk,
	}
	imps := []string{"devZ/lox/internal/expressions", "devZ/lox/internal/tokens"}

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

	pkgInfo := tools.PkgTemplateData{
		Package:    "statements",
		Imports:    imps,
		Structs:    statement_types,
		Interfaces: []tools.InterfaceStr{statement},
	}

	tools.GenerateVisitorPkgFile("stmt_structs_ints.go", &pkgInfo)
}
