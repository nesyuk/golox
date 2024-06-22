package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GeneratorError struct {
	Msg string
	Err error
}

func (e *GeneratorError) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.Err)
}

func (e *GeneratorError) wrap(err error) *GeneratorError {
	e.Err = err
	return e
}

var productions = []string{
	"Literal: Value interface{}",
	"Unary: Operation scanner.Token, Right Expr",
	"Binary: Left Expr, Operation scanner.Token, Right Expr",
	"Grouping: Expression Expr",
}

func generate(filename string) error {
	cwd, err := os.Getwd()
	fmt.Println(cwd)
	if err != nil {
		return fmt.Errorf("failed to determine current dir: %w", err)
	}
	if _, err := os.Stat(filepath.Join(cwd, "token")); os.IsNotExist(err) {
		return fmt.Errorf("target folder doesn't exist: %v", filepath.Join(cwd, "token"))
	}
	f, err := os.Create(filepath.Join(cwd, "token", fmt.Sprintf("%v.go", filename)))
	defer f.Close()
	if err != nil {
		return fmt.Errorf("failed to create file: %w\n", err)
	}

	f.WriteString("package token\n\n")
	declareImports(f, "github.com/nesyuk/golox/scanner")

	f.WriteString("type Expr interface {\n\tAccept(visitor Visitor) interface{}\n}\n\n")

	declareVisitorInterface(f)

	declareExpressions(f)
	return nil
}

func declareExpressions(f *os.File) {
	for _, prod := range productions {
		prodArr := strings.Split(prod, ":")
		head := prodArr[0]
		body := prodArr[1]
		f.WriteString(fmt.Sprintf("type %v struct {", head))
		for _, part := range strings.Split(body, ",") {
			part = strings.TrimLeft(part, " ")
			f.WriteString(fmt.Sprintf("\n\t%v", part))
		}
		f.WriteString("\n}\n\n")
		f.WriteString(fmt.Sprintf("func (e *%v) Accept(visitor Visitor) interface{} {\n\treturn visitor.Visit%v(e)\n}\n\n", head, head))
	}
	/*
		type Literal struct {
			Value interface{}
		}

		func (e *Literal) Accept(visitor Visitor) interface{} {
		return visitor.VisitLiteral(e)
		}*/
}

func declareVisitorInterface(f *os.File) {
	f.WriteString("type Visitor interface {\n")
	for _, prod := range productions {
		prodArr := strings.Split(prod, ":")
		head := prodArr[0]
		f.WriteString(fmt.Sprintf("\tVisit%v(expr *%v) interface{}\n", head, head))
	}
	f.WriteString("}\n\n")
}

func declareImports(f *os.File, imports ...string) {
	f.WriteString("import (\n")
	for _, imp := range imports {
		f.WriteString("\t\"" + imp + "\"" + "\n")
	}
	f.WriteString(")\n\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: tokengen [filename]")
		os.Exit(64)
	}
	if err := generate(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
