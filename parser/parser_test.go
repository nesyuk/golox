package parser

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/scanner/testutil"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestParseDeclaration(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if stmts == nil {
		t.Fatal("expect not nil")
	}
	if len(stmts) != 1 {
		t.Fatalf("expect len(1) got %v", len(stmts))
	}
	got, ok := stmts[0].(*token.Var)
	if !ok {
		t.Fatalf("expect *token.Var got %T", got)
	}
	if got.Name.Lexeme == nil || *got.Name.Lexeme != "a" {
		t.Fatalf("expect 'a' got '%v'", got.Name.Lexeme)
	}
	init, ok := got.Initializer.(*token.Literal)
	if !ok {
		t.Fatalf("expect *token.Literal got %T", init)
	}
	if init.Value != "before" {
		t.Fatalf("expect 'before' got '%v'", init.Value)
	}
}

func TestParseAssign(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if stmts == nil {
		t.Fatal("expect not nil")
	}
	if len(stmts) != 1 {
		t.Fatalf("expect len(1) got %v", len(stmts))
	}
	got, ok := stmts[0].(*token.Expression)
	if !ok {
		t.Fatalf("expect *token.Expression got %T", got)
	}
	expr, ok := got.Expression.(*token.Assign)
	if !ok {
		t.Fatalf("expect *token.Assign got %T", got)
	}
	if expr.Name.Lexeme == nil || *expr.Name.Lexeme != "a" {
		t.Fatalf("expect 'a' got '%v'", expr.Name.Lexeme)
	}
	value, ok := expr.Value.(*token.Literal)
	if !ok {
		t.Fatalf("expect *token.Literal got %T", got)
	}
	if value == nil || value.Value != "after" {
		t.Fatalf("expect 'after' got '%v'", value.Value)
	}
}

func TestParseAssignError(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) > 0 {
		t.Fatalf("expect empty, got %v", stmts)
	}
	if len(errors) != 1 {
		t.Log(errors)
		t.Fatalf("expect len(1) got %v", len(errors))
	}
	if errors[0] != "Invalid assignment target." {
		t.Fatalf("expect 'Invalid assignment target.' got '%v'", errors[0])
	}
}

func TestParseLiteral(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if len(errors) != 0 {
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if stmts == nil {
		t.Fatalf("expect not nil")
	}
}

func TestParseUnary(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if len(errors) != 0 {
		t.Fatalf("expect empty len(error) got %d, %v", len(errors), errors[0])
	}
	if stmts == nil {
		t.Fatal("expect not nil")
	}
}

func TestParseTerm(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if stmts == nil {
		t.Fatal("expect not nil")
	}
}

func TestParseFactorial(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if len(errors) != 0 {
		t.Log(errors)
		t.Fatalf("expect empty len(error) got %d", len(errors))
	}
	if stmts == nil {
		t.Fatal("expect not nil")
	}
}

func TestParseGrouping(t *testing.T) {
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
}

func TestParseGroupingError(t *testing.T) {
	errors := make([]string, 0)
	p := NewParser(
		[]scanner.Token{
			testutil.LeftParen(),
			testutil.Number(321.0),
			testutil.Star(),
			testutil.Number(123.0),
			testutil.Semicolon(),
			testutil.Eof(),
		},
		testCallBack(&errors),
	)
	stmts, err := p.Parse()
	if err != nil {
		t.Fatal("expect nil")
	}
	if len(errors) != 1 {
		t.Fatalf("expect len(1) got %d", len(errors))
	}
	if len(stmts) != 0 {
		t.Log(stmts)
		t.Fatalf("expect empty got %v", len(stmts))
	}
	if errors[0] != "expect ')' after expression." {
		t.Fatalf("expect %v got %v", "expect ')' after expression.", errors[0])
	}
}

func TestParseGroupingErrors(t *testing.T) {
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
	if err != nil {
		t.Fatalf("expect empty, got %v", err.Error())
	}
	if len(stmts) > 0 {
		t.Log(stmts)
		t.Fatalf("Expecting empty got: %v", len(stmts))
	}
	if len(errors) != 1 {
		t.Fatalf("expect len(1) got %d", len(errors))
	}
	if errors[0] != "expect ')' after expression." {
		t.Fatalf("expect: 'expect ')' after expression.', got: '%v'", errors[0])
	}
}

func TestParseBlock(t *testing.T) {
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
	got, ok := stmts[0].(*token.Block)
	if !ok {
		t.Fatalf("expect *token.Block got %T", got)
	}
	if len(got.Statements) != 1 {
		t.Fatalf("expect len(1) got: %d", len(got.Statements))
	}
	gotStmt := got.Statements[0]
	expr, ok := gotStmt.(*token.Expression)
	if !ok {
		t.Fatalf("expect *token.Expression got %T", gotStmt)
	}
	_, ok = expr.Expression.(*token.Binary)
	if !ok {
		t.Fatalf("expect *token.Expression got %T", expr.Expression)
	}
}

func TestParseBlockError(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}
	if len(stmts) > 0 {
		t.Fatalf("expect empty, got %v", stmts)
	}
	if len(errors) != 1 {
		t.Log(errors)
		t.Fatalf("expect len(1) got '%v'", len(errors))
	}
	if errors[0] != "Expect '}' after block." {
		t.Fatalf("expect 'Expect '}' after block.' got %v", errors[0])
	}
}

var testCallBack = func(errs *[]string) ErrorCallback {
	return func(token scanner.Token, message string) {
		*errs = append(*errs, message)
	}
}
