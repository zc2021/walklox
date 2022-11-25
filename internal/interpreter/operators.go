package interpreter

import (
	"devZ/lox/internal/tokens"
	"errors"
)

func (i *Interpreter) SetBiFunc(t *tokens.Token) error {
	fn, ok := i.env.Get(t).(tokens.BiFunc)
	if !ok || fn == nil {
		return errors.New("unrecognized token type for binary function")
	}
	tokens.SetOp(fn, t)
	return nil
}

func (i *Interpreter) SetUnFunc(t *tokens.Token) error {
	fn, ok := i.env.Get(t).(tokens.UnFunc)
	if !ok || fn == nil {
		return errors.New("unrecognized token type for unary function")
	}
	tokens.SetOp(fn, t)
	return nil
}
