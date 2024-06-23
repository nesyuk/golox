package parser

import (
	"github.com/nesyuk/golox/scanner"
	"github.com/nesyuk/golox/token"
)

/*
	expression -> equality
	equality -> comparison ((!= | ==) comparison)*
	comparison -> term  ((> | >= | < | <=) term)*
	term -> factor ((- | +) factor)*
    factor -> unary (( / | *) unary)*
    unary -> (! | -) unary | primary
    primary -> NUMBER | STRING | "true" | "false" | nil | "(" + expression + ")"
*/

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() (token.Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (token.Expr, error) {
	return p.equality()
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
		expr = &token.Binary{Left: expr, Operation: op, Right: right}
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
		expr = &token.Binary{Left: expr, Operation: op, Right: right}
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
		expr = &token.Binary{Left: expr, Operation: op, Right: right}
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
		expr = &token.Binary{Left: expr, Operation: op, Right: right}
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
		return &token.Unary{Operation: op, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (token.Expr, error) {
	switch {
	case p.match(scanner.FALSE):
		return &token.Literal{Value: false}, nil
	case p.match(scanner.TRUE):
		return &token.Literal{Value: true}, nil
	case p.match(scanner.NIL):
		return &token.Literal{Value: nil}, nil
	case p.match(scanner.NUMBER) || p.match(scanner.STRING):
		return &token.Literal{Value: p.previous().Literal}, nil
	case p.match(scanner.LEFT_BRACE):
		expr, err := p.expression()
		if err != nil {
		}
		_, err = p.consume(scanner.RIGHT_BRACE, "Expect ')' after expression.")
		if err != nil {
			// TODO: handle error
			return nil, err
		}
		return &token.Grouping{Expression: expr}, nil

	}
	return nil, &Error{"Expect expression"}
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

func (p *Parser) consume(tt scanner.TokenType, msg string) (*scanner.Token, error) {
	if p.check(tt) {
		t := p.advance()
		return &t, nil
	}
	// TODO: Lox.error()
	return nil, &Error{Message: msg}
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

type Error struct {
	Message string
}

func (err *Error) Error() string {
	return ""
}
