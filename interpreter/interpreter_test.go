package interpreter

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/scanner/testutil"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestInterpretLiteralFloat(t *testing.T) {
	expr := &token.LiteralExpr{Value: 123.}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.eval(expr)
	if err != nil {
		t.Error(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}
	if got != 123.0 {
		t.Fatalf("expected '123' got '%v'", got)
	}
}

func TestInterpretLiteralString(t *testing.T) {
	expr := &token.LiteralExpr{Value: "abc"}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.eval(expr)
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
	expr := &token.LiteralExpr{Value: true}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.eval(expr)
	if err != nil {
		t.Error(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}
	if got != true {
		t.Fatalf("expected 'true' got '%v'", got)
	}
}

func TestInterpretUnaryMinus(t *testing.T) {
	expr := &token.UnaryExpr{
		Operator: scanner.Token{TokenType: scanner.MINUS},
		Right:    &token.LiteralExpr{Value: 123.0},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.eval(expr)
	if err != nil {
		t.Error(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expect len(1), got: %d", len(errs))
	}
	if got != -123.0 {
		t.Fatalf("expect: '-123' got: '%v'", got)
	}
}

func TestInterpretUnaryError(t *testing.T) {
	expr := &token.UnaryExpr{
		Operator: scanner.Token{TokenType: scanner.MINUS},
		Right:    &token.LiteralExpr{Value: "muffin"},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.eval(expr)
	if err == nil {
		t.Fatalf("expect not nil")
	}
	if err.Error() != "Operand must be a number." {
		t.Fatalf("expect 'Operand must be a number.', got '%v'", err.Error())
	}
	if len(errs) != 0 {
		t.Fatalf("expect empty, got: %d", len(errs))
	}
	if got == "" {
		t.Fatalf("expect not empty")
	}
}

func TestInterpretUnaryBang(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.UnaryExpr{
			Operator: scanner.Token{TokenType: scanner.BANG},
			Right:    &token.LiteralExpr{Value: 123.0},
		}, false},
		{&token.UnaryExpr{
			Operator: scanner.Token{TokenType: scanner.BANG},
			Right:    &token.LiteralExpr{Value: "abc"},
		}, false},
		{&token.UnaryExpr{
			Operator: scanner.Token{TokenType: scanner.BANG},
			Right:    &token.LiteralExpr{Value: true},
		}, false},
		{&token.UnaryExpr{
			Operator: scanner.Token{TokenType: scanner.BANG},
			Right:    &token.LiteralExpr{Value: false},
		}, true},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	for _, test := range tests {
		got, err := i.eval(test.expr)
		if err != nil {
			t.Error(err)
		}
		if len(errs) != 0 {
			t.Fatalf("expect len(1), got: %d", len(errs))
		}
		if got != test.expect {
			t.Fatalf("expect: '%v' got: '%v'", test.expect, got)
		}
	}
}

func TestInterpretBinary(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 40.02},
			Operator: scanner.Token{TokenType: scanner.PLUS},
			Right:    &token.LiteralExpr{Value: 2.40},
		}, 42.42},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "Hello, "},
			Operator: scanner.Token{TokenType: scanner.PLUS},
			Right:    &token.LiteralExpr{Value: "World!"},
		}, "Hello, World!"},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 16.0},
			Operator: scanner.Token{TokenType: scanner.STAR},
			Right:    &token.LiteralExpr{Value: 2.0},
		}, 32.0},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 256.0},
			Operator: scanner.Token{TokenType: scanner.SLASH},
			Right:    &token.LiteralExpr{Value: 2.0},
		}, 128.0},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.0},
			Operator: scanner.Token{TokenType: scanner.GREATER},
			Right:    &token.LiteralExpr{Value: 1.0},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: scanner.Token{TokenType: scanner.GREATER_EQUAL},
			Right:    &token.LiteralExpr{Value: 2.4},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.0},
			Operator: scanner.Token{TokenType: scanner.LESS},
			Right:    &token.LiteralExpr{Value: 1.0},
		}, false},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: scanner.Token{TokenType: scanner.LESS_EQUAL},
			Right:    &token.LiteralExpr{Value: 2.4},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: scanner.Token{TokenType: scanner.EQUAL_EQUAL},
			Right:    &token.LiteralExpr{Value: 2.4},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "lox"},
			Operator: scanner.Token{TokenType: scanner.EQUAL_EQUAL},
			Right:    &token.LiteralExpr{Value: "lox"},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:    &token.LiteralExpr{Value: 2.4},
		}, false},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "lox"},
			Operator: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:    &token.LiteralExpr{Value: "lox"},
		}, false},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "lox"},
			Operator: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:    &token.LiteralExpr{Value: nil},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: scanner.Token{TokenType: scanner.BANG_EQUAL},
			Right:    &token.LiteralExpr{Value: "lox"},
		}, true},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	for _, test := range tests {
		got, err := i.eval(test.expr)
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
	expr := &token.BinaryExpr{
		Left:     &token.LiteralExpr{Value: 16.0},
		Operator: scanner.Token{TokenType: scanner.STAR},
		Right:    &token.LiteralExpr{Value: "muffin"},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.eval(expr)
	if err == nil {
		t.Fatalf("expect not nil")
	}
	if err.Error() != "Operands must be a numbers." {
		t.Fatalf("expect 'Operands must be a numbers.', got %v", err.Error())
	}
	if len(errs) != 0 {
		t.Fatalf("expect empty got %d", len(errs))
	}
	if got != nil {
		t.Fatalf("expect empty")
	}
}

func TestDeclVar(t *testing.T) {
	tok := testutil.Identifier("a")
	stmt := &token.VarStmt{
		Name: tok,
		Initializer: &token.LiteralExpr{
			Value: testutil.Str("before"),
		},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	got, err := i.exec(stmt)
	if got != nil {
		t.Fatalf("expect nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	assertToken(t, i, &tok, "before")
}

func TestInterpretAssign(t *testing.T) {
	errs := make([]string, 0)
	i := New(testCallBack(&errs))

	tok := testutil.Identifier("a")
	declStmt := &token.VarStmt{
		Name: tok,
		Initializer: &token.LiteralExpr{
			Value: testutil.Str("before"),
		},
	}
	if _, err := i.exec(declStmt); err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	assertToken(t, i, &tok, "before")

	assign := &token.AssignExpr{
		Name: tok,
		Value: &token.LiteralExpr{
			Value: testutil.Str("after"),
		},
	}
	got, err := i.eval(assign)
	if got == nil {
		t.Fatalf("expect not nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	gotLiteral, ok := got.(scanner.Token)
	if !ok {
		t.Fatalf("expect *token.Literal got %T", got)
	}
	if gotLiteral.Literal != "after" {
		t.Fatalf("expect 'after' got %v", gotLiteral.Literal)
	}
	assertToken(t, i, &tok, "after")
}

func TestBlock(t *testing.T) {
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	tokA, tokB, tokC := testutil.Identifier("a"), testutil.Identifier("b"), testutil.Identifier("c")
	block := &token.BlockStmt{Statements: []token.Stmt{
		&token.VarStmt{Name: tokA, Initializer: &token.LiteralExpr{Value: testutil.Str("global a")}},
		&token.VarStmt{Name: tokB, Initializer: &token.LiteralExpr{Value: testutil.Str("global b")}},
		&token.VarStmt{Name: tokC, Initializer: &token.LiteralExpr{Value: testutil.Str("global c")}},
		&token.BlockStmt{Statements: []token.Stmt{
			&token.VarStmt{Name: tokA, Initializer: &token.LiteralExpr{Value: testutil.Str("outer a")}},
			&token.VarStmt{Name: tokB, Initializer: &token.LiteralExpr{Value: testutil.Str("outer b")}},
			&token.BlockStmt{Statements: []token.Stmt{
				&token.VarStmt{Name: tokA, Initializer: &token.LiteralExpr{Value: testutil.Str("inner a")}},
				&token.PrintStmt{Expression: &token.VariableExpr{Name: tokA}},
				&token.PrintStmt{Expression: &token.VariableExpr{Name: tokB}},
				&token.PrintStmt{Expression: &token.VariableExpr{Name: tokC}},
			}},
			&token.PrintStmt{Expression: &token.VariableExpr{Name: tokA}},
			&token.PrintStmt{Expression: &token.VariableExpr{Name: tokB}},
			&token.PrintStmt{Expression: &token.VariableExpr{Name: tokC}},
		}},
		&token.PrintStmt{Expression: &token.VariableExpr{Name: tokA}},
		&token.PrintStmt{Expression: &token.VariableExpr{Name: tokB}},
		&token.PrintStmt{Expression: &token.VariableExpr{Name: tokC}},
	}}
	got, err := i.exec(block)
	if got != nil {
		t.Fatalf("expect nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
}

func TestInterpretBinaryDivisionByZero(t *testing.T) {
	t.Skip()
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 256.0},
			Operator: scanner.Token{TokenType: scanner.SLASH},
			Right:    &token.LiteralExpr{Value: 0.0},
		}, nil},
	}
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	for _, test := range tests {
		got, err := i.eval(test.expr)
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

func assertToken(t *testing.T, i *Interpreter, tok *scanner.Token, expect interface{}) {
	value, err := i.env.Get(tok)
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	gotToken, ok := value.(scanner.Token)
	if !ok {
		t.Fatalf("expect *token.Literal got %T", gotToken)
	}
	if gotToken.Literal != expect {
		t.Fatalf("expect '%v' got %v", expect, gotToken.Literal)
	}
}

var testCallBack = func(errs *[]string) ErrorCallback {
	return func(err *RuntimeError) {
		*errs = append(*errs, err.Error())
	}
}
