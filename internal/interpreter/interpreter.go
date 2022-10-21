package interpreter

import (
	"devZ/lox/internal/expressions"
	"devZ/lox/internal/reporters"
	"devZ/lox/internal/tokens"
	"fmt"
)

const CTX = reporters.INTERPRETING

type Interpreter struct {
	accum   *reporters.Accumulator
	Printer *reporters.PrettyPrinter
}

func New(a *reporters.Accumulator) *Interpreter {
	return &Interpreter{
		accum: a,
	}
}

func (i *Interpreter) Interpret(e expressions.Expr) {
	val := i.evaluate(e)
	if i.accum.HasErrors() {
		return
	}
	fmt.Println(stringify(val))
}

func stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}
	switch t := val.(type) {
	case string, bool, float64:
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("Cannot interpret %s", t)
	}
}

func (i *Interpreter) VisitBinary(bi *expressions.Binary) interface{} {
	lf := i.evaluate(bi.Left())
	rt := i.evaluate(bi.Right())
	opLoc := bi.Operator().Line()
	lfMsgNum := fmt.Sprintf("Expect a number before %s", bi.Operator())
	rtMsgNum := fmt.Sprintf("Expect a number after %s", bi.Operator())
	rtMsgStr := fmt.Sprintf("Expect a string after %s", bi.Operator())
	if bi.Operator().ID() == tokens.PLUS {
		lstr, ok := lf.(string)
		if ok {
			rstr, ok := i.checkStr(rt, opLoc, rtMsgStr)
			if !ok {
				return nil
			}
			return bi.Operator().BiFunc(lstr, rstr)
		}
	}
	ln, ok := i.checkNum(lf, opLoc, lfMsgNum)
	if !ok {
		return nil
	}
	rn, ok := i.checkNum(rt, opLoc, rtMsgNum)
	if !ok {
		return nil
	}
	return bi.Operator().BiFunc(ln, rn)
}

func (i *Interpreter) VisitGrouping(gr *expressions.Grouping) interface{} {
	return i.evaluate(gr.Expression())
}

func (i *Interpreter) VisitLiteral(li *expressions.Literal) interface{} {
	return li.Value()
}

func (i *Interpreter) VisitUnary(un *expressions.Unary) interface{} {
	rt := i.evaluate(un.Right())
	opLoc := un.Operator().Line()
	rtMsgNum := fmt.Sprintf("Expect a number after %s", un.Operator().String())
	if un.Operator().ID() == tokens.MINUS {
		_, ok := i.checkNum(rt, opLoc, rtMsgNum)
		if !ok {
			return nil
		}
	}
	return un.Operator().UnFunc(rt)
}

func (i *Interpreter) evaluate(e expressions.Expr) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) checkNum(x interface{}, loc int, msg string) (float64, bool) {
	n, ok := x.(float64)
	if !ok {
		i.accum.AddError(loc, msg, CTX)
		return -1, false
	}
	return n, true
}

func (i *Interpreter) checkStr(x interface{}, loc int, msg string) (string, bool) {
	s, ok := x.(string)
	if !ok {
		i.accum.AddError(loc, msg, CTX)
		return "", false
	}
	return s, true
}
