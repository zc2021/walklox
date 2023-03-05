package interpreter

import (
	"devZ/lox/internal/tokens"
	"errors"
)

func (i *Interpreter) SetBiFunc(t *tokens.Token) error {
	fn := i.env.BinaryOp(t)
	if fn == nil {
		return errors.New("unrecognized token type for binary function")
	}
	tokens.SetOp(fn, t)
	return nil
}

func (i *Interpreter) SetUnFunc(t *tokens.Token) error {
	fn := i.env.UnaryOp(t)
	if fn == nil {
		return errors.New("unrecognized token type for unary function")
	}
	tokens.SetOp(fn, t)
	return nil
}
