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
}

func New(onError ErrorCallback) *Interpreter {
	return &Interpreter{onError}
}

func (i *Interpreter) Interpret(expression token.Expr) (string, error) {
	result, err := i.eval(expression)
	var intErr *RuntimeError
	if err != nil {
		if !errors.As(err, &intErr) {
			return "", err
		}
		i.errorCallback(err.(*RuntimeError))
		return "", nil
	}
	return stringify(result), nil
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

func (i *Interpreter) eval(expr token.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteral(expr *token.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnary(expr *token.Unary) (interface{}, error) {
	right, err := i.eval(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operation.TokenType {
	case scanner.BANG:
		return !i.isTruthy(right), nil
	case scanner.MINUS:
		if err = checkNumberOperand(expr.Operation, right); err != nil {
			return nil, err
		}
		return -1 * right.(float64), nil
	}
	//TODO: handle other cases error
	return nil, nil
}

func (i *Interpreter) VisitBinary(expr *token.Binary) (interface{}, error) {
	left, err := i.eval(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.eval(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operation.TokenType {
	case scanner.MINUS:
		if err := checkNumberOperands(expr.Operation, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case scanner.SLASH:
		if err := checkNumberOperands(expr.Operation, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case scanner.STAR:
		if err := checkNumberOperands(expr.Operation, left, right); err != nil {
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
				return nil, &RuntimeError{Token: &expr.Operation, Message: "Operands must be numbers."}
			}
		case string:
			r, ok := right.(string)
			if ok {
				return l + r, nil
			} else {
				return nil, &RuntimeError{Token: &expr.Operation, Message: "Operands must be strings."}
			}
		default:
			return nil, &RuntimeError{Token: &expr.Operation, Message: "Operands must be two numbers or two strings."}
		}
	case scanner.GREATER:
		if err := checkNumberOperands(expr.Operation, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case scanner.GREATER_EQUAL:
		if err := checkNumberOperands(expr.Operation, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case scanner.LESS:
		if err := checkNumberOperands(expr.Operation, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case scanner.LESS_EQUAL:
		if err := checkNumberOperands(expr.Operation, left, right); err != nil {
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

func (i *Interpreter) VisitGrouping(expr *token.Grouping) (interface{}, error) {
	return i.eval(expr.Expression)
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
