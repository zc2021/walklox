package interpreter

import (
	"devZ/lox/internal/environment"
	"devZ/lox/internal/statements"
	"devZ/lox/internal/tokens"
	"fmt"
)

type callable interface {
	arity() int
	call(*Interpreter, []interface{}) interface{}
}

type userFunction struct {
	closure     *environment.Execution
	declaration *statements.Function
}

func newFunc(dec *statements.Function, i *Interpreter) *userFunction {
	env := environment.Copy(i.env)
	fn := &userFunction{
		declaration: dec,
		closure:     env,
	}
	fn.closure.Define(dec.Name().Lexeme(), fn)
	return fn
}

func (uf *userFunction) arity() int {
	return len(uf.declaration.Params())
}

func (uf *userFunction) call(i *Interpreter, args []interface{}) interface{} {
	env := uf.closure.Block()
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
