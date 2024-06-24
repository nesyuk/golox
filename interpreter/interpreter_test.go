package interpreter

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestInterpretLiteralFloat(t *testing.T) {
	expr := &token.Literal{Value: 123.}
	i := &Interpreter{}
	got := i.Eval(expr)
	result, ok := got.(float64)
	if !ok {
		t.Fatalf("expected float64 got %T", got)
	}
	if result != 123. {
		t.Fatalf("expected 123. got %f", result)
	}
}

func TestInterpretLiteralString(t *testing.T) {
	expr := &token.Literal{Value: "abc"}
	i := &Interpreter{}
	got := i.Eval(expr)
	result, ok := got.(string)
	if !ok {
		t.Fatalf("expected string got %T", got)
	}
	if result != "abc" {
		t.Fatalf("expected abc got %v", result)
	}
}

func TestInterpretLiteralBoolean(t *testing.T) {
	expr := &token.Literal{Value: true}
	i := &Interpreter{}
	got := i.Eval(expr)
	result, ok := got.(bool)
	if !ok {
		t.Fatalf("expected bool got %T", got)
	}
	if result != true {
		t.Fatalf("expected true got %v", result)
	}
}

func TestInterpretUnaryMinus(t *testing.T) {
	expr := &token.Unary{
		Operation: scanner.Token{TokenType: scanner.MINUS},
		Right:     &token.Literal{Value: 123.0},
	}
	i := &Interpreter{}
	got := i.Eval(expr)
	if got != -123.0 {
		t.Fatalf("expect: %v got: %v", -123.0, got)
	}
}

func TestInterpretUnaryBang(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect bool
	}{
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: 123.0},
		}, false},
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: "abc"},
		}, false},
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: true},
		}, false},
		{&token.Unary{
			Operation: scanner.Token{TokenType: scanner.BANG},
			Right:     &token.Literal{Value: false},
		}, true},
	}
	i := &Interpreter{}
	for _, test := range tests {
		got := i.Eval(test.expr)
		if got != test.expect {
			t.Fatalf("expect: %v got: %v", test.expect, got)
		}
	}
}

func TestInterpretBinary(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.Binary{
			Left:      &token.Literal{Value: 40.02},
			Operation: scanner.Token{TokenType: scanner.PLUS},
			Right:     &token.Literal{Value: 2.40},
		}, 42.42},
		{&token.Binary{
			Left:      &token.Literal{Value: "Hello, "},
			Operation: scanner.Token{TokenType: scanner.PLUS},
			Right:     &token.Literal{Value: "World!"},
		}, "Hello, World!"},
		{&token.Binary{
			Left:      &token.Literal{Value: 16.0},
			Operation: scanner.Token{TokenType: scanner.STAR},
			Right:     &token.Literal{Value: 2.0},
		}, 32.0},
		{&token.Binary{
			Left:      &token.Literal{Value: 256.0},
			Operation: scanner.Token{TokenType: scanner.SLASH},
			Right:     &token.Literal{Value: 2.0},
		}, 128.0},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.0},
			Operation: scanner.Token{TokenType: scanner.GREATER},
			Right:     &token.Literal{Value: 1.0},
		}, true},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.GREATER_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, true},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.0},
			Operation: scanner.Token{TokenType: scanner.LESS},
			Right:     &token.Literal{Value: 1.0},
		}, false},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.LESS_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, true},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.EQUAL_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, true},
		{&token.Binary{
			Left:      &token.Literal{Value: "lox"},
			Operation: scanner.Token{TokenType: scanner.EQUAL_EQUAL},
			Right:     &token.Literal{Value: "lox"},
		}, true},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: 2.4},
		}, false},
		{&token.Binary{
			Left:      &token.Literal{Value: "lox"},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: "lox"},
		}, false},
		{&token.Binary{
			Left:      &token.Literal{Value: "lox"},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: nil},
		}, true},
		{&token.Binary{
			Left:      &token.Literal{Value: 2.4},
			Operation: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:     &token.Literal{Value: "lox"},
		}, true},
	}
	i := &Interpreter{}
	for _, test := range tests {
		got := i.Eval(test.expr)
		if got != test.expect {
			t.Fatalf("expect: %v got: %v", test.expect, got)
		}
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
	i := &Interpreter{}
	for _, test := range tests {
		got := i.Eval(test.expr)
		if got != test.expect {
			t.Fatalf("expect: %v got: %v", test.expect, got)
		}
	}
}

func TestIsTruthy(t *testing.T) {
	i := &Interpreter{}
	if i.isTruthy(nil) {
		t.Fatalf("expected false")
	}
}
