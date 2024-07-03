package token

import (
	"github.com/nesyuk/golox/scanner"
)

type Expr interface {
	Accept(visitor VisitorExpr) (interface{}, error)
}

type VisitorExpr interface {
	VisitAssign(expr *Assign) (interface{}, error)
	VisitLiteral(expr *Literal) (interface{}, error)
	VisitUnary(expr *Unary) (interface{}, error)
	VisitVariable(expr *Variable) (interface{}, error)
	VisitBinary(expr *Binary) (interface{}, error)
	VisitGrouping(expr *Grouping) (interface{}, error)
}

type Assign struct {
	Name scanner.Token
	Value Expr
}

func (e *Assign) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitAssign(e)
}

type Literal struct {
	Value interface{}
}

func (e *Literal) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitLiteral(e)
}

type Unary struct {
	Operator scanner.Token
	Right Expr
}

func (e *Unary) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitUnary(e)
}

type Variable struct {
	Name scanner.Token
}

func (e *Variable) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitVariable(e)
}

type Binary struct {
	Left Expr
	Operator scanner.Token
	Right Expr
}

func (e *Binary) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitBinary(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitGrouping(e)
}

type Stmt interface {
	Accept(visitor VisitorStmt) (interface{}, error)
}

type VisitorStmt interface {
	VisitBlock(stmt *Block) (interface{}, error)
	VisitExpression(stmt *Expression) (interface{}, error)
	VisitPrint(stmt *Print) (interface{}, error)
	VisitVar(stmt *Var) (interface{}, error)
}

type Block struct {
	Statements []Stmt
}

func (e *Block) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitBlock(e)
}

type Expression struct {
	Expression Expr
}

func (e *Expression) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitExpression(e)
}

type Print struct {
	Expression Expr
}

func (e *Print) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitPrint(e)
}

type Var struct {
	Name scanner.Token
	Initializer Expr
}

func (e *Var) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitVar(e)
}

