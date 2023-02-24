package interpreter

import "time"

type clock struct{}

func (c *clock) arity() int { return 0 }

func (c *clock) call(i *Interpreter, args []interface{}) interface{} {
	return time.Now().UnixMilli() / 1000
}

func (c *clock) String() string { return "<native fn: clock>" }
