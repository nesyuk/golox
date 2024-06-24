package lox

import (
	"bufio"
	"fmt"
	"github.com/nesyuk/golox/parser"
	"github.com/nesyuk/golox/printer"
	"github.com/nesyuk/golox/scanner"
	"io"
	"log"
	"os"
)

var hadError = false

func run(source string) error {
	sc := scanner.NewScanner(source)
	tokens, err := sc.ScanTokens()
	if err != nil {
		// TODO: handle scanner error
		return err
	}
	p := parser.NewParser(tokens)
	ast, err := p.Parse()
	if err != nil && hadError {
		// TODO handle parser error
		return err
	}
	if err != nil {

	}
	printer := printer.Ast{}
	printer.Print(ast)
	return nil
}

func logError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Printf("[line %d] Error%v: %v\n", line, where, message)
	hadError = true
}

func RunFile(f string) error {
	s, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	if err := run(string(s)); err != nil {
		os.Exit(65)
	}
	return nil
}

func RunPrompt() {
	for {
		fmt.Print("> ")
		r := bufio.NewReader(os.Stdin)
		s, err := r.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err = run(s); err != nil {
			fmt.Printf("failed to interpret: %v\n", err)
			hadError = false
		}
	}
}
