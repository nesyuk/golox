package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nesyuk/golox/parser"
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
	"io"
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
	printer := token.AstPrinter{}
	printer.Print(ast)
	return nil
}

func runFile(f string) error {
	s, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	if err := run(string(s)); err != nil {
		os.Exit(65)
	}
	return nil
}

func runPrompt() {
	for {
		fmt.Print("> ")
		r := bufio.NewReader(os.Stdin)
		s, err := r.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err := run(s); err != nil {
			fmt.Printf("failed to interpret: %v\n", err)
			hadError = false
		}
	}
}

func report(line int, loc string, msg string) {
	fmt.Printf("[line %d] Error %v: %v\n", line, loc, msg)
	hadError = true
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println(errors.New("usage: glox [script]"))
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Printf("failed to read a file: %v\n", err)
		}
	} else {
		runPrompt()
	}
}
