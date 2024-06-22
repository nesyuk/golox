package token

import (
	"github.com/nesyuk/golox/scanner"
)

type Expr interface {
	Accept(visitor Visitor) interface{}
}

type Visitor interface {
	VisitLiteral(expr *Literal) interface{}
	VisitUnary(expr *Unary) interface{}
	VisitBinary(expr *Binary) interface{}
	VisitGrouping(expr *Grouping) interface{}
}

type Literal struct {
	Value interface{}
}

func (e *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(e)
}

type Unary struct {
	Operation scanner.Token
	Right Expr
}

func (e *Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnary(e)
}

type Binary struct {
	Left Expr
	Operation scanner.Token
	Right Expr
}

func (e *Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinary(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGrouping(e)
}

