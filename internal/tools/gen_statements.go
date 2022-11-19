//go:build generate

package main

import (
	"devZ/lox/internal/tools"
)

func main() {
	fields := map[string]tools.FieldStr{
		"expression": tools.FieldStr{
			Name:  "expression",
			Param: "ex",
			Type:  "expressions.Expr",
		},
		"nm_tok": tools.FieldStr{
			Name:  "name",
			Param: "nm",
			Type:  "*tokens.Token",
		},
		"init_expr": tools.FieldStr{
			Name:  "initializer",
			Param: "init",
			Type:  "expressions.Expr",
		},
		"stmt_list": tools.FieldStr{
			Name:  "statements",
			Param: "sts",
			Type:  "[]Stmt",
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

	exprStmt := tools.StructStr{
		Name:  "Expression",
		Param: "expst",
		Fields: []tools.FieldStr{
			fields["expression"],
		},
	}

	prnStmt := tools.StructStr{
		Name:  "Print",
		Param: "prnst",
		Fields: []tools.FieldStr{
			fields["expression"],
		},
	}

	varStmt := tools.StructStr{
		Name:  "VarStmt",
		Param: "varst",
		Fields: []tools.FieldStr{
			fields["nm_tok"],
			fields["init_expr"],
		},
	}

	block := tools.StructStr{
		Name:  "Block",
		Param: "blk",
		Fields: []tools.FieldStr{
			fields["stmt_list"],
		},
	}

	statement_types := []tools.StructStr{exprStmt, prnStmt, varStmt, block}
	imps := []string{"devZ/lox/internal/expressions", "devZ/lox/internal/tokens"}

	pkgInfo := tools.PkgTemplateData{
		Package:    "statements",
		Imports:    imps,
		Structs:    statement_types,
		Interfaces: []tools.InterfaceStr{statement},
	}

	tools.GenerateVisitorPkgFile("stmt_structs_ints.go", &pkgInfo, true)
}
