package interpreter

import (
	"devZ/lox/internal/environment"
	"devZ/lox/internal/expressions"
	"devZ/lox/internal/reporters"
	"devZ/lox/internal/statements"
	"devZ/lox/internal/tokens"
	"fmt"
)

const CTX = reporters.INTERPRETING

type Interpreter struct {
	accum   *reporters.Accumulator
	env     *environment.Execution
	Printer *reporters.PrettyPrinter
}

func New(a *reporters.Accumulator, e *environment.Execution) *Interpreter {
	e.SetCtx(CTX)
	e.SetAccum(a)
	return &Interpreter{
		accum: a,
		env:   e,
	}
}

func (i *Interpreter) Interpret(stmts []statements.Stmt) interface{} {
	var val interface{}
	for _, s := range stmts {
		val = i.execute(s)
		if i.accum.HasErrors() {
			return nil
		}
	}
	return val
}

func stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}
	switch t := val.(type) {
	case string, bool, float64:
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("cannot interpret %s", t)
	}
}

func biMsg(tp, loc string, val *tokens.Token) string {
	return fmt.Sprintf("expect a %s %s %s", tp, loc, val)
}

func (i *Interpreter) VisitBinary(bi *expressions.Binary) interface{} {
	lf := i.evaluate(bi.Left())
	rt := i.evaluate(bi.Right())
	opLoc := bi.Operator().Line()
	err := i.SetBiFunc(bi.Operator())
	if err != nil {
		i.accum.AddError(opLoc, err.Error(), CTX)
	}
	lfMsgNum := biMsg("number", "before", bi.Operator())
	rtMsgNum := biMsg("number", "after", bi.Operator())
	rtMsgStr := biMsg("string", "after", bi.Operator())
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

func (i *Interpreter) VisitUnary(un *expressions.Unary) interface{} {
	rt := i.evaluate(un.Right())
	opLoc := un.Operator().Line()
	err := i.SetUnFunc(un.Operator())
	if err != nil {
		i.accum.AddError(opLoc, err.Error(), CTX)
	}
	rtMsgNum := fmt.Sprintf("expect a number after %s", un.Operator().String())
	if un.Operator().ID() == tokens.MINUS {
		_, ok := i.checkNum(rt, opLoc, rtMsgNum)
		if !ok {
			return nil
		}
	}
	return un.Operator().UnFunc(rt)
}

func (i *Interpreter) VisitGrouping(gr *expressions.Grouping) interface{} {
	return i.evaluate(gr.Expression())
}

func (i *Interpreter) VisitLiteral(li *expressions.Literal) interface{} {
	return li.Value()
}

func (i *Interpreter) VisitVarExpr(vr *expressions.VarExpr) interface{} {
	return i.env.Get(vr.Name())
}

func (i *Interpreter) VisitAssignment(as *expressions.Assignment) interface{} {
	val := i.evaluate(as.Value())
	i.env.Assign(as, val)
	return val
}

func (i *Interpreter) VisitExpression(expst *statements.Expression) interface{} {
	return i.evaluate(expst.Expression())
}

func (i *Interpreter) VisitPrint(prnst *statements.Print) interface{} {
	val := i.evaluate(prnst.Expression())
	if i.accum.HasErrors() {
		return nil
	}
	fmt.Println(stringify(val))
	return nil
}

func (i *Interpreter) VisitVarStmt(varst *statements.VarStmt) interface{} {
	var val interface{}
	if varst.Initializer() != nil {
		val = i.evaluate(varst.Initializer())
	}
	i.env.Define(varst.Name().Lexeme(), val)
	return val
}

func (i *Interpreter) VisitBlock(blk *statements.Block) interface{} {
	i.env = i.env.Block()
	val := i.executeBlock(blk.Statements())
	i.env = i.env.Up()
	return val
}

func (i *Interpreter) VisitIf(ifst *statements.If) interface{} {
	var val interface{}
	if environment.IsTruthy(i.evaluate(ifst.Condition())).(bool) {
		val = i.execute(ifst.ThenBranch())
	} else if ifst.ElseBranch() != nil {
		val = i.execute(ifst.ElseBranch())
	}
	return val
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

func (i *Interpreter) evaluate(e expressions.Expr) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) execute(s statements.Stmt) interface{} {
	return s.Accept(i)
}

func (i *Interpreter) executeBlock(stmts []statements.Stmt) interface{} {
	var val interface{}
	err := i.accum.LastError()
	for _, stmt := range stmts {
		val = i.execute(stmt)
		if i.accum.LastError() != err {
			break
		}
	}
	return val
}
