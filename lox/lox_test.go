package lox

import "testing"

func TestRun(t *testing.T) {
	tests := []struct {
		expr            string
		expect          string
		hadError        bool
		hadRuntimeError bool
	}{
		{"print 3 + 2;", "5", false, false},
		{"print 3.3 + 2.2;", "5.5", false, false},
		{"print 4 * 5;", "20", false, false},
		{"print \"hello,\" + \" world!\";", "hello, world!", false, false},
		{"print (3 + 2;", "", true, false},
		{"print 3 + \"2\";", "", false, true},
		{"var a = 1; var b = 2; print a + b;", "", false, false},
		{"var a = 1; {var a = 2; print a;} print a;", "", false, false},
		{"var i = 1; while(i < 3) {print i; i = i + 1;}", "", false, false},
		{"var a = 0; var temp; for (var b = 1; a < 10000; b = temp + b) { print a; temp = a; a = b;}", "", false, false},
	}
	for _, test := range tests {
		lox := newLox()
		err := lox.run(test.expr)
		if err != nil {
			t.Error(err)
		}
		if test.hadError != lox.hadError {
			t.Fatalf("expect: '%v' got '%v'", test.hadError, lox.hadError)
		}
		if test.hadRuntimeError != lox.hadRuntimeError {
			t.Fatalf("expect: %v got %v", test.hadRuntimeError, lox.hadRuntimeError)
		}
		/*		if got != test.expect {
				t.Fatalf("expect: %v, got: %v", test.expect, got)
			}*/
	}
}
