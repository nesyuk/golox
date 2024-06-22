package token

import (
	"github.com/nesyuk/golox/scanner"
	"testing"
)

func TestPrint(t *testing.T) {
	printer := &AstPrinter{}
	// - 123 * (45.67)
	star := "*"
	minus := "-"
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
	_ = expr
	got := expr.Accept(printer)
	if got != "(* (- 123) (grouping 45.67))" {
		t.Fatalf("expect: %v got: %v", "", got)
	}
}
