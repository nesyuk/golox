package token

import (
	"github.com/nesyuk/golox/scanner"
	"testing"
)

var (
	minus = "-"
	star  = "*"
)

func TestPrintPrimary(t *testing.T) {
	printer := &AstPrinter{}
	expr := &Literal{Value: 123}
	got, err := expr.Accept(printer)
	if err != nil {
		t.Error(err)
	}
	if got != "123" {
		t.Fatalf("expect: %v got: %v", "", got)
	}
}

func TestPrintUnary(t *testing.T) {
	expr := &Unary{
		Operation: scanner.Token{TokenType: scanner.MINUS, Lexeme: &minus},
		Right:     &Literal{Value: 123},
	}
	printer := &AstPrinter{}
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
	expr := &Binary{
		Left: &Unary{
			Operation: scanner.Token{TokenType: scanner.MINUS, Lexeme: &minus},
			Right:     &Literal{Value: 123},
		},
		Operation: scanner.Token{
			TokenType: scanner.STAR,
			Lexeme:    &star,
		},
		Right: &Grouping{&Literal{Value: 45.67}},
	}
	printer := &AstPrinter{}
	got, err := expr.Accept(printer)
	if err != nil {
		t.Error(err)
	}
	if got != "(* (- 123) (grouping 45.67))" {
		t.Fatalf("expect: %v got: %v", "", got)
	}
}
