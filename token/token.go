// Code generated by "go run token/gen/main.go token"; DO NOT EDIT.

package token

import (
	"github.com/nesyuk/golox/scanner"
)

type Expr interface {
	Accept(visitor VisitorExpr) (interface{}, error)
}

type VisitorExpr interface {
	VisitAssignExpr(expr *AssignExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitLogicalExpr(expr *LogicalExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitVariableExpr(expr *VariableExpr) (interface{}, error)
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
}

type AssignExpr struct {
	Name scanner.Token
	Value Expr
}

func (e *AssignExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitAssignExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitLiteralExpr(e)
}

type LogicalExpr struct {
	Left Expr
	Operator scanner.Token
	Right Expr
}

func (e *LogicalExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitLogicalExpr(e)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right Expr
}

func (e *UnaryExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct {
	Name scanner.Token
}

func (e *VariableExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitVariableExpr(e)
}

type BinaryExpr struct {
	Left Expr
	Operator scanner.Token
	Right Expr
}

func (e *BinaryExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitGroupingExpr(e)
}

type Stmt interface {
	Accept(visitor VisitorStmt) (interface{}, error)
}

type VisitorStmt interface {
	VisitBlockStmt(stmt *BlockStmt) (interface{}, error)
	VisitExpressionStmt(stmt *ExpressionStmt) (interface{}, error)
	VisitIfStmt(stmt *IfStmt) (interface{}, error)
	VisitPrintStmt(stmt *PrintStmt) (interface{}, error)
	VisitVarStmt(stmt *VarStmt) (interface{}, error)
}

type BlockStmt struct {
	Statements []Stmt
}

func (e *BlockStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitBlockStmt(e)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e *ExpressionStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitExpressionStmt(e)
}

type IfStmt struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *IfStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitIfStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (e *PrintStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitPrintStmt(e)
}

type VarStmt struct {
	Name scanner.Token
	Initializer Expr
}

func (e *VarStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitVarStmt(e)
}

