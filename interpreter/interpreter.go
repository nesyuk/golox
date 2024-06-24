package interpreter

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
)

type Interpreter struct {
}

func (i *Interpreter) Eval(expr token.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteral(expr *token.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnary(expr *token.Unary) interface{} {
	right := i.Eval(expr.Right)
	switch expr.Operation.TokenType {
	case scanner.BANG:
		return !i.isTruthy(right)
	case scanner.MINUS:

		return -1 * right.(float64)
	}
	//TODO: handle other cases error
	return nil
}

func (i *Interpreter) VisitBinary(expr *token.Binary) interface{} {
	left := i.Eval(expr.Left)
	right := i.Eval(expr.Right)

	switch expr.Operation.TokenType {
	case scanner.MINUS:
		return left.(float64) - right.(float64)
	case scanner.SLASH:
		return left.(float64) / right.(float64)
	case scanner.STAR:
		return left.(float64) * right.(float64)
	case scanner.PLUS:
		switch l := left.(type) {
		case float64:
			r, ok := right.(float64)
			if ok {
				return l + r
			}
			//TODO: handle error if !ok
		case string:
			r, ok := right.(string)
			if ok {
				return l + r
			}
			//TODO: handle error if !ok
		}
	case scanner.GREATER:
		return left.(float64) > right.(float64)
	case scanner.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case scanner.LESS:
		return left.(float64) < right.(float64)
	case scanner.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case scanner.BANG_EQUAL:
		return !i.isEqual(left, right)
	case scanner.EQUAL_EQUAL:
		return i.isEqual(left, right)
	}
	//TODO: handle error
	return nil
}

func (i *Interpreter) VisitGrouping(expr *token.Grouping) interface{} {
	return i.Eval(expr.Expression)
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
