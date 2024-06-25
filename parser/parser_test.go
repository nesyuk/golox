package parser

import (
	"github.com/nesyuk/golox/printer"
	"github.com/nesyuk/golox/scanner"
	"testing"
)

func TestParserLiteral(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
		testCallBack(&errors),
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(errors) != 0 {
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if ast == nil {
		t.Fatalf("expect not nil")
	}
	pp := printer.Ast{}
	result, err := pp.Print(ast)
	if err != nil {
		t.Error(err)
	}
	if result != "123" {
		t.Fatalf("expect: %v, got: %v", "123", result)
	}
}

func TestParserUnary(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			{scanner.MINUS, getStrPtr("-"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
		testCallBack(&errors),
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(errors) != 0 {
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if ast == nil {
		t.Fatal("expect not nil")
	}
	pp := printer.Ast{}
	result, err := pp.Print(ast)
	if err != nil {
		t.Error(err)
	}
	if result != "(- 123)" {
		t.Fatalf("expect: %v, got: %v", "(- 123)", result)
	}
}

func TestParserTerm(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("+"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
		testCallBack(&errors),
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if ast == nil {
		t.Fatal("expect not nil")
	}
	pp := printer.Ast{}
	result, err := pp.Print(ast)
	if err != nil {
		t.Error(err)
	}
	if result != "(+ 321 123)" {
		t.Fatalf("expect: %v, got: %v", "(+ 321 123)", result)
	}
}

func TestParserFactorial(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("*"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
		testCallBack(&errors),
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(errors) != 0 {
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if ast == nil {
		t.Fatal("expect not nil")
	}
	pp := printer.Ast{}
	result, err := pp.Print(ast)
	if err != nil {
		t.Error(err)
	}
	if result != "(* 321 123)" {
		t.Fatalf("expect: %v, got: %v", "(* 321 123)", result)
	}
}

func TestParserGrouping(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			{scanner.LEFT_BRACE, getStrPtr("("), nil, 1},
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("*"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.RIGHT_BRACE, getStrPtr(")"), nil, 1},
			{scanner.EOF, nil, nil, 1},
		},
		testCallBack(&errors),
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("expect: nil got: %v", err)
	}
	if len(errors) != 0 {
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if ast == nil {
		t.Fatal("expect not nil")
	}
	pp := printer.Ast{}
	result, err := pp.Print(ast)
	if err != nil {
		t.Error(err)
	}
	if result != "(grouping (* 321 123))" {
		t.Fatalf("expect: %v, got: %v", "(grouping (* 321 123)(", result)
	}
}

func TestParserGroupingError(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			{scanner.LEFT_BRACE, getStrPtr("("), nil, 1},
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("*"), nil, 1},
			{scanner.NUMBER, getStrPtr("123"), 123, 1},
			{scanner.EOF, nil, nil, 1},
		},
		testCallBack(&errors),
	)
	ast, err := p.Parse()
	if err != nil {
		t.Fatal("expect nil")
	}
	if len(errors) != 1 {
		t.Fatalf("expect len(1) got %d", len(errors))
	}
	if ast != nil {
		t.Fatal("expect nil")
	}
	if errors[0] != "expect ')' after expression." {
		t.Fatalf("expect %v got %v", "expect ')' after expression.", errors[0])
	}
}

func TestParserGroupingErrors(t *testing.T) {
	errors := make([]string, 0)
	t.Skip("sync is not implemented yet")
	p := NewParser(
		[]scanner.Token{
			{scanner.LEFT_BRACE, getStrPtr("("), nil, 1},
			{scanner.NUMBER, getStrPtr("321"), 321, 1},
			{scanner.MINUS, getStrPtr("-"), nil, 1},
			{scanner.LEFT_BRACE, getStrPtr("("), nil, 1},
			{scanner.EOF, nil, nil, 1},
		},
		testCallBack(&errors),
	)
	ast, err := p.Parse()
	if err == nil {
		t.Fatal("Expecting parsing error")
	}
	if len(errors) != 0 {
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if ast == nil {
		t.Fatal("Expecting nil ast")
	}
	t.Log(err)
	if err.Error() != "[line 1] Error at end: Expect ')' after expression." {
		t.Fatalf("expect: %v, got: %v", "[line 1] Error at end: Expect ')' after expression.", err.Error())
	}
}

func getStrPtr(s string) *string {
	return &s
}

var testCallBack = func(errs *[]string) ErrorCallback {
	return func(token scanner.Token, message string) {
		*errs = append(*errs, message)
	}
}
