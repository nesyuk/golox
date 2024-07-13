package interpreter

import (
	"errors"
	"fmt"
	"github.com/nesyuk/golox/token"
)

type LoxCallable interface {
	Arity() int
	Call(*Interpreter, []interface{}) (interface{}, error)
}

type loxFunction struct {
	declaration *token.FunctionStmt
}

func NewLoxFunction(decl *token.FunctionStmt) LoxCallable {
	return &loxFunction{decl}
}

func (fn *loxFunction) Arity() int {
	return len(fn.declaration.Params)
}

func (fn *loxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	env := NewScopeEnvironment(interpreter.globals)
	for i := range arguments {
		env.Define(*fn.declaration.Params[i].Lexeme, arguments[i])
	}
	_, err := interpreter.execBlock(fn.declaration.Body, env)
	var returnValue *ReturnException
	if err != nil && errors.As(err, &returnValue) {
		return returnValue.Value, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (fn *loxFunction) String() string {
	return fmt.Sprintf("<fn '%v'.>", *fn.declaration.Name.Lexeme)
}