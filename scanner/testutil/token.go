package testutil

import (
	"fmt"
	"github.com/nesyuk/golox/scanner"
)

func VarDecl() scanner.Token {
	lexeme := "var"
	return scanner.Token{TokenType: scanner.VAR, Lexeme: &lexeme, Line: 1}
}

func Equal() scanner.Token {
	lexeme := "="
	return scanner.Token{TokenType: scanner.EQUAL, Lexeme: &lexeme, Line: 1}
}

func Identifier(name string) scanner.Token {
	return scanner.Token{TokenType: scanner.IDENTIFIER, Lexeme: &name, Literal: name, Line: 1}
}

func LeftParen() scanner.Token {
	lexeme := "("
	return scanner.Token{TokenType: scanner.LEFT_PAREN, Lexeme: &lexeme, Line: 1}
}

func RightParen() scanner.Token {
	lexeme := ")"
	return scanner.Token{TokenType: scanner.RIGHT_PAREN, Lexeme: &lexeme, Line: 1}
}

func LeftBrace() scanner.Token {
	lexeme := "{"
	return scanner.Token{TokenType: scanner.LEFT_BRACE, Lexeme: &lexeme, Line: 1}
}

func RightBrace() scanner.Token {
	lexeme := "}"
	return scanner.Token{TokenType: scanner.RIGHT_BRACE, Lexeme: &lexeme, Line: 1}
}

func Minus() scanner.Token {
	lexeme := "-"
	return scanner.Token{TokenType: scanner.MINUS, Lexeme: &lexeme, Line: 1}
}

func Plus() scanner.Token {
	lexeme := "+"
	return scanner.Token{TokenType: scanner.PLUS, Lexeme: &lexeme, Line: 1}
}

func Star() scanner.Token {
	lexeme := "*"
	return scanner.Token{TokenType: scanner.STAR, Lexeme: &lexeme, Line: 1}
}

func Number(n float64) scanner.Token {
	lexeme := fmt.Sprintf("%v", n)
	return scanner.Token{TokenType: scanner.NUMBER, Lexeme: &lexeme, Literal: n, Line: 1}
}

func Str(s string) scanner.Token {
	return scanner.Token{TokenType: scanner.STRING, Lexeme: &s, Literal: s, Line: 1}
}

func Semicolon() scanner.Token {
	lexeme := ";"
	return scanner.Token{TokenType: scanner.SEMICOLON, Lexeme: &lexeme, Line: 1}
}

func Eof() scanner.Token {
	return scanner.Token{TokenType: scanner.EOF, Line: 1}
}
