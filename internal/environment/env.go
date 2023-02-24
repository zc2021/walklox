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

func Copy(e *Execution) *Execution {
	unfns, bifns := copyOps(e)
	return &Execution{
		values:     make(map[string]interface{}),
		unary_ops:  unfns,
		binary_ops: bifns,
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

func (ex *Execution) Get(vr *tokens.Token) interface{} {
	un := ex.unary_ops[vr.ID()]
	if un != nil {
		return un
	}
	bi := ex.binary_ops[vr.ID()]
	if bi != nil {
		return bi
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
