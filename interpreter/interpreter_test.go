package interpreter

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestInterpretLiteralFloat(t *testing.T) {
	expr := &token.Literal{Value: 123.}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.Interpret(expr)
	if err != nil {
		t.Error(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}
	if got != "123" {
		t.Fatalf("expected '123' got '%v'", got)
	}
}

func TestInterpretLiteralString(t *testing.T) {
	expr := &token.Literal{Value: "abc"}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.Interpret(expr)
	if err != nil {
		t.Error(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}

	if got != "abc" {
		t.Fatalf("expected 'abc' got '%v'", got)
	}
}

func TestInterpretLiteralBoolean(t *testing.T) {
	expr := &token.Literal{Value: true}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.Interpret(expr)
	if err != nil {
		t.Error(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}
	if got != "true" {
		t.Fatalf("expected 'true' got '%v'", got)
	}
}

func TestInterpretUnaryMinus(t *testing.T) {
	expr := &token.Unary{
		Operation: scanner.Token{TokenType: scanner.MINUS},
		Right:     &token.Literal{Value: 123.0},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.Interpret(expr)
	if err != nil {
		t.Error(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}
	if got != "-123" {
		t.Fatalf("expect: '-123' got: '%v'", got)
	}
}

func TestInterpretUnaryError(t *testing.T) {
	expr := &token.Unary{
		Operation: scanner.Token{TokenType: scanner.MINUS},
		Right:     &token.Literal{Value: "muffin"},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.Interpret(expr)
	if err != nil {
		t.Fatalf("expect nil")
	}
	if len(errs) != 1 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}
	if got != "" {
		t.Fatalf("expect empty")
	}
	if errs[0] != "Operand must be a number." {
		t.Fatalf("expect %v got %v", "Operand must be a number.", err.Error())
	}
}

func TestInterpretUnaryBang(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect string
	}{
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: 123.0},
		}, "false"},
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: "abc"},
		}, "false"},
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: true},
		}, "false"},
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: false},
		}, "true"},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	for _, test := range tests {
		got, err := i.Interpret(test.expr)
		if err != nil {
			t.Error(err)
		}
		if len(errs) != 0 {
			t.Fatalf("expect len(1), got: %d", len(errs))
		}
		if got != test.expect {
			t.Fatalf("expect: %v got: %v", test.expect, got)
		}
	}
}

func TestInterpretBinary(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect string
	}{
		{&token.Binary{
			Left:      &token.Literal{Value: 40.02},
			Operation: scanner.Token{TokenType: scanner.PLUS},
			Right:     &token.Literal{Value: 2.40},
		}, "42.42"},
		{&token.Binary{
			Left:      &token.Literal{Value: "Hello, "},
			Operation: scanner.Token{TokenType: scanner.PLUS},
			Right:     &token.Literal{Value: "World!"},
		}, "Hello, World!"},
		{&token.Binary{
			Left:      &token.Literal{Value: 16.0},
			Operation: scanner.Token{TokenType: scanner.STAR},
			Right:     &token.Literal{Value: 2.0},
		}, "32"},
		{&token.Binary{
			Left:      &token.Literal{Value: 256.0},
			Operation: scanner.Token{TokenType: scanner.SLASH},
			Right:     &token.Literal{Value: 2.0},
		}, "128"},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.0},
			Operation: scanner.Token{TokenType: scanner.GREATER},
			Right:     &token.Literal{Value: 1.0},
		}, "true"},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.GREATER_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, "true"},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.0},
			Operation: scanner.Token{TokenType: scanner.LESS},
			Right:     &token.Literal{Value: 1.0},
		}, "false"},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.LESS_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, "true"},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.EQUAL_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, "true"},
		{&token.Binary{
			Left:      &token.Literal{Value: "lox"},
			Operation: scanner.Token{TokenType: scanner.EQUAL_EQUAL},
			Right:     &token.Literal{Value: "lox"},
		}, "true"},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, "false"},
		{&token.Binary{
			Left:      &token.Literal{Value: "lox"},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: "lox"},
		}, "false"},
		{&token.Binary{
			Left:      &token.Literal{Value: "lox"},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: nil},
		}, "true"},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: "lox"},
		}, "true"},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	for _, test := range tests {
		got, err := i.Interpret(test.expr)
		if err != nil {
			t.Error(err)
		}
		if len(errs) != 0 {
			t.Fatalf("expect len(1), got: %d", len(errs))
		}
		if got != test.expect {
			t.Fatalf("expect: %v got: %v", test.expect, got)
		}
	}
}

func TestInterpretBinaryError(t *testing.T) {
	expr := &token.Binary{
		Left:      &token.Literal{Value: 16.0},
		Operation: scanner.Token{TokenType: scanner.STAR},
		Right:     &token.Literal{Value: "muffin"},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.Interpret(expr)
	if err != nil {
		t.Fatalf("expect nil")
	}
	if len(errs) != 1 {
		t.Fatalf("expect len(%d) got %d", 1, len(errs))
	}
	if got != "" {
		t.Fatalf("expect empty")
	}
	if errs[0] != "Operands must be a numbers." {
		t.Fatalf("expect %v got %v", "Operands must be a numbers.", err.Error())
	}
}

func TestInterpretBinaryDivisionByZero(t *testing.T) {
	t.Skip()
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.Binary{
			Left:      &token.Literal{Value: 256.0},
			Operation: scanner.Token{TokenType: scanner.SLASH},
			Right:     &token.Literal{Value: 0.0},
		}, nil},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	for _, test := range tests {
		got, err := i.Interpret(test.expr)
		if err != nil {
			t.Error(err)
		}
		if got != test.expect {
			t.Fatalf("expect: %v got: %v", test.expect, got)
		}
	}
}

func TestIsTruthy(t *testing.T) {
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	if i.isTruthy(nil) {
		t.Fatalf("expected false")
	}
}

var testCallBack = func(errs *[]string) ErrorCallback {
	return func(err *RuntimeError) {
		*errs = append(*errs, err.Error())
	}
}
