package lox

import (
	"bufio"
	"fmt"
	"github.com/nesyuk/golox/interpreter"
	"github.com/nesyuk/golox/parser"
	"github.com/nesyuk/golox/scanner"
	"io"
	"os"
)

type golox struct {
	interpret       *interpreter.Interpreter
	hadError        bool
	hadRuntimeError bool
}

func newLox() *golox {
	return &golox{&interpreter.Interpreter{}, false, false}
}

func (l *golox) run(source string) (string, error) {
	if l.hadError {
		os.Exit(65)
	}
	if l.hadRuntimeError {
		os.Exit(70)
	}
	sc := scanner.NewScanner(source, l.error)
	tokens := sc.ScanTokens()

	p := parser.NewParser(tokens, l.parseError)
	ast, err := p.Parse()
	if err != nil || l.hadError {
		return "", err
	}

	i := interpreter.New(l.runtimeError)
	result, err := i.Interpret(ast)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (l *golox) runtimeError(err *interpreter.RuntimeError) {
	fmt.Printf("%v\n[line %d]\n", err.Error(), err.Token.Line)
	l.hadRuntimeError = true
}

func (l *golox) parseError(token scanner.Token, message string) {
	if token.TokenType == scanner.EOF {
		l.report(token.Line, " at end", message)
	} else {
		l.report(token.Line, fmt.Sprintf(" at '%v'", *token.Lexeme), message)
	}
}

func (l *golox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *golox) report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%v: %v\n", line, where, message)
	l.hadError = true
}

func (l *golox) ResetError() {
	l.hadError = false
}

func RunFile(f string) error {
	s, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	lox := newLox()
	result, err := lox.run(string(s))
	if err != nil {
		os.Exit(65)
	}
	fmt.Println(result)
	return nil
}

func RunPrompt() {
	lox := newLox()
	for {
		fmt.Print("> ")
		r := bufio.NewReader(os.Stdin)
		s, err := r.ReadString('\n')
		if err == io.EOF {
			return
		}
		result, err := lox.run(s)
		if err != nil {
			fmt.Printf("failed to interpret: %v\n", err)
			lox.ResetError()
		} else {
			fmt.Println(result)
		}
	}
}
