package token

import (
	"github.com/nesyuk/golox/scanner"
)

type Expr interface {
	Accept(visitor Visitor) (interface{}, error)
}

type Visitor interface {
	VisitLiteral(expr *Literal) (interface{}, error)
	VisitUnary(expr *Unary) (interface{}, error)
	VisitBinary(expr *Binary) (interface{}, error)
	VisitGrouping(expr *Grouping) (interface{}, error)
}

type Literal struct {
	Value interface{}
}

func (e *Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteral(e)
}

type Unary struct {
	Operation scanner.Token
	Right     Expr
}

func (e *Unary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnary(e)
}

type Binary struct {
	Left      Expr
	Operation scanner.Token
	Right     Expr
}

func (e *Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinary(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGrouping(e)
}
