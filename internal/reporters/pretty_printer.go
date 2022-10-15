package reporters

import (
	"fmt"
	"strings"

	"devZ/lox/internal/expressions"
)

type PrettyPrinter struct{}

func (pp *PrettyPrinter) Print(ex expressions.Expr) {
	stringVal := ex.Accept(pp).(string)
	fmt.Println(stringVal)
}

func (pp *PrettyPrinter) VisitLiteral(li *expressions.Literal) interface{} {
	if li.Value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", li.Value)
}

func (pp *PrettyPrinter) VisitBinary(bi *expressions.Binary) interface{} {
	return pp.parenthesize(bi.Operator.Lexeme(), bi.Left, bi.Right)
}

func (pp *PrettyPrinter) VisitGrouping(gr *expressions.Grouping) interface{} {
	return pp.parenthesize("group", gr.Expression)
}

func (pp *PrettyPrinter) VisitUnary(un *expressions.Unary) interface{} {
	return pp.parenthesize(un.Operator.Lexeme(), un.Right)
}

func (pp *PrettyPrinter) parenthesize(name string, exprs ...expressions.Expr) interface{} {
	var bld strings.Builder
	bld.WriteString(fmt.Sprintf("(%s", name))
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		bld.WriteRune(' ')
		bld.WriteString(expr.Accept(pp).(string))
	}
	bld.WriteRune(')')
	return bld.String()
}
