package interpreter

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/scanner/testutil"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestInterpretLiteralExprFloat(t *testing.T) {
	expr := &token.LiteralExpr{Value: 123.}
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	got, err := i.eval(expr)
	if err != nil {
		t.Error(err)
	}
	v.validateNoErrors(t)
	if got != 123.0 {
		t.Fatalf("expected '123' got '%v'", got)
	}
}

func TestInterpretLiteralExprString(t *testing.T) {
	expr := &token.LiteralExpr{Value: "abc"}
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	got, err := i.eval(expr)
	if err != nil {
		t.Error(err)
	}
	v.validateNoErrors(t)

	if got != "abc" {
		t.Fatalf("expected 'abc' got '%v'", got)
	}
}

func TestInterpretLiteralExprBoolean(t *testing.T) {
	expr := &token.LiteralExpr{Value: true}
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	got, err := i.eval(expr)
	if err != nil {
		t.Error(err)
	}
	v.validateNoErrors(t)
	if got != true {
		t.Fatalf("expected 'true' got '%v'", got)
	}
}

func TestInterpretUnaryExprMinus(t *testing.T) {
	expr := &token.UnaryExpr{
		Operator: testutil.Minus(),
		Right:    &token.LiteralExpr{Value: 123.0},
	}
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	got, err := i.eval(expr)
	if err != nil {
		t.Error(err)
	}
	v.validateNoErrors(t)
	if got != -123.0 {
		t.Fatalf("expect: '-123' got: '%v'", got)
	}
}

func TestInterpretUnaryExprError(t *testing.T) {
	expr := &token.UnaryExpr{
		Operator: testutil.Minus(),
		Right:    &token.LiteralExpr{Value: "muffin"},
	}
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	got, err := i.eval(expr)
	if err == nil {
		t.Fatalf("expect not nil")
	}
	if err.Error() != "Operand must be a number." {
		t.Fatalf("expect 'Operand must be a number.', got '%v'", err.Error())
	}
	v.validateNoErrors(t)
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
	for _, test := range tests {
		v := NewValidator()
		i := New(v.onError, v.onPrint)
		got, err := i.eval(test.expr)
		if err != nil {
			t.Error(err)
		}
		v.validateNoErrors(t)
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

	for _, test := range tests {
		v := NewValidator()
		i := New(v.onError, v.onPrint)
		got, err := i.eval(test.expr)
		if err != nil {
			t.Error(err)
		}
		v.validateNoErrors(t)
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
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	got, err := i.eval(expr)
	if err == nil {
		t.Fatalf("expect not nil")
	}
	if err.Error() != "Operands must be a numbers." {
		t.Fatalf("expect 'Operands must be a numbers.', got %v", err.Error())
	}
	v.validateNoErrors(t)
	if got != nil {
		t.Fatalf("expect empty")
	}
}

func TestDeclVarStmt(t *testing.T) {
	tok := testutil.Identifier("a")
	stmt := &token.VarStmt{
		Name: tok,
		Initializer: &token.LiteralExpr{
			Value: "before",
		},
	}
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	got, err := i.exec(stmt)
	if got != nil {
		t.Fatalf("expect nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	assertValue(t, i, &tok, "before")
}

func TestInterpretAssignExpr(t *testing.T) {
	v := NewValidator()
	i := New(v.onError, v.onPrint)

	tok := testutil.Identifier("a")
	declStmt := &token.VarStmt{
		Name: tok,
		Initializer: &token.LiteralExpr{
			Value: "before",
		},
	}
	if _, err := i.exec(declStmt); err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	assertValue(t, i, &tok, "before")

	assign := &token.AssignExpr{
		Name: tok,
		Value: &token.LiteralExpr{
			Value: "after",
		},
	}
	got, err := i.eval(assign)
	if got == nil {
		t.Fatalf("expect not nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	assertValue(t, i, &tok, "after")
}

func TestIfStmt(t *testing.T) {
	v := NewValidator()
	i := New(v.onError, v.onPrint)
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

func TestWhileStmt(t *testing.T) {
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	tok := testutil.Identifier("i")
	variable := token.LiteralExpr{Value: 1.0}
	varExpr := token.VariableExpr{Name: tok}
	i.locals = map[token.Expr]int{
		&varExpr: 0,
	}
	stmt := &token.BlockStmt{
		Statements: []token.Stmt{
			&token.VarStmt{Name: tok, Initializer: &variable},
			&token.WhileStmt{
				Condition: &token.BinaryExpr{Left: &varExpr, Operator: testutil.Less(), Right: &token.LiteralExpr{Value: 1.0}},
				Body: &token.BlockStmt{
					Statements: []token.Stmt{
						&token.PrintStmt{Expression: &varExpr},
						&token.ExpressionStmt{
							Expression: &token.BinaryExpr{
								Left:     &varExpr,
								Operator: testutil.Plus(),
								Right:    &token.LiteralExpr{Value: 1.0}}},
					},
				},
			},
		},
	}
	got, err := i.exec(stmt)
	if got != nil {
		t.Fatalf("expect nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
}

func TestLogicalOrExpr(t *testing.T) {
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	tests := []struct {
		expr   *token.LogicalExpr
		expect interface{}
	}{
		{&token.LogicalExpr{
			Left:     &token.LiteralExpr{Value: "hi"},
			Operator: testutil.Or(),
			Right:    &token.LiteralExpr{Value: 2.0},
		}, "hi"},
		{&token.LogicalExpr{
			Left:     &token.LiteralExpr{Value: nil},
			Operator: testutil.Or(),
			Right:    &token.LiteralExpr{Value: "yes"},
		}, "yes"},
	}
	for _, test := range tests {
		got, err := i.eval(test.expr)
		if got == nil {
			t.Fatalf("expect not empty")
		}
		if got != test.expect {
			t.Fatalf("expect '%v' got '%v'", test.expect, got)
		}
		if err != nil {
			t.Fatalf("expect nil got '%v'", err)
		}
	}
}

func TestBlockStmt(t *testing.T) {
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	tokA, tokB, tokC := testutil.Identifier("a"), testutil.Identifier("b"), testutil.Identifier("c")
	block := &token.BlockStmt{Statements: []token.Stmt{
		&token.VarStmt{Name: tokA, Initializer: &token.LiteralExpr{Value: "global a"}},
		&token.VarStmt{Name: tokB, Initializer: &token.LiteralExpr{Value: "global b"}},
		&token.VarStmt{Name: tokC, Initializer: &token.LiteralExpr{Value: "global c"}},
		&token.BlockStmt{Statements: []token.Stmt{
			&token.VarStmt{Name: tokA, Initializer: &token.LiteralExpr{Value: "outer a"}},
			&token.VarStmt{Name: tokB, Initializer: &token.LiteralExpr{Value: "outer b"}},
			&token.BlockStmt{Statements: []token.Stmt{
				&token.VarStmt{Name: tokA, Initializer: &token.LiteralExpr{Value: "inner a"}},
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

func TestVariableInEqualityExpr(t *testing.T) {
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	tokI := testutil.Identifier("i")
	varExpr := token.VariableExpr{Name: tokI}
	assignExpr := token.AssignExpr{
		Name: tokI,
		Value: &token.BinaryExpr{
			Left:     &varExpr,
			Operator: testutil.Less(),
			Right:    &token.LiteralExpr{Value: 2.0},
		}}
	i.locals = map[token.Expr]int{
		&varExpr:    0,
		&assignExpr: 1,
	}
	block := &token.BlockStmt{Statements: []token.Stmt{
		&token.VarStmt{Name: tokI, Initializer: &token.LiteralExpr{Value: 1.0}},
		&token.ExpressionStmt{Expression: &assignExpr},
		&token.PrintStmt{Expression: &varExpr},
	}}
	got, err := i.exec(block)
	if got != nil {
		t.Fatalf("expect nil")
	}
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
}

func TestVariableInBinaryExpr(t *testing.T) {
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	tokI := testutil.Identifier("i")
	varExpr := token.VariableExpr{Name: tokI}
	assignExpr := token.AssignExpr{
		Name: tokI,
		Value: &token.BinaryExpr{
			Left:     &varExpr,
			Operator: testutil.Plus(),
			Right:    &token.LiteralExpr{Value: 1.0},
		}}
	i.locals = map[token.Expr]int{
		&varExpr:    0,
		&assignExpr: 1,
	}
	block := &token.BlockStmt{Statements: []token.Stmt{
		&token.VarStmt{Name: tokI, Initializer: &token.LiteralExpr{Value: 1.0}},
		&token.ExpressionStmt{Expression: &assignExpr},
		&token.PrintStmt{Expression: &varExpr},
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

	for _, test := range tests {
		v := NewValidator()
		i := New(v.onError, v.onPrint)
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
	v := NewValidator()
	i := New(v.onError, v.onPrint)
	if i.isTruthy(nil) {
		t.Fatalf("expected false")
	}
}

func assertValue(t *testing.T, i *Interpreter, tok *scanner.Token, expect interface{}) {
	value, err := i.env.Get(tok)
	if err != nil {
		t.Fatalf("expect nil got %v", err)
	}
	if value != expect {
		t.Fatalf("expect '%v' got %v", expect, value)
	}
}

type validator struct {
	errors  []string
	results []string
}

func NewValidator() *validator {
	return &validator{make([]string, 0), make([]string, 0)}
}

func (v *validator) onError(err *RuntimeError) {
	v.errors = append(v.errors, err.Error())
}

func (v *validator) onPrint(s string) {
	v.results = append(v.results, s)
}

func (v *validator) validateNoErrors(t *testing.T) {
	v.validateErrors(t, make([]string, 0))
}

func (v *validator) validateErrors(t *testing.T, expect []string) {
	v.validate(t, "errors", v.errors, expect)
}

func (v *validator) validateNoResult(t *testing.T) {
	v.validateResult(t, make([]string, 0))
}

func (v *validator) validateResult(t *testing.T, expect []string) {
	v.validate(t, "results", v.results, expect)
}

func (v *validator) validate(t *testing.T, prefix string, got, expect []string) {
	if len(got) != len(expect) {
		t.Errorf("%s: expect %v, got: %v", prefix, len(expect), len(got))
		t.Log(got)
		return
	}
	for i := range got {
		if got[i] != expect[i] {
			t.Errorf("%s expect: %v, got: %v", prefix, expect[i], got)
		}
	}
}

func numberIdentifier(name string, value float64) *token.VarStmt {
	tok := testutil.Identifier(name)
	return &token.VarStmt{Name: tok, Initializer: &token.LiteralExpr{Value: testutil.Number(value)}}
}
