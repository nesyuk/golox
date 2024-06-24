package printer

import (
	"fmt"
	"github.com/nesyuk/golox/token"
	"strings"
)

type Ast struct {
}

func (p *Ast) Print(e token.Expr) (string, error) {
	result, err := e.Accept(p)
	return result.(string), err
}

func (p *Ast) parenthesize(op string, exprs ...token.Expr) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(op)
	for _, expr := range exprs {
		sb.WriteString(" ")
		result, err := expr.Accept(p)
		if err != nil {
			return "", err
		}
		sb.WriteString(result.(string))
	}
	sb.WriteString(")")
	return sb.String(), nil
}

func (p *Ast) VisitLiteral(e *token.Literal) (interface{}, error) {
	if e.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", e.Value), nil
}

func (p *Ast) VisitUnary(e *token.Unary) (interface{}, error) {
	op := ""
	if e.Operation.Lexeme != nil {
		op = *e.Operation.Lexeme
	}
	return p.parenthesize(op, e.Right)
}

func (p *Ast) VisitBinary(e *token.Binary) (interface{}, error) {
	op := ""
	if e.Operation.Lexeme != nil {
		op = *e.Operation.Lexeme
	}
	return p.parenthesize(op, e.Left, e.Right)
}

func (p *Ast) VisitGrouping(e *token.Grouping) (interface{}, error) {
	return p.parenthesize("grouping", e.Expression)
}
