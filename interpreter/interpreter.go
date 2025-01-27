package interpreter

import (
	"errors"
	"fmt"
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
	"strings"
)

type Interpreter struct {
	errorCallback ErrorCallback
	printCallback PrintCallback
	globals       *Environment
	env           *Environment
	locals        map[token.Expr]int
}

func New(onError ErrorCallback, onPrint PrintCallback) *Interpreter {
	globals := NewEnvironment()
	globals.Define("clock", clock{})
	return &Interpreter{onError, onPrint, globals, globals, make(map[token.Expr]int, 0)}
}

func (i *Interpreter) Interpret(statements []token.Stmt) error {
	for _, stmt := range statements {
		_, err := i.exec(stmt)

		var intErr *RuntimeError
		if err != nil {
			if !errors.As(err, &intErr) {
				return err
			}
			i.errorCallback(err.(*RuntimeError))
			return nil
		}
	}
	return nil
}

func stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}
	str := fmt.Sprintf("%v", value)
	if strings.HasSuffix(str, ".0") {
		str = str[:2]
	}
	return str
}

func (i *Interpreter) Resolve(expr token.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) eval(expr token.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) exec(stmt token.Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) execBlock(statements []token.Stmt, env *Environment) (interface{}, error) {
	prev := i.env
	i.env = env

	defer func() {
		i.env = prev
	}()

	for _, stmt := range statements {
		if _, err := i.exec(stmt); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *token.LogicalExpr) (interface{}, error) {
	left, err := i.eval(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.TokenType == scanner.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else if !i.isTruthy(left) {
		// logical 'and'
		return left, nil
	}
	return i.eval(expr.Right)
}

func (i *Interpreter) VisitAssignExpr(expr *token.AssignExpr) (interface{}, error) {
	value, err := i.eval(expr.Value)
	if err != nil {
		return nil, err
	}
	if distance, exist := i.locals[expr]; exist {
		i.env.AssignAt(distance, &expr.Name, value)
		return nil, nil
	}
	if err := i.globals.Assign(&expr.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitSetExpr(expr *token.SetExpr) (interface{}, error) {
	obj, err := i.eval(expr.Object)
	if err != nil {
		return nil, err
	}
	inst, ok := obj.(*loxInstance)
	if !ok {
		return nil, &RuntimeError{
			Token:   expr.Name,
			Message: "Only instances have fields.",
		}
	}
	val, err := i.eval(expr.Value)
	if err != nil {
		return nil, err
	}
	inst.Set(expr.Name, val)
	return val, nil
}

func (i *Interpreter) VisitClassStmt(stmt *token.ClassStmt) (interface{}, error) {
	var supercls *loxClass
	if stmt.Superclass != nil {
		obj, err := i.eval(stmt.Superclass)
		if err != nil {
			return nil, err
		}
		var ok bool
		supercls, ok = obj.(*loxClass)
		if !ok {
			return nil, &RuntimeError{&stmt.Superclass.Name, "Superclass must be a class."}
		}
	}

	i.env.Define(*stmt.Name.Lexeme, nil)
	// Defining in two steps allows methods to use their class name

	if stmt.Superclass != nil {
		i.env = NewScopeEnvironment(i.env)
		i.env.Define("super", supercls)
	}

	methods := make(map[string]*loxFunction, 0)
	for _, method := range stmt.Methods {
		fn := NewLoxFunction(method, i.env, *method.Name.Lexeme == "init")
		methods[*method.Name.Lexeme] = fn.(*loxFunction)
	}

	class := NewLoxClass(*stmt.Name.Lexeme, supercls, methods)

	if stmt.Superclass != nil {
		i.env = i.env.enclosing
	}

	if err := i.env.Assign(stmt.Name, class); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitGetExpr(expr *token.GetExpr) (interface{}, error) {
	obj, err := i.eval(expr.Object)
	if err != nil {
		return nil, err
	}
	inst, ok := obj.(*loxInstance)
	if !ok {
		return nil, &RuntimeError{
			Token:   expr.Name,
			Message: "Only instances have properties.",
		}
	}
	return inst.Get(expr.Name)
}

func (i *Interpreter) VisitSuperExpr(expr *token.SuperExpr) (interface{}, error) {
	distance := i.locals[expr]
	superCls := i.env.GetAt(distance, "super").(*loxClass)
	instance := i.env.GetAt(distance-1, "this").(*loxInstance)
	method := superCls.findMethod(*expr.Method.Lexeme)
	if method == nil {
		return nil, &RuntimeError{&expr.Method, fmt.Sprintf("Undefined property '%v'.", *expr.Method.Lexeme)}
	}
	return method.bind(instance), nil
}

func (i *Interpreter) VisitThisExpr(expr *token.ThisExpr) (interface{}, error) {
	return i.lookupVariable(&expr.Keyword, expr)
}

func (i *Interpreter) VisitExpressionStmt(stmt *token.ExpressionStmt) (interface{}, error) {
	_, err := i.eval(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(stmt *token.PrintStmt) (interface{}, error) {
	result, err := i.eval(stmt.Expression)
	if err != nil {
		return nil, err
	}
	i.printCallback(stringify(result))
	return nil, nil
}

func (i *Interpreter) VisitReturnStmt(stmt *token.ReturnStmt) (interface{}, error) {
	var value interface{}
	if stmt.Value != nil {
		var err error
		value, err = i.eval(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, &ReturnException{Value: value}
}

func (i *Interpreter) VisitWhileStmt(stmt *token.WhileStmt) (interface{}, error) {
	cond, err := i.eval(stmt.Condition)
	for ; err == nil && i.isTruthy(cond); cond, err = i.eval(stmt.Condition) {
		if _, err = i.exec(stmt.Body); err != nil {
			return nil, err
		}
	}
	return nil, err
}

func (i *Interpreter) VisitIfStmt(stmt *token.IfStmt) (interface{}, error) {
	val, err := i.eval(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(val) {
		return i.exec(stmt.ThenBranch)
	}
	if stmt.ElseBranch != nil {
		return i.exec(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *token.FunctionStmt) (interface{}, error) {
	fn := NewLoxFunction(stmt, i.env, false)
	i.env.Define(*stmt.Name.Lexeme, fn)
	return fn, nil
}

func (i *Interpreter) VisitVarStmt(stmt *token.VarStmt) (interface{}, error) {
	var value interface{}
	if stmt.Initializer != nil {
		var err error
		value, err = i.eval(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.env.Define(*stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(block *token.BlockStmt) (interface{}, error) {
	return i.execBlock(block.Statements, NewScopeEnvironment(i.env))
}

func (i *Interpreter) VisitLiteralExpr(expr *token.LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *token.UnaryExpr) (interface{}, error) {
	right, err := i.eval(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.TokenType {
	case scanner.BANG:
		return !i.isTruthy(right), nil
	case scanner.MINUS:
		if err = checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -1 * right.(float64), nil
	}
	//TODO: handle other cases error
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *token.BinaryExpr) (interface{}, error) {
	left, err := i.eval(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanner.MINUS:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case scanner.SLASH:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case scanner.STAR:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case scanner.PLUS:
		switch l := left.(type) {
		case float64:
			r, ok := right.(float64)
			if ok {
				return l + r, nil
			} else {
				return nil, &RuntimeError{Token: &expr.Operator, Message: fmt.Sprintf("Operands must be numbers: %v", right)}
			}
		case string:
			r, ok := right.(string)
			if ok {
				return l + r, nil
			} else {
				return nil, &RuntimeError{Token: &expr.Operator, Message: fmt.Sprintf("Operands must be strings: %v", right)}
			}
		default:
			return nil, &RuntimeError{Token: &expr.Operator, Message: "Operands must be two numbers or two strings."}
		}
	case scanner.GREATER:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case scanner.GREATER_EQUAL:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case scanner.LESS:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case scanner.LESS_EQUAL:
		if err := checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case scanner.BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case scanner.EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	}
	//TODO: handle error
	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr *token.CallExpr) (interface{}, error) {
	callee, err := i.eval(expr.Callee)
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, 0)
	for _, exArg := range expr.Arguments {
		arg, err := i.eval(exArg)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, &RuntimeError{Token: expr.Paren, Message: "Can only call functions and classes."}
	}
	if len(args) != function.Arity() {
		return nil, &RuntimeError{
			Token:   expr.Paren,
			Message: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(args)),
		}
	}
	return function.Call(i, args)
}

func (i *Interpreter) VisitGroupingExpr(expr *token.GroupingExpr) (interface{}, error) {
	return i.eval(expr.Expression)
}

func (i *Interpreter) VisitVariableExpr(expr *token.VariableExpr) (interface{}, error) {
	return i.lookupVariable(&expr.Name, expr)
}

func (i *Interpreter) lookupVariable(name *scanner.Token, expr token.Expr) (interface{}, error) {
	if distance, exist := i.locals[expr]; exist {
		return i.env.GetAt(distance, *name.Lexeme), nil
	}
	val, _ := i.globals.Get(name)
	return val, nil
}

func (i *Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if val, isBool := value.(bool); isBool {
		return val
	}
	return true
}

func (i *Interpreter) isEqual(left, right interface{}) bool {
	switch v1 := left.(type) {
	case nil:
		if right == nil {
			return true
		}
		return false
	case bool:
		if v2, ok := right.(bool); ok {
			return v1 == v2
		}
		return false
	case float64:
		if v2, ok := right.(float64); ok {
			return v1 == v2
		}
		return false
	case string:
		if v2, ok := right.(string); ok {
			return v1 == v2
		}
		return false
	}

	return false
}

func checkNumberOperand(operator scanner.Token, operand interface{}) error {
	if _, ok := operand.(float64); !ok {
		return &RuntimeError{Token: &operator, Message: "Operand must be a number."}
	}
	return nil
}

func checkNumberOperands(operator scanner.Token, left, right interface{}) error {
	_, lok := left.(float64)
	_, rok := right.(float64)
	if !lok || !rok {
		return &RuntimeError{Token: &operator, Message: "Operands must be a numbers."}
	}
	return nil
}

type RuntimeError struct {
	Token   *scanner.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return e.Message
}

type ErrorCallback func(*RuntimeError)

type ReturnException struct {
	Value interface{}
}

func (e *ReturnException) Error() string {
	return fmt.Sprintf("%v", e.Value)
}

type PrintCallback = func(string)
