package resolver

import (
	"github.com/nesyuk/golox/interpreter"
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
)

type Resolver struct {
	scopes        []map[string]bool
	currentFn     FunctionType
	currentCls    ClassType
	interpreter   *interpreter.Interpreter
	errorCallback ErrorCallback
}

func New(i *interpreter.Interpreter, onError ErrorCallback) *Resolver {
	return &Resolver{make([]map[string]bool, 0), FN_NONE, CLS_NONE, i, onError}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) resolveStmt(stmt token.Stmt) (interface{}, error) {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr token.Expr) (interface{}, error) {
	return expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr token.Expr, name scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, exist := r.scopes[i][*name.Lexeme]; exist {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) declare(name *scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := &r.scopes[len(r.scopes)-1]
	if _, exist := (*scope)[*name.Lexeme]; exist {
		r.errorCallback(*name, "Already a variable with this name in this scope.")
		return
	}
	r.scopes[len(r.scopes)-1][*name.Lexeme] = false
}

func (r *Resolver) define(name *scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][*name.Lexeme] = true
}

func (r *Resolver) Resolve(stmts []token.Stmt) (interface{}, error) {
	for _, s := range stmts {
		if _, err := r.resolveStmt(s); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *token.AssignExpr) (interface{}, error) {
	if _, err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr *token.SetExpr) (interface{}, error) {
	if _, err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}
	if _, err := r.resolveExpr(expr.Object); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr *token.GetExpr) (interface{}, error) {
	return r.resolveExpr(expr.Object)
}

func (r *Resolver) VisitThisExpr(expr *token.ThisExpr) (interface{}, error) {
	if r.currentCls == CLS_NONE {
		r.errorCallback(expr.Keyword, "Can't use 'this' outside of a class.")
		return nil, nil
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *token.LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *token.LogicalExpr) (interface{}, error) {
	if _, err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if _, err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *token.UnaryExpr) (interface{}, error) {
	if _, err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *token.CallExpr) (interface{}, error) {
	if _, err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}
	for _, arg := range expr.Arguments {
		if _, err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *token.VariableExpr) (interface{}, error) {
	if len(r.scopes) != 0 {
		isDefined, exist := r.scopes[len(r.scopes)-1][*expr.Name.Lexeme]
		if exist && !isDefined {
			r.errorCallback(expr.Name, "Can't read local variable in its own initializer.")
			return nil, nil
		}
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *token.BinaryExpr) (interface{}, error) {
	if _, err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if _, err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *token.GroupingExpr) (interface{}, error) {
	return r.resolveExpr(expr.Expression)
}

func (r *Resolver) VisitBlockStmt(stmt *token.BlockStmt) (interface{}, error) {
	r.beginScope()
	if _, err := r.Resolve(stmt.Statements); err != nil {
		return nil, err
	}
	r.endScope()
	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt *token.ClassStmt) (interface{}, error) {
	enclosingCls := r.currentCls
	r.currentCls = CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	if stmt.Superclass != nil && *(stmt.Name.Lexeme) == *(stmt.Superclass.Name.Lexeme) {
		r.errorCallback(stmt.Superclass.Name, "A class can't inherit from itself.")
		return nil, nil
	}

	if stmt.Superclass != nil {
		r.resolveExpr(stmt.Superclass)
	}

	for _, met := range stmt.Methods {
		declaration := METHOD
		if *met.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		if _, err := r.resolveFunction(met, declaration); err != nil {
			return nil, err
		}
	}

	r.endScope()
	r.currentCls = enclosingCls
	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt *token.ExpressionStmt) (interface{}, error) {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitFunctionStmt(stmt *token.FunctionStmt) (interface{}, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	return r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) resolveFunction(stmt *token.FunctionStmt, fnType FunctionType) (interface{}, error) {
	enclosingFn := r.currentFn
	r.currentFn = fnType
	r.beginScope()
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)
	}
	if _, err := r.Resolve(stmt.Body); err != nil {
		return nil, err
	}
	r.endScope()
	r.currentFn = enclosingFn
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt *token.IfStmt) (interface{}, error) {
	if _, err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	if _, err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		if _, err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt *token.PrintStmt) (interface{}, error) {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt *token.ReturnStmt) (interface{}, error) {
	if r.currentFn == FN_NONE {
		r.errorCallback(*stmt.Keyword, "Can't return from top-level code.")
		return nil, nil
	}
	if stmt.Value != nil {
		if r.currentFn == INITIALIZER {
			r.errorCallback(*stmt.Keyword, "Can't return a value from initializer.")
			return nil, nil
		}
		return r.resolveExpr(stmt.Value)
	}
	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt *token.WhileStmt) (interface{}, error) {
	if _, err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	if _, err := r.resolveStmt(stmt.Body); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt *token.VarStmt) (interface{}, error) {
	r.declare(&stmt.Name)
	if stmt.Initializer != nil {
		if _, err := r.resolveExpr(stmt.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(&stmt.Name)
	return nil, nil
}

type ErrorCallback = func(scanner.Token, string)

type FunctionType uint8

const (
	FN_NONE FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

type ClassType uint8

const (
	CLS_NONE ClassType = iota
	CLASS
)
