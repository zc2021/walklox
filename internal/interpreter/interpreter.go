package interpreter

import "devZ/lox/internal/expressions"

type Interpreter struct{}

func (i *Interpreter) VisitBinary(bi *expressions.Binary) interface{} {

}

func (i *Interpreter) VisitGrouping(gr *expressions.Grouping) interface{} {

}

func (i *Interpreter) VisitLiteral(li *expressions.Literal) interface{} {
	return li.Value()
}

func (i *Interpreter) VisitUnary(un *expressions.Unary) interface{} {

}
