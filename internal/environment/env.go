package environment

import (
	"devZ/lox/internal/expressions"
	"devZ/lox/internal/reporters"
	"devZ/lox/internal/tokens"
	"fmt"
)

type Execution struct {
	enclosing *Execution
	values    map[string]interface{}
	accum     *reporters.Accumulator
	ctx       reporters.ErrCtx
}

func NewGlobal() *Execution {
	return &Execution{
		values: make(map[string]interface{}),
	}
}

func (ex *Execution) Block() *Execution {
	return &Execution{
		enclosing: ex,
		values:    make(map[string]interface{}),
		accum:     ex.accum,
		ctx:       ex.ctx,
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
	msg := fmt.Sprintf("Undefined variable %s referenced at %d.", nm, loc)
	ex.accum.AddError(loc, msg, ex.ctx)
}
