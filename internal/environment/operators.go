package environment

import (
	"devZ/lox/internal/tokens"
	"fmt"
)

var unaryFuncs = [...]tokens.UnFunc{
	tokens.MINUS: negate,
	tokens.BANG:  isNotTruthy,
}

var binaryFuncs = [...]tokens.BiFunc{
	tokens.PLUS:          add,
	tokens.MINUS:         subtract,
	tokens.SLASH:         divide,
	tokens.STAR:          multiply,
	tokens.GREATER:       greater,
	tokens.LESS:          less,
	tokens.GREATER_EQUAL: greaterEq,
	tokens.LESS_EQUAL:    lessEq,
	tokens.EQUAL_EQUAL:   equal,
	tokens.BANG_EQUAL:    notEqual,
}

func add(x, y interface{}) interface{} {
	n, isNum := x.(float64)
	if !isNum {
		return fmt.Sprintf("%s%s", x.(string), y.(string))
	}
	return n + y.(float64)
}

func subtract(x, y interface{}) interface{} {
	return x.(float64) - y.(float64)
}

func divide(x, y interface{}) interface{} {
	return x.(float64) / y.(float64)
}

func multiply(x, y interface{}) interface{} {
	return x.(float64) * y.(float64)
}

func greater(x, y interface{}) interface{} {
	return x.(float64) > y.(float64)
}

func less(x, y interface{}) interface{} {
	return x.(float64) < y.(float64)
}

func greaterEq(x, y interface{}) interface{} {
	return x.(float64) >= y.(float64)
}

func lessEq(x, y interface{}) interface{} {
	return x.(float64) <= y.(float64)
}

func checkVal[T any](y interface{}) interface{} {
	yT, ok := y.(T)
	if !ok {
		return nil
	}
	return yT
}

func equal(x, y interface{}) interface{} {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	var yt interface{}
	switch x.(type) {
	case bool:
		yt = checkVal[bool](y)
	case float64:
		yt = checkVal[float64](y)
	case string:
		yt = checkVal[string](y)
	default:
		return false
	}
	if yt == nil {
		return false
	}
	return yt == x
}

func notEqual(x, y interface{}) interface{} {
	return !equal(x, y).(bool)
}

func IsTruthy(x interface{}) interface{} {
	if x == nil {
		return false
	}
	xb, ok := x.(bool)
	if !ok {
		return true
	}
	return xb
}

func isNotTruthy(x interface{}) interface{} {
	return !IsTruthy(x).(bool)
}

func negate(x interface{}) interface{} {
	return -1 * x.(float64)
}

func defaultOps() (unfns []tokens.UnFunc, bifns []tokens.BiFunc) {
	unfns = unaryFuncs[:]
	bifns = binaryFuncs[:]
	return
}

func copyOps(e *Execution) (unfns []tokens.UnFunc, bifns []tokens.BiFunc) {
	unfns = e.unary_ops[:]
	bifns = e.binary_ops[:]
	return
}
