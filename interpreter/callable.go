package interpreter

import (
	"errors"
	"fmt"
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
)

type LoxCallable interface {
	Arity() int
	Call(*Interpreter, []interface{}) (interface{}, error)
}

type loxFunction struct {
	declaration *token.FunctionStmt
	closure     *Environment
}

func NewLoxFunction(decl *token.FunctionStmt, env *Environment) LoxCallable {
	return &loxFunction{decl, env}
}

func (fn *loxFunction) bind(int *loxInstance) LoxCallable {
	env := NewScopeEnvironment(fn.closure)
	env.Define("this", int)
	return NewLoxFunction(fn.declaration, env)
}

func (fn *loxFunction) Arity() int {
	return len(fn.declaration.Params)
}

func (fn *loxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	env := NewScopeEnvironment(fn.closure)
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

type loxClass struct {
	name    string
	methods map[string]*loxFunction
}

func NewLoxClass(name string, methods map[string]*loxFunction) LoxCallable {
	return &loxClass{name, methods}
}

func (cls *loxClass) findMethod(name *scanner.Token) *loxFunction {
	return cls.methods[*name.Lexeme]
}

func (cls *loxClass) Arity() int {
	return 0
}

func (cls *loxClass) Call(*Interpreter, []interface{}) (interface{}, error) {
	return NewLoxInstance(cls), nil
}

func (cls *loxClass) String() string {
	return fmt.Sprintf("<class '%v'.>", cls.name)
}

type loxInstance struct {
	class  *loxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *loxClass) LoxCallable {
	return &loxInstance{class: class, fields: make(map[string]interface{}, 0)}
}

func (i *loxInstance) Get(name *scanner.Token) (interface{}, error) {
	if val, exist := i.fields[*name.Lexeme]; exist {
		return val, nil
	}
	method := i.class.findMethod(name)
	if method != nil {
		return method.bind(i), nil
	}
	return nil, &RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined property '%v'.", *name.Lexeme),
	}
}

func (i *loxInstance) Set(name *scanner.Token, value interface{}) {
	i.fields[*name.Lexeme] = value
}

func (i *loxInstance) Arity() int {
	return 0
}

func (i *loxInstance) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return nil, nil
}

func (i *loxInstance) String() string {
	return fmt.Sprintf("<'%v' instance.>", i.class.name)
}
