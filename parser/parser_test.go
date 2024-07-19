package parser

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/scanner/testutil"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestParseDeclarationStmt(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.VarDecl(),
			testutil.Identifier("a"),
			testutil.Equal(),
			testutil.Str("before"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
	got, ok := stmts[0].(*token.VarStmt)
	if !ok {
		t.Fatalf("expect *token.VarStmt got %T", got)
	}
	if got.Name.Lexeme == nil || *got.Name.Lexeme != "a" {
		t.Fatalf("expect 'a' got '%v'", got.Name.Lexeme)
	}
	init, ok := got.Initializer.(*token.LiteralExpr)
	if !ok {
		t.Fatalf("expect *token.LiteralExpr got %T", init)
	}
	if init.Value != "before" {
		t.Fatalf("expect 'before' got '%v'", init.Value)
	}
}

func TestParseAssignExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Identifier("a"),
			testutil.Equal(),
			testutil.Str("after"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
	got, ok := stmts[0].(*token.ExpressionStmt)
	if !ok {
		t.Fatalf("expect *token.Expression got %T", got)
	}
	expr, ok := got.Expression.(*token.AssignExpr)
	if !ok {
		t.Fatalf("expect *token.Assign got %T", got)
	}
	if expr.Name.Lexeme == nil || *expr.Name.Lexeme != "a" {
		t.Fatalf("expect 'a' got '%v'", expr.Name.Lexeme)
	}
	value, ok := expr.Value.(*token.LiteralExpr)
	if !ok {
		t.Fatalf("expect *token.LiteralExpr got %T", got)
	}
	if value == nil || value.Value != "after" {
		t.Fatalf("expect 'after' got '%v'", value.Value)
	}
}

func TestParseAssignExprError(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Identifier("a"),
			testutil.Plus(),
			testutil.Identifier("b"),
			testutil.Equal(),
			testutil.Number(10),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateHasErrors(t, stmts, errors, err, "Invalid assignment target.")
}

func TestParseLiteralExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Number(123.0),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
}

func TestParseUnaryExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Minus(),
			testutil.Number(123.0),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
}

func TestParseTermExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Number(321.0),
			testutil.Plus(),
			testutil.Number(123.0),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
}

func TestParseFactorialExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Number(321.0),
			testutil.Star(),
			testutil.Number(123.0),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
}

func TestParseGroupingExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.LeftParen(),
			testutil.Number(321.0),
			testutil.Star(),
			testutil.Number(123.0),
			testutil.RightParen(),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
}

func TestParseGroupingExprError(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.LeftParen(),
			testutil.Number(321.0),
			testutil.Minus(),
			testutil.Number(123.0),
			testutil.LeftParen(),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateHasErrors(t, stmts, errors, err, "expect expression")
}

func TestParseIfStmt(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.IfIdentifier(),
			testutil.LeftParen(),
			testutil.Number(1.0), // any non-nil value is truthy
			testutil.RightParen(),
			testutil.Print(),
			testutil.Str("is true"),
			testutil.Semicolon(),
			testutil.ElseIdentifier(),
			testutil.Str("is false"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)

	got, ok := stmts[0].(*token.IfStmt)
	if !ok {
		t.Fatalf("expect *token.IfStmt got %T", got)
	}
}

func TestParseIfStmtError(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.IfIdentifier(),
			testutil.LeftParen(),
			testutil.Number(1.0), // any non-nil value is truthy
			// missing ')'
			testutil.Print(),
			testutil.Str("is true"),
			testutil.Semicolon(),
			testutil.ElseIdentifier(),
			testutil.Str("is false"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateHasErrors(t, stmts, errors, err, "expect ')' after if condition")
}

func TestParseBlockStmt(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.LeftBrace(),
			testutil.Number(321.0),
			testutil.Star(),
			testutil.Number(123.0),
			testutil.Semicolon(),
			testutil.RightBrace(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
	got, ok := stmts[0].(*token.BlockStmt)
	if !ok {
		t.Fatalf("expect *token.Block got %T", got)
	}
	if len(got.Statements) != 1 {
		t.Fatalf("expect len(1) got: %d", len(got.Statements))
	}
	gotStmt := got.Statements[0]
	expr, ok := gotStmt.(*token.ExpressionStmt)
	if !ok {
		t.Fatalf("expect *token.Expression got %T", gotStmt)
	}
	_, ok = expr.Expression.(*token.BinaryExpr)
	if !ok {
		t.Fatalf("expect *token.Expression got %T", expr.Expression)
	}
}

func TestParseBlockStmtError(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.LeftBrace(),
			testutil.Number(321.0),
			testutil.Star(),
			testutil.Number(123.0),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateHasErrors(t, stmts, errors, err, "Expect '}' after block.")
}

func TestParseLogicalOrExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Str("5"),
			testutil.Or(),
			testutil.Str("4"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
	got, ok := stmts[0].(*token.ExpressionStmt)
	if !ok {
		t.Fatalf("expect *token.ExpressionStmt got %T", got)
	}
	expr, ok := got.Expression.(*token.LogicalExpr)
	if !ok {
		t.Fatalf("expect *token.LogicalExpr got %T", expr)
	}
	if expr.Operator.TokenType != scanner.OR {
		t.Fatalf("expect scanner.OR got %T", expr.Operator.TokenType)
	}
}

func TestParseLogicalAndExpr(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.Str("5"),
			testutil.And(),
			testutil.Str("4"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
	got, ok := stmts[0].(*token.ExpressionStmt)
	if !ok {
		t.Fatalf("expect *token.ExpressionStmt got %T", got)
	}
	expr, ok := got.Expression.(*token.LogicalExpr)
	if !ok {
		t.Fatalf("expect *token.LogicalExpr got %T", expr)
	}
	if expr.Operator.TokenType != scanner.AND {
		t.Fatalf("expect scanner.AND got %T", expr.Operator.TokenType)
	}
}

func TestWhileStmt(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.While(),
			testutil.LeftParen(),
			testutil.Str("true"),
			testutil.RightParen(),
			testutil.Print(),
			testutil.Str("Hurray!"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
	got, ok := stmts[0].(*token.WhileStmt)
	if !ok {
		t.Fatalf("expect *token.WhileStmt got %T", got)
	}
}

func TestForStmt(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.For(),

			testutil.LeftParen(),
			// initialization
			testutil.VarDecl(),
			testutil.Identifier("i"),
			testutil.Equal(),
			testutil.Number(1.0),
			testutil.Semicolon(),
			// condition
			testutil.Identifier("i"),
			testutil.Less(),
			testutil.Number(3.0),
			testutil.Semicolon(),
			// increment
			testutil.Identifier("i"),
			testutil.Equal(),
			testutil.Identifier("i"),
			testutil.Plus(),
			testutil.Number(1.0),
			testutil.RightParen(),

			testutil.Print(),
			testutil.Str("Hurray!"),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	validateNoError(t, stmts, errors, err)
	stmt, ok := stmts[0].(*token.BlockStmt)
	if !ok {
		t.Fatalf("expect *token.WhileStmt got %T", stmts[0])
	}
	_, ok = stmt.Statements[0].(*token.VarStmt)
	if !ok {
		t.Fatalf("expect *token.VarStmt got %T", stmt.Statements[0])
	}
	_, ok = stmt.Statements[1].(*token.WhileStmt)
	if !ok {
		t.Fatalf("expect *token.WhileStmt got %T", stmt.Statements[1])
	}
}

var testCallBack = func(errs *[]string) ErrorCallback {
	return func(token scanner.Token, message string) {
		*errs = append(*errs, message)
	}
}

func validateNoError(t *testing.T, stmts []token.Stmt, errors []string, err error) {
	if err != nil {
		t.Fatalf("expect: nil got: %v", err)
	}
	if len(errors) != 0 {
		t.Log(errors)
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if stmts == nil {
		t.Fatal("expect not nil")
	}
	if len(stmts) != 1 {
		t.Fatalf("expect len(1) got %v", len(stmts))
	}
}

func validateHasErrors(t *testing.T, stmts []token.Stmt, errors []string, err error, expectErrors ...string) {
	if len(stmts) != 0 {
		t.Fatalf("expect empty got %v", len(stmts))
	}
	if err != nil {
		t.Fatalf("expect: nil got: %v", err)
	}
	if len(errors) != 1 {
		t.Log(errors)
		t.Fatalf("expect not empty")
	}
	if len(errors) != len(expectErrors) {
		t.Fatalf("expect len(%d) got len(%d)", len(expectErrors), len(errors))
	}
	for i := range errors {
		if errors[i] != expectErrors[i] {
			t.Fatalf("expect '%v' got '%v'", expectErrors[i], errors[i])
		}
	}
}
