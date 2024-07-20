package runtime

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		expr         string
		expect       []string
		errors       []string
		runtimeError bool
	}{
		{"print 3 + 2;", []string{"5"}, []string{}, false},
		{"print 3.3 + 2.2;", []string{"5.5"}, []string{}, false},
		{"print 4 * 5;", []string{"20"}, []string{}, false},
		{"print \"hello,\" + \" world!\";", []string{"hello, world!"}, []string{}, false},
		{"print (3 + 2;", []string{}, []string{"[line 1] Error at ';': expect ')' after expression.\n"}, false},
		{"print 3 + \"2\";", []string{}, []string{"Operands must be numbers: 2\n[line 1]\n"}, true},
		{"var a = 1; var b = 2; print a + b;", []string{"3"}, []string{}, false},
		{"var a = 1; {var a = 2; print a;} print a;", []string{"2", "1"}, []string{}, false},
		{"var i = 1; while(i < 3) {print i; i = i + 1;}", []string{"1", "2"}, []string{}, false},
		{"var a = 0; var temp; for (var b = 1; a < 10000; b = temp + b) { print a; temp = a; a = b;}", []string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89", "144", "233", "377", "610", "987", "1597", "2584", "4181", "6765"}, []string{}, false},
		{"fun sayHi(first, last) { print \"Hi, \" + first + \" \" + last + \"!\"; }\n sayHi(\"Mr.\", \"Bean\");", []string{"Hi, Mr. Bean!"}, []string{}, false},
		{"fun fib(n) {\nif (n <= 1) return n;\nreturn fib(n-2) + fib(n-1);\n}\n\nfor (var i = 0; i < 20; i = i + 1) {\nprint fib(i);\n}", []string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89", "144", "233", "377", "610", "987", "1597", "2584", "4181"}, []string{}, false},
		{"fun makeCounter() {\nvar i = 0;\nfun count() {\ni = i + 1;\nprint i;\n }\nreturn count;\n}\n\nvar counter = makeCounter();\ncounter(); // 1\ncounter(); // 2", []string{"1", "2"}, []string{}, false},
		{"class Greeting {\n\thello() {\n\t\treturn \"Hello\";\n\t}\n}\n\nprint Greeting;", []string{"<class 'Greeting'.>"}, []string{}, false},
	}

	for _, test := range tests {
		reporter := newTestReporter()
		lox := NewLox(reporter)
		err := lox.run(test.expr)
		if err != nil {
			t.Error(err)
		}
		if len(test.errors) != 0 != lox.hadError && !lox.hadRuntimeError {
			t.Errorf("expect no error, got '%v'", lox.hadError)
		}
		if test.runtimeError != lox.hadRuntimeError {
			t.Errorf("expect runtime error: %v got %v (in %v)", test.runtimeError, lox.hadRuntimeError, test.expr)
		}
		reporter.Validate(t, test.expect, test.errors, test.expr)
	}
}

type testReporter struct {
	errors []string
	got    []string
}

func newTestReporter() *testReporter {
	return &testReporter{make([]string, 0), make([]string, 0)}
}

func (r *testReporter) Error(format string, args ...interface{}) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

func (r *testReporter) Print(s string) {
	r.got = append(r.got, s)
}

func (r *testReporter) Validate(t *testing.T, expect []string, expectErrors []string, test string) {
	r.validate(t, "result", r.got, expect, test)
	r.validate(t, "error", r.errors, expectErrors, test)
}

func (r *testReporter) validate(t *testing.T, prefix string, got []string, expect []string, test string) {
	if len(got) != len(expect) {
		t.Errorf("%s: expect %v, got: %v (test: '%v')", prefix, len(expect), len(got), test)
		t.Log(got)
		return
	}
	for i := range got {
		if got[i] != expect[i] {
			t.Errorf("%s: expect: %v, got: %v", prefix, expect[i], got[i])
		}
	}
}
