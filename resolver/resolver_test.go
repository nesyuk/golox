package resolver

import (
	"github.com/nesyuk/golox/interpreter"
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/scanner/testutil"
	"github.com/nesyuk/golox/token"
	"testing"
)

func TestResolver_Resolve(t *testing.T) {
	tests := []struct {
		stmts       []token.Stmt
		staticErrs  []string
		runtimeErrs []string
	}{
		{[]token.Stmt{
			&token.ReturnStmt{
				Keyword: testutil.Return(),
				Value:   &token.LiteralExpr{Value: "nasty stuff"},
			},
		}, []string{"Can't return from top-level code."},
			[]string{}},
		{
			[]token.Stmt{
				&token.VarStmt{
					Name:        scanner.Token{},
					Initializer: nil,
				},
			},
			[]string{},
			[]string{},
		},
	}
	for _, test := range tests {
		runtimeErrs := make([]string, 0)
		interpr := interpreter.New(runtimeTestCallBack(&runtimeErrs))
		staticErrs := make([]string, 0)
		res := New(interpr, testCallBack(&staticErrs))
		res.Resolve(test.stmts)
		checkErrors(t, test.runtimeErrs, runtimeErrs)
		checkErrors(t, test.staticErrs, staticErrs)
	}
}

func checkErrors(t *testing.T, expected []string, got []string) {
	if len(got) != len(expected) {
		t.Log(got)
		t.Fatalf("expect %v, got: %v", len(expected), len(got))
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Logf("expect: %v, got: %v", expected[i], got[i])
		}
	}
}

var testCallBack = func(errs *[]string) ErrorCallback {
	return func(token scanner.Token, message string) {
		*errs = append(*errs, message)
	}
}

var runtimeTestCallBack = func(errs *[]string) interpreter.ErrorCallback {
	return func(err *interpreter.RuntimeError) {
		*errs = append(*errs, err.Error())
	}
}
