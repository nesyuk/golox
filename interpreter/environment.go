package interpreter

import (
	"fmt"
	"github.com/nesyuk/golox/scanner"
)

type Environment struct {
	variables map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{make(map[string]interface{})}
}

func (e *Environment) Define(name string, value interface{}) {
	e.variables[name] = value
}

func (e *Environment) Get(name *scanner.Token) (interface{}, error) {
	val, exist := e.variables[*name.Lexeme]
	if !exist {
		return nil, &RuntimeError{
			Token:   name,
			Message: fmt.Sprintf("Undefined variable '%v'.", *name.Lexeme),
		}
	}
	return val, nil
}

func (e *Environment) Assign(name *scanner.Token, value interface{}) error {
	if _, exist := e.variables[*name.Lexeme]; exist {
		e.variables[*name.Lexeme] = value
		return nil
	}
	return &RuntimeError{name, fmt.Sprintf("Undefined variable '%v'", *name.Lexeme)}
}
