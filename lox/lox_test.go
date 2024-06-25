package lox

import "testing"

func TestRun(t *testing.T) {
	tests := []struct {
		expr            string
		expect          string
		hadError        bool
		hadRuntimeError bool
	}{
		{"3 + 2", "5", false, false},
		{"3.3 + 2.2", "5.5", false, false},
		{"4 * 5", "20", false, false},
		{"\"hello,\" + \" world!\"", "hello, world!", false, false},
		{"(3 + 2", "", true, false},
		{"3 + \"2\"", "", false, true},
	}
	for _, test := range tests {
		lox := newLox()
		got, err := lox.run(test.expr)
		if err != nil {
			t.Error(err)
		}
		if test.hadError != lox.hadError {
			t.Fatalf("expect: %v got %v", test.hadError, lox.hadError)
		}
		if test.hadRuntimeError != lox.hadRuntimeError {
			t.Fatalf("expect: %v got %v", test.hadRuntimeError, lox.hadRuntimeError)
		}
		if got != test.expect {
			t.Fatalf("expect: %v, got: %v", test.expect, got)
		}
	}
}
