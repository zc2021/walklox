//go:build generate

package main

import (
	"devZ/lox/internal/tools"
)

func main() {
	fields := map[string]tools.FieldStr{
		"nm_tok": tools.FieldStr{
			Name:  "name",
			Param: "nm",
			Type:  "*tokens.Token",
		},
		"op_tok": tools.FieldStr{
			Name:  "operator",
			Param: "op",
			Type:  "*tokens.Token",
		},
		"expr": tools.FieldStr{
			Name:  "expression",
			Param: "ex",
			Type:  "Expr",
		},
		"right_expr": tools.FieldStr{
			Name:  "right",
			Param: "rt",
			Type:  "Expr",
		},
		"left_expr": tools.FieldStr{
			Name:  "left",
			Param: "lf",
			Type:  "Expr",
		},
		"val_expr": tools.FieldStr{
			Name:  "value",
			Param: "vl",
			Type:  "Expr",
		},
		"value": tools.FieldStr{
			Name:  "value",
			Param: "val",
			Type:  "interface{}",
		},
	}

	op_bodies := map[string]string{
		"unary":  "op.SetUnFunc()",
		"binary": "op.SetBiFunc()",
	}

	binary := tools.StructStr{
		Name:  "Binary",
		Param: "bi",
		Fields: []tools.FieldStr{
			fields["left_expr"],
			fields["op_tok"],
			fields["right_expr"],
		},
		CnstBd: []string{op_bodies["binary"]},
	}
	binary.Fields[1].SetBd = []string{op_bodies["binary"]}

	group := tools.StructStr{
		Name:  "Grouping",
		Param: "gr",
		Fields: []tools.FieldStr{
			fields["expr"],
		},
	}

	literal := tools.StructStr{
		Name:  "Literal",
		Param: "li",
		Fields: []tools.FieldStr{
			fields["value"],
		},
	}

	unary := tools.StructStr{
		Name:  "Unary",
		Param: "un",
		Fields: []tools.FieldStr{
			fields["op_tok"],
			fields["right_expr"],
		},
		CnstBd: []string{op_bodies["unary"]},
	}
	unary.Fields[0].SetBd = []string{op_bodies["unary"]}

	variable := tools.StructStr{
		Name:  "VarExpr",
		Param: "vr",
		Fields: []tools.FieldStr{
			fields["nm_tok"],
		},
	}

	assignment := tools.StructStr{
		Name:  "Assignment",
		Param: "as",
		Fields: []tools.FieldStr{
			fields["nm_tok"],
			fields["val_expr"],
		},
	}

	expr_types := []tools.StructStr{
		binary,
		group,
		literal,
		unary,
		variable,
		assignment,
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

	tools.GenerateVisitorPkgFile("expr_structs_ints.go", &pkgInfo, false)
}
