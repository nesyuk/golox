package parser

import (
	"errors"
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
)

/*
    program -> declaration* EOF
    declaration -> varDecl | statement
    varDecl -> "var" IDENTIFIER ("=" expression)? ";"
    statement -> exprStmt | printStmt | block
    exprStmt -> expression ";"
	ifStmt ->  "if" "(" expression ")" statement ( "else" statement )?
    printStmt -> "print" expression
    block -> "{" declaration* "}"
	expression -> assignment
    assignment -> IDENTIFIER "=" assignment | logicOr
	logicOr -> logicAnd ( "or" logicAnd )*
	logicAnd-> equality ( "and" equality )*
	equality -> comparison (("!=" | "==") comparison)*
	comparison -> term  ((">" | ">=" | "<" | "<=") term)*
	term -> factor (("-" | "+") factor)*
    factor -> unary (( "/" | "*") unary)*
    unary -> ("!" | "-") unary | primary
    primary -> NUMBER | STRING | "true" | "false" | nil | "(" + expression + ")" | IDENTIFIER
*/

type Parser struct {
	tokens        []scanner.Token
	current       int
	errorCallback ErrorCallback
}

func NewParser(tokens []scanner.Token, onError ErrorCallback) *Parser {
	return &Parser{tokens: tokens, current: 0, errorCallback: onError}
}

func (p *Parser) Parse() ([]token.Stmt, error) {
	stmts := make([]token.Stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil || stmt == nil {
			return make([]token.Stmt, 0), err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) declaration() (token.Stmt, error) {
	if p.match(scanner.VAR) {
		return p.statementSync(p.variableDeclaration())
	}
	return p.statementSync(p.statement())
}

func (p *Parser) statementSync(stmt token.Stmt, err error) (token.Stmt, error) {
	var parseErr *ParseError
	if err != nil && errors.As(err, &parseErr) {
		p.sync()
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) variableDeclaration() (token.Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "expect variable name")
	if err != nil || name == nil {
		return nil, err
	}
	var initializer token.Expr
	if p.match(scanner.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.SEMICOLON, "expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return &token.VarStmt{Name: *name, Initializer: initializer}, nil
}

func (p *Parser) statement() (token.Stmt, error) {
	switch {
	case p.match(scanner.IF):
		return p.ifStatement()
	case p.match(scanner.PRINT):
		return p.printStmt()
	case p.match(scanner.LEFT_BRACE):
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return &token.BlockStmt{Statements: stmts}, nil
	}
	return p.expressionStmt()
}

func (p *Parser) ifStatement() (token.Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "expect ')' after if condition")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch token.Stmt
	if p.match(scanner.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &token.IfStmt{Condition: cond, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) printStmt() (token.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(scanner.SEMICOLON, "expect ';' after value."); err != nil {
		return nil, err
	}
	return &token.PrintStmt{Expression: expr}, err
}

func (p *Parser) block() ([]token.Stmt, error) {
	stmts := make([]token.Stmt, 0)
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if _, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}
	return stmts, nil
}

func (p *Parser) expressionStmt() (token.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(scanner.SEMICOLON, "expect ';' after expression."); err != nil {
		return nil, err
	}
	return &token.ExpressionStmt{Expression: expr}, err
}

func (p *Parser) expression() (token.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (token.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(scanner.EQUAL) {
		// it's an assigment, so left part (l-value) must be a variable.
		tokenEquals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		name, ok := expr.(*token.VariableExpr)
		if !ok {
			return nil, p.error(tokenEquals, "Invalid assignment target.")
		}
		return &token.AssignExpr{Name: name.Name, Value: value}, nil
	}
	// it's not an 'or' expression, return it
	return expr, nil
}

func (p *Parser) or() (token.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(scanner.OR) {
		op := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &token.LogicalExpr{Operator: op, Left: expr, Right: right}
	}
	return expr, nil
}

func (p *Parser) and() (token.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(scanner.AND) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &token.LogicalExpr{Operator: op, Left: expr, Right: right}
	}
	return expr, nil
}

func (p *Parser) equality() (token.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &token.BinaryExpr{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (token.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &token.BinaryExpr{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) term() (token.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(scanner.MINUS, scanner.PLUS) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &token.BinaryExpr{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) factor() (token.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.SLASH, scanner.STAR) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &token.BinaryExpr{Left: expr, Operator: op, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (token.Expr, error) {
	if p.match(scanner.BANG, scanner.MINUS) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &token.UnaryExpr{Operator: op, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (token.Expr, error) {
	switch {
	case p.match(scanner.FALSE):
		return &token.LiteralExpr{Value: false}, nil
	case p.match(scanner.TRUE):
		return &token.LiteralExpr{Value: true}, nil
	case p.match(scanner.NIL):
		return &token.LiteralExpr{Value: nil}, nil
	case p.match(scanner.NUMBER) || p.match(scanner.STRING):
		return &token.LiteralExpr{Value: p.previous().Literal}, nil
	case p.match(scanner.IDENTIFIER):
		return &token.VariableExpr{Name: p.previous()}, nil
	case p.match(scanner.LEFT_PAREN):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err = p.consume(scanner.RIGHT_PAREN, "expect ')' after expression."); err != nil {
			return nil, err
		}
		return &token.GroupingExpr{Expression: expr}, err

	}
	return nil, p.error(p.peek(), "expect expression")
}

func (p *Parser) match(tokens ...scanner.TokenType) bool {
	for _, t := range tokens {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tt scanner.TokenType, message string) (*scanner.Token, error) {
	if p.check(tt) {
		t := p.advance()
		return &t, nil
	}
	return nil, p.error(p.peek(), message)
}

func (p *Parser) error(t scanner.Token, message string) error {
	p.errorCallback(t, message)
	return &ParseError{Message: message}
}

func (p *Parser) sync() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == scanner.SEMICOLON {
			return
		}
		switch p.peek().TokenType {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR, scanner.IF, scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) check(t scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}

type ErrorCallback = func(scanner.Token, string)
