package token

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (p *AstPrinter) Print(e Expr) string {
	return e.Accept(p).(string)
}

func (p *AstPrinter) parenthesize(op string, exprs ...Expr) string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(op)
	for _, expr := range exprs {
		sb.WriteString(" ")
		result := expr.Accept(p).(string)
		sb.WriteString(result)
	}
	sb.WriteString(")")
	return sb.String()
}

func (p *AstPrinter) VisitLiteral(e *Literal) interface{} {
	if e.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}

func (p *AstPrinter) VisitUnary(e *Unary) interface{} {
	op := ""
	if e.Operation.Lexeme != nil {
		op = *e.Operation.Lexeme
	}
	return p.parenthesize(op, e.Right)
}

func (p *AstPrinter) VisitBinary(e *Binary) interface{} {
	op := ""
	if e.Operation.Lexeme != nil {
		op = *e.Operation.Lexeme
	}
	return p.parenthesize(op, e.Left, e.Right)
}

func (p *AstPrinter) VisitGrouping(e *Grouping) interface{} {
	return p.parenthesize("grouping", e.Expression)
}
