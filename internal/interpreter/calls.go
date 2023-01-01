package interpreter

type callable interface {
	arity() int
	call(*Interpreter, []interface{}) interface{}
}
