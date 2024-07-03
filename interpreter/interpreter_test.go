package interpreter

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/scanner/testutil"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestInterpretLiteralExprFloat(t *testing.T) {
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

func TestInterpretLiteralExprString(t *testing.T) {
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

func TestInterpretLiteralExprBoolean(t *testing.T) {
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

func TestInterpretUnaryExprMinus(t *testing.T) {
	expr := &token.UnaryExpr{
		Operator: testutil.Minus(),
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

func TestInterpretUnaryExprError(t *testing.T) {
	expr := &token.UnaryExpr{
		Operator: testutil.Minus(),
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

func TestInterpretUnaryExprBang(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.UnaryExpr{
			Operator: testutil.Bang(),
			Right:    &token.LiteralExpr{Value: 123.0},
		}, false},
		{&token.UnaryExpr{
			Operator: testutil.Bang(),
			Right:    &token.LiteralExpr{Value: "abc"},
		}, false},
		{&token.UnaryExpr{
			Operator: testutil.Bang(),
			Right:    &token.LiteralExpr{Value: true},
		}, false},
		{&token.UnaryExpr{
			Operator: testutil.Bang(),
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

func TestInterpretBinaryExpr(t *testing.T) {
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 40.02},
			Operator: testutil.Plus(),
			Right:    &token.LiteralExpr{Value: 2.40},
		}, 42.42},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "Hello, "},
			Operator: testutil.Plus(),
			Right:    &token.LiteralExpr{Value: "World!"},
		}, "Hello, World!"},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 16.0},
			Operator: testutil.Star(),
			Right:    &token.LiteralExpr{Value: 2.0},
		}, 32.0},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 256.0},
			Operator: testutil.Slash(),
			Right:    &token.LiteralExpr{Value: 2.0},
		}, 128.0},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.0},
			Operator: testutil.Greater(),
			Right:    &token.LiteralExpr{Value: 1.0},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: testutil.GreaterEqual(),
			Right:    &token.LiteralExpr{Value: 2.4},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.0},
			Operator: testutil.Less(),
			Right:    &token.LiteralExpr{Value: 1.0},
		}, false},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: testutil.LessEqual(),
			Right:    &token.LiteralExpr{Value: 2.4},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: testutil.EqualEqual(),
			Right:    &token.LiteralExpr{Value: 2.4},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "lox"},
			Operator: testutil.EqualEqual(),
			Right:    &token.LiteralExpr{Value: "lox"},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: testutil.BangEqual(),
			Right:    &token.LiteralExpr{Value: 2.4},
		}, false},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "lox"},
			Operator: testutil.BangEqual(),
			Right:    &token.LiteralExpr{Value: "lox"},
		}, false},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: "lox"},
			Operator: testutil.BangEqual(),
			Right:    &token.LiteralExpr{Value: nil},
		}, true},
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 2.4},
			Operator: testutil.BangEqual(),
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

func TestInterpretBinaryExprError(t *testing.T) {
	expr := &token.BinaryExpr{
		Left:     &token.LiteralExpr{Value: 16.0},
		Operator: testutil.Star(),
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

func TestDeclVarStmt(t *testing.T) {
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

func TestInterpretAssignExpr(t *testing.T) {
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

func TestIfStmt(t *testing.T) {
	errs := make([]string, 0)
	i := New(testCallBack(&errs))
	ifStmt := &token.IfStmt{
		Condition:  &token.LiteralExpr{Value: 1.0},
		ThenBranch: &token.PrintStmt{Expression: &token.LiteralExpr{Value: "is true"}},
		ElseBranch: &token.PrintStmt{Expression: &token.LiteralExpr{Value: "is false"}},
	}
	got, err := i.exec(ifStmt)
	if got != nil {
		t.Fatalf("expect nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
}

func TestBlockStmt(t *testing.T) {
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

func TestInterpretBinaryExprDivisionByZero(t *testing.T) {
	t.Skip()
	tests := []struct {
		expr   token.Expr
		expect interface{}
	}{
		{&token.BinaryExpr{
			Left:     &token.LiteralExpr{Value: 256.0},
			Operator: testutil.Slash(),
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
