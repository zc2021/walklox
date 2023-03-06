package environment

import (
	"devZ/lox/internal/expressions"
	"devZ/lox/internal/reporters"
	"devZ/lox/internal/tokens"
	"fmt"
)

type Execution struct {
	enclosing  *Execution
	values     map[string]interface{}
	returning  bool
	accum      *reporters.Accumulator
	ctx        reporters.ErrCtx
	unary_ops  []tokens.UnFunc
	binary_ops []tokens.BiFunc
}

func NewGlobal() *Execution {
	unfns, bifns := defaultOps()
	return &Execution{
		values:     make(map[string]interface{}),
		unary_ops:  unfns,
		binary_ops: bifns,
	}
}

func Copy(ex *Execution) *Execution {
	unfns, bifns := copyOps(ex)
	vals := make(map[string]interface{})
	for nm, v := range ex.values {
		vals[nm] = v
	}
	return &Execution{
		enclosing:  ex.enclosing,
		values:     vals,
		unary_ops:  unfns,
		binary_ops: bifns,
		accum:      ex.accum,
		ctx:        ex.ctx,
	}
}

func (ex *Execution) Block() *Execution {
	unfns, bifns := copyOps(ex)
	return &Execution{
		enclosing:  ex,
		values:     make(map[string]interface{}),
		accum:      ex.accum,
		ctx:        ex.ctx,
		unary_ops:  unfns,
		binary_ops: bifns,
	}
}

func (ex *Execution) Up() *Execution {
	return ex.enclosing
}

func (ex *Execution) SetCtx(ctx reporters.ErrCtx) {
	ex.ctx = ctx
}

func (ex *Execution) SetAccum(a *reporters.Accumulator) {
	ex.accum = a
}

func (ex *Execution) Define(nm string, val interface{}) {
	ex.values[nm] = val
}

func (ex *Execution) BinaryOp(t *tokens.Token) tokens.BiFunc {
	op := ex.binary_ops[t.ID()]
	if op != nil {
		return op
	}
	if ex.enclosing != nil {
		return ex.enclosing.BinaryOp(t)
	}
	return nil
}

func (ex *Execution) UnaryOp(t *tokens.Token) tokens.UnFunc {
	op := ex.unary_ops[t.ID()]
	if op != nil {
		return op
	}
	if ex.enclosing != nil {
		return ex.enclosing.UnaryOp(t)
	}
	return nil
}

func (ex *Execution) Get(vr *tokens.Token) interface{} {
	tid := vr.ID()
	if tid == tokens.STRING || tid == tokens.NUMBER {
		return vr.Literal()
	}
	nm := vr.Lexeme()
	val, set := ex.values[nm]
	if set {
		return val
	}
	if ex.enclosing != nil {
		return ex.enclosing.Get(vr)
	}
	loc := vr.Line()
	ex.undefinedVar(nm, loc)
	return nil
}

func (ex *Execution) Assign(as *expressions.Assignment, val interface{}) {
	nm := as.Name().Lexeme()
	_, set := ex.values[nm]
	if set {
		ex.values[nm] = val
		return
	}
	if ex.enclosing != nil {
		ex.enclosing.Assign(as, val)
		return
	}
	loc := as.Name().Line()
	ex.undefinedVar(nm, loc)
}

func (ex *Execution) undefinedVar(nm string, loc int) {
	msg := fmt.Sprintf("undefined variable %s referenced", nm)
	ex.accum.AddError(loc, msg, ex.ctx)
}

func (ex *Execution) StartRet() {
	ex.returning = true
}

func (ex *Execution) StopRet() {
	ex.returning = false
}

func (ex *Execution) Returning() bool {
	return ex.returning
}
