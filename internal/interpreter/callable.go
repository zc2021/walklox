package interpreter

import (
	"devZ/lox/internal/statements"
	"devZ/lox/internal/tokens"
	"fmt"
)

type callable interface {
	arity() int
	call(*Interpreter, []interface{}) interface{}
}

type userFunction struct {
	declaration *statements.Function
}

func newFunc(dec *statements.Function, i *Interpreter) *userFunction {
	return &userFunction{
		declaration: dec,
	}
}

func (uf *userFunction) arity() int {
	return len(uf.declaration.Params())
}

func (uf *userFunction) call(i *Interpreter, args []interface{}) interface{} {
	env := i.env.Block()
	for j := 0; j < len(uf.declaration.Params()); j++ {
		val := args[j]
		tok, ok := val.(*tokens.Token)
		if ok {
			val = i.env.Get(tok)
		}
		env.Define(uf.declaration.Params()[j].Lexeme(), val)
	}
	val := i.executeBlockIn(uf.declaration.Body(), env)
	return val
}

func (uf *userFunction) String() string {
	return fmt.Sprintf("<fn %s>", uf.declaration.Name().Lexeme())
}
