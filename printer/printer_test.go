package printer

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
	"testing"
)

var (
	minus = "-"
	star  = "*"
)

func TestPrintPrimary(t *testing.T) {
	printer := &Ast{}
	expr := &token.Literal{Value: 123}
	got, err := expr.Accept(printer)
	if err != nil {
		t.Error(err)
	}
	if got != "123" {
		t.Fatalf("expect: %v got: %v", "", got)
	}
}

func TestPrintUnary(t *testing.T) {
	expr := &token.Unary{
		Operation: scanner.Token{TokenType: scanner.MINUS, Lexeme: &minus},
		Right:     &token.Literal{Value: 123},
	}
	printer := &Ast{}
	got, err := expr.Accept(printer)
	if err != nil {
		t.Error(err)
	}
	if got != "(- 123)" {
		t.Fatalf("expect: %v got: %v", "", got)
	}
}

func TestPrintExpr(t *testing.T) {
	// -123 * (45.67)
	expr := &token.Binary{
		Left: &token.Unary{
			Operation: scanner.Token{TokenType: scanner.MINUS, Lexeme: &minus},
			Right:     &token.Literal{Value: 123},
		},
		Operation: scanner.Token{
			TokenType: scanner.STAR,
			Lexeme:    &star,
		},
		Right: &token.Grouping{Expression: &token.Literal{Value: 45.67}},
	}
	printer := &Ast{}
	got, err := expr.Accept(printer)
	if err != nil {
		t.Error(err)
	}
	if got != "(* (- 123) (grouping 45.67))" {
		t.Fatalf("expect: %v got: %v", "", got)
	}
}
