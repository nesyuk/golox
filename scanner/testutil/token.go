package testutil

import (
	"fmt"
	"github.com/nesyuk/golox/scanner"
)

func IfIdentifier() scanner.Token {
	lexeme := "if"
	return scanner.Token{TokenType: scanner.IF, Lexeme: &lexeme, Line: 1}
}

func ElseIdentifier() scanner.Token {
	lexeme := "else"
	return scanner.Token{TokenType: scanner.ELSE, Lexeme: &lexeme, Line: 1}
}

func PrintIdentifier() scanner.Token {
	lexeme := "print"
	return scanner.Token{TokenType: scanner.PRINT, Lexeme: &lexeme, Line: 1}
}

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

func Bang() scanner.Token {
	lexeme := "!"
	return scanner.Token{TokenType: scanner.BANG, Lexeme: &lexeme, Line: 1}
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

func Slash() scanner.Token {
	lexeme := "/"
	return scanner.Token{TokenType: scanner.SLASH, Lexeme: &lexeme, Line: 1}
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

func Greater() scanner.Token {
	lexeme := ">"
	return scanner.Token{TokenType: scanner.GREATER, Lexeme: &lexeme, Line: 1}
}

func GreaterEqual() scanner.Token {
	lexeme := ">="
	return scanner.Token{TokenType: scanner.GREATER_EQUAL, Lexeme: &lexeme, Line: 1}
}

func Less() scanner.Token {
	lexeme := "<"
	return scanner.Token{TokenType: scanner.LESS, Lexeme: &lexeme, Line: 1}
}

func LessEqual() scanner.Token {
	lexeme := "<="
	return scanner.Token{TokenType: scanner.LESS_EQUAL, Lexeme: &lexeme, Line: 1}
}

func EqualEqual() scanner.Token {
	lexeme := "=="
	return scanner.Token{TokenType: scanner.EQUAL_EQUAL, Lexeme: &lexeme, Line: 1}
}

func BangEqual() scanner.Token {
	lexeme := "!="
	return scanner.Token{TokenType: scanner.BANG_EQUAL, Lexeme: &lexeme, Line: 1}
}

func Eof() scanner.Token {
	return scanner.Token{TokenType: scanner.EOF, Line: 1}
}
