package parser

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestParserLiteral(t *testing.T) {
	p := NewParser(
		[]scanner.Token{
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	pp := token.AstPrinter{}
	if pp.Print(ast) != "123" {
		t.Fatalf("expect: %v, got: %v", "123", pp.Print(ast))
	}
}

func TestParserUnary(t *testing.T) {
	p := NewParser(
		[]scanner.Token{
			{scanner.MINUS, getStrPtr("-"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	pp := token.AstPrinter{}
	if pp.Print(ast) != "(- 123)" {
		t.Fatalf("expect: %v, got: %v", "(- 123)", pp.Print(ast))
	}
}

func TestParserTerm(t *testing.T) {
	p := NewParser(
		[]scanner.Token{
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("+"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	pp := token.AstPrinter{}
	if pp.Print(ast) != "(+ 321 123)" {
		t.Fatalf("expect: %v, got: %v", "(+ 321 123)", pp.Print(ast))
	}
}

func TestParserFactorial(t *testing.T) {
	p := NewParser(
		[]scanner.Token{
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("*"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	pp := token.AstPrinter{}
	if pp.Print(ast) != "(* 321 123)" {
		t.Fatalf("expect: %v, got: %v", "(* 321 123)", pp.Print(ast))
	}
}

func TestParserGrouping(t *testing.T) {
	p := NewParser(
		[]scanner.Token{
			{scanner.LEFT_BRACE, getStrPtr("("), nil, 1},
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("*"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.RIGHT_BRACE, getStrPtr(")"), nil, 1},
			{scanner.EOF, nil, nil, 1},
		},
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	pp := token.AstPrinter{}
	if pp.Print(ast) != "(grouping (* 321 123))" {
		t.Fatalf("expect: %v, got: %v", "(grouping (* 321 123)(", pp.Print(ast))
	}
}

func TestParserGroupingError(t *testing.T) {
	p := NewParser(
		[]scanner.Token{
			{scanner.LEFT_BRACE, getStrPtr("("), nil, 1},
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("*"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
	)
	ast, err := p.Parse()
	if err == nil {
		t.Fatal("Expecting parsing error")
	}
	if ast != nil {
		t.Fatal("Expecting nil ast")
	}
	t.Log(err)
	//TODO: check error
}

func getStrPtr(s string) *string {
	return &s
}
