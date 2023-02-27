package interpreter

import (
	"devZ/lox/internal/environment"
	"devZ/lox/internal/statements"
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
	env := environment.Copy(i.env)
	for i := 0; i < len(uf.declaration.Params()); i++ {
		env.Define(uf.declaration.Params()[i].Lexeme(), args[i])
	}
	return i.executeBlockIn(uf.declaration.Body(), env)
}

func (uf *userFunction) String() string {
	return fmt.Sprintf("<fn %s>", uf.declaration.Name().Lexeme())
}
