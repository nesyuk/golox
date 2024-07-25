package parser

import (
	"errors"
	"fmt"
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
)

/*
    program -> declaration* EOF
    declaration -> classDecl | funDecl | varDecl | statement
    classDecl -> "class" IDENTIFIER ("<" IDENTIFIER)? "{" function* "}"
    funDecl -> "fun" function
    function -> IDENTIFIER "(" parameters? ")" block
    parameters -> IDENTIFIER ("," IDENTIFIER)*
    varDecl -> "var" IDENTIFIER ("=" expression)? ";"
    statement -> exprStmt | printStmt | block
    exprStmt -> expression ";"
    forStmt -> "for" "( (varDecl | exprStmt) ";" expr? ";" expr? ")" statement
	ifStmt ->  "if" "(" expression ")" statement ( "else" statement )?
    printStmt -> "print" expression
    returnStmt -> "return" expression? ";"
	whileStmt -> "while" "(" expression ")" statement
    block -> "{" declaration* "}"
	expression -> assignment
    assignment -> (call ".")? IDENTIFIER "=" assignment | logicOr
	logicOr -> logicAnd ( "or" logicAnd )*
	logicAnd-> equality ( "and" equality )*
	equality -> comparison (("!=" | "==") comparison)*
	comparison -> term  ((">" | ">=" | "<" | "<=") term)*
	term -> factor (("-" | "+") factor)*
    factor -> unary (( "/" | "*") unary)*
    unary -> ("!" | "-") unary | call
    call -> primary ( "(" arguments? ")" | "." IDENTIFIER )*
    arguments -> expression ( "," expression )*
    primary -> NUMBER | STRING | "true" | "false" | nil | "(" + expression + ")" | IDENTIFIER | "this"
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
	if p.match(scanner.CLASS) {
		return p.class()
	} else if p.match(scanner.FUN) {
		return p.function("function")
	} else if p.match(scanner.VAR) {
		return p.variableDeclaration() //TODO: verify p.statementSync(p.variableDeclaration()) or p.variableDeclaration()
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

func (p *Parser) class() (token.Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expect class name")
	if err != nil {
		return nil, err
	}
	var supercls *token.VariableExpr
	if p.match(scanner.LESS) {
		if _, err = p.consume(scanner.IDENTIFIER, "Expect superclass name"); err != nil {
			return nil, nil
		}
		supercls = &token.VariableExpr{Name: p.previous()}
	}
	if _, err := p.consume(scanner.LEFT_BRACE, "Expect '{' before class body."); err != nil {
		return nil, err
	}
	methods := make([]*token.FunctionStmt, 0)
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.function("method")
		method := stmt.(*token.FunctionStmt)
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}
	if _, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after class body."); err != nil {
		return nil, err
	}
	return &token.ClassStmt{
		Name:       name,
		Superclass: supercls,
		Methods:    methods,
	}, nil
}

func (p *Parser) function(kind string) (token.Stmt, error) {
	tok, err := p.consume(scanner.IDENTIFIER, fmt.Sprintf("expect %v name.", kind))
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.LEFT_PAREN, fmt.Sprintf(fmt.Sprintf("expect '(' after %v name.", kind)))
	if err != nil {
		return nil, err
	}
	params := make([]*scanner.Token, 0)
	if !p.check(scanner.RIGHT_PAREN) {
		param, err := p.consume(scanner.IDENTIFIER, fmt.Sprintf("expect %v name.", kind))
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		for p.match(scanner.COMMA) {
			if len(params) >= 255 {
				return nil, p.error(p.peek(), "can't have more than 255 parameters.")
			}
			param, err = p.consume(scanner.IDENTIFIER, fmt.Sprintf("expect parameter name."))
			if err != nil {
				return nil, err
			}
			params = append(params, param)
		}
	}

	if _, err = p.consume(scanner.RIGHT_PAREN, fmt.Sprintf(fmt.Sprintf("expect ')' after parameters."))); err != nil {
		return nil, err
	}
	if _, err = p.consume(scanner.LEFT_BRACE, fmt.Sprintf(fmt.Sprintf("expect '{' before %v body.", kind))); err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &token.FunctionStmt{
		Name:   tok,
		Params: params,
		Body:   body,
	}, nil
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
	case p.match(scanner.FOR):
		return p.forStatement()
	case p.match(scanner.IF):
		return p.ifStatement()
	case p.match(scanner.PRINT):
		return p.printStmt()
	case p.match(scanner.RETURN):
		return p.returnStmt()
	case p.match(scanner.WHILE):
		return p.whileStatement()
	case p.match(scanner.LEFT_BRACE):
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return &token.BlockStmt{Statements: stmts}, nil
	}
	return p.expressionStmt()
}

func (p *Parser) forStatement() (token.Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}

	var initializer token.Stmt
	var err error
	switch {
	case p.match(scanner.SEMICOLON):
	// no initialization
	case p.match(scanner.VAR):
		initializer, err = p.variableDeclaration()
		if err != nil {
			return nil, err
		}
	default:
		initializer, err = p.expressionStmt()
		if err != nil {
			return nil, err
		}
	}

	var condition token.Expr
	if !p.check(scanner.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(scanner.SEMICOLON, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	var increment token.Expr
	if !p.check(scanner.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after for clauses."); err != nil {
		return nil, err
	}

	if condition == nil {
		condition = &token.LiteralExpr{Value: true}
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &token.BlockStmt{Statements: []token.Stmt{body, &token.ExpressionStmt{Expression: increment}}}
	}

	body = &token.WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = &token.BlockStmt{Statements: []token.Stmt{initializer, body}}
	}

	return body, nil
}

func (p *Parser) whileStatement() (token.Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after condition. "); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &token.WhileStmt{Condition: condition, Body: body}, nil
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

func (p *Parser) returnStmt() (token.Stmt, error) {
	keyword := p.previous()
	var value token.Expr
	var err error
	if !p.check(scanner.SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(scanner.SEMICOLON, "Expect ';' after return value."); err != nil {
		return nil, err
	}
	return &token.ReturnStmt{
		Keyword: &keyword,
		Value:   value,
	}, nil
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

		switch t := expr.(type) {
		case *token.VariableExpr:
			return &token.AssignExpr{Name: t.Name, Value: value}, nil
		case *token.GetExpr:
			return &token.SetExpr{Object: t.Object, Name: t.Name, Value: value}, nil
		}
		return nil, p.error(tokenEquals, "Invalid assignment target.")
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
	return p.call()
}

func (p *Parser) call() (token.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(scanner.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(scanner.DOT) {
			name, err := p.consume(scanner.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &token.GetExpr{Object: expr, Name: name}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee token.Expr) (token.Expr, error) {
	args := make([]token.Expr, 0)
	if !p.check(scanner.RIGHT_PAREN) {
		arg, err := p.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		for p.match(scanner.COMMA) {
			arg, err = p.expression()
			if err != nil {
				return nil, err
			}
			if len(args) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			args = append(args, arg)
		}
	}
	tok, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return &token.CallExpr{
		Callee:    callee,
		Paren:     tok,
		Arguments: args,
	}, nil
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
	case p.match(scanner.THIS):
		return &token.ThisExpr{Keyword: p.previous()}, nil
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
