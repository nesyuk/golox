package interpreter

import "time"

type clock struct {
}

func (fn *clock) Arity() int {
	return 0
}

func (fn *clock) Call(_ *Interpreter, _ []interface{}) interface{} {
	return time.Now().UnixMilli() / 1000.0
}

func (fn *clock) String() string {
	return "<native fn 'clock'>"
}
