package interpreter

import (
	"fmt"
	"github.com/nesyuk/golox/scanner"
)

type Environment struct {
	enclosing *Environment
	variables map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{nil, make(map[string]interface{})}
}

func NewScopeEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing, make(map[string]interface{})}
}

func (e *Environment) Define(name string, value interface{}) {
	e.variables[name] = value
}

func (e *Environment) Get(name *scanner.Token) (interface{}, error) {
	val, exist := e.variables[*name.Lexeme]
	if exist {
		return val, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, &RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined variable '%v'.", *name.Lexeme),
	}
}

func (e *Environment) Assign(name *scanner.Token, value interface{}) error {
	if _, exist := e.variables[*name.Lexeme]; exist {
		e.variables[*name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return &RuntimeError{name, fmt.Sprintf("Undefined variable '%v'", *name.Lexeme)}
}
