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

var expressions = []string{
	"AssignExpr: Name scanner.Token, Value Expr",
	"LiteralExpr: Value interface{}",
	"LogicalExpr: Left Expr, Operator scanner.Token, Right Expr",
	"UnaryExpr: Operator scanner.Token, Right Expr",
	"VariableExpr: Name scanner.Token",
	"BinaryExpr: Left Expr, Operator scanner.Token, Right Expr",
	"GroupingExpr: Expression Expr",
}

var statements = []string{
	"BlockStmt: Statements []Stmt",
	"ExpressionStmt: Expression Expr",
	"IfStmt: Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
	"PrintStmt: Expression Expr",
	"VarStmt: Name scanner.Token, Initializer Expr",
}

func generateAst(filename string) error {
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

	defineExpressions(f, expressions, "Expr")
	defineExpressions(f, statements, "Stmt")

	return nil
}

func defineExpressions(f *os.File, productions []string, name string) {
	f.WriteString(fmt.Sprintf("type %v interface {\n\tAccept(visitor Visitor%v) (interface{}, error)\n}\n\n", name, name))

	declareVisitorInterface(f, productions, name)

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
		f.WriteString(fmt.Sprintf("func (e *%v) Accept(visitor Visitor%v) (interface{}, error) {\n\treturn visitor.Visit%v(e)\n}\n\n", head, name, head))
	}
}

func declareVisitorInterface(f *os.File, productions []string, name string) {
	f.WriteString(fmt.Sprintf("type Visitor%v interface {\n", name))
	for _, prod := range productions {
		prodArr := strings.Split(prod, ":")
		head := prodArr[0]
		f.WriteString(fmt.Sprintf("\tVisit%v(%v *%v) (interface{}, error)\n", head, strings.ToLower(name), head))
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
	if err := generateAst(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
