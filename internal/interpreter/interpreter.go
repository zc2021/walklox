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
	globals *environment.Execution
	env     *environment.Execution
	accum   *reporters.Accumulator
	Printer *reporters.PrettyPrinter
}

func New(a *reporters.Accumulator, e *environment.Execution) *Interpreter {
	e.SetCtx(CTX)
	e.SetAccum(a)
	e.Define("clock", &clock{})
	return &Interpreter{
		accum:   a,
		env:     e,
		globals: e,
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

func (i *Interpreter) VisitCall(cl *expressions.Call) interface{} {
	loc := cl.Paren().Line()
	callee := i.evaluate(cl.Callee())
	args := make([]interface{}, 0)
	for _, a := range cl.Arguments() {
		args = append(args, i.evaluate(a))
	}
	fn, ok := callee.(callable)
	if !ok {
		i.accum.AddError(loc, "attempt to call uncallable entity", CTX)
		return nil
	}
	if fn.arity() != len(args) {
		arityMsg := fmt.Sprintf("%d arguments", fn.arity())
		i.accum.AddError(loc, reporters.Expectation(arityMsg, "in", "call"), CTX)
	}
	return fn.call(i, args)
}

func (i *Interpreter) VisitBinary(bi *expressions.Binary) interface{} {
	opLoc := bi.Operator().Line()
	err := i.SetBiFunc(bi.Operator())
	if err != nil {
		i.accum.AddError(opLoc, err.Error(), CTX)
		return nil
	}
	lf := i.evaluate(bi.Left())
	rt := i.evaluate(bi.Right())
	lfMsgNum := reporters.Expectation("a number", "before", bi.Operator().String())
	rtMsgNum := reporters.Expectation("a number", "after", bi.Operator().String())
	rtMsgStr := reporters.Expectation("a string", "after", bi.Operator().String())
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
	opLoc := un.Operator().Line()
	err := i.SetUnFunc(un.Operator())
	if err != nil {
		i.accum.AddError(opLoc, err.Error(), CTX)
		return nil
	}
	rt := i.evaluate(un.Right())
	rtMsgNum := reporters.Expectation("a number", "after", un.Operator().String())
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

func (i *Interpreter) VisitLogical(lg *expressions.Logical) interface{} {
	left := i.evaluate(lg.Left())
	if lg.Operator().ID() == tokens.OR {
		if environment.IsTruthy(left).(bool) {
			return left
		}
	} else {
		if !environment.IsTruthy(left).(bool) {
			return left
		}
	}
	return i.evaluate(lg.Right())
}

func (i *Interpreter) VisitExpression(expst *statements.Expression) interface{} {
	return i.evaluate(expst.Expression())
}

func (i *Interpreter) VisitFunction(fnst *statements.Function) interface{} {
	fn := newFunc(fnst, i)
	i.env.Define(fn.declaration.Name().Lexeme(), fn)
	return nil
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

func (i *Interpreter) VisitWhile(whst *statements.While) interface{} {
	var val interface{}
	for environment.IsTruthy(i.evaluate(whst.Condition())).(bool) {
		val = i.execute(whst.Body())
		if i.checkBreak(val) {
			break
		}
	}
	return val
}

func (i *Interpreter) VisitBreak(brkst *statements.Break) interface{} {
	return brkst.Tok()
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

func (i *Interpreter) checkBreak(x interface{}) bool {
	tok, ok := x.(*tokens.Token)
	if !ok {
		return false
	}
	return tok.ID() == tokens.BREAK
}

func (i *Interpreter) evaluate(e expressions.Expr) interface{} {
	return e.Accept(i)
}

func (i *Interpreter) execute(s statements.Stmt) interface{} {
	return s.Accept(i)
}

func (i *Interpreter) executeBlock(stmts []statements.Stmt) interface{} {
	var val interface{}
	for _, stmt := range stmts {
		val = i.execute(stmt)
		if i.accum.HasErrors() {
			break
		}
	}
	return val
}

func (i *Interpreter) executeBlockIn(stmts []statements.Stmt, block_env *environment.Execution) interface{} {
	prev := i.env
	i.env = block_env
	val := i.executeBlock(stmts)
	i.env = prev
	return val
}
