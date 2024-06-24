package token

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (p *AstPrinter) Print(e Expr) (string, error) {
	result, err := e.Accept(p)
	return result.(string), err
}

func (p *AstPrinter) parenthesize(op string, exprs ...Expr) (string, error) {
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

func (p *AstPrinter) VisitLiteral(e *Literal) (interface{}, error) {
	if e.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", e.Value), nil
}

func (p *AstPrinter) VisitUnary(e *Unary) (interface{}, error) {
	op := ""
	if e.Operation.Lexeme != nil {
		op = *e.Operation.Lexeme
	}
	return p.parenthesize(op, e.Right)
}

func (p *AstPrinter) VisitBinary(e *Binary) (interface{}, error) {
	op := ""
	if e.Operation.Lexeme != nil {
		op = *e.Operation.Lexeme
	}
	return p.parenthesize(op, e.Left, e.Right)
}

func (p *AstPrinter) VisitGrouping(e *Grouping) (interface{}, error) {
	return p.parenthesize("grouping", e.Expression)
}
