package scanner

import (
	"errors"
	"github.com/nesyuk/golox/formatter"
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source         string
	tokens         []Token
	start, current int
	line           int
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, tokens: make([]Token, 0), start: 0, current: 0, line: 1}
}

func (sc *Scanner) ScanTokens() ([]Token, error) {
	errs := make([]error, 0)
	sc.start = sc.current
	for !sc.isAtEnd() {
		sc.start = sc.current
		if err := sc.scanToken(); err != nil {
			errs = append(errs, err)
		}
	}
	sc.tokens = append(sc.tokens, Token{TokenType: EOF})
	if len(errs) == 0 {
		return sc.tokens, nil
	}
	return sc.tokens, errors.Join(errs...)
}

func (sc *Scanner) scanToken() error {
	char := sc.advance()
	switch {
	case char == '(':
		sc.addToken(LEFT_PAREN)
	case char == ')':
		sc.addToken(RIGHT_PAREN)
	case char == '{':
		sc.addToken(LEFT_BRACE)
	case char == '}':
		sc.addToken(RIGHT_BRACE)
	case char == ',':
		sc.addToken(COMMA)
	case char == '.':
		sc.addToken(DOT)
	case char == '-':
		sc.addToken(MINUS)
	case char == '+':
		sc.addToken(PLUS)
	case char == ';':
		sc.addToken(SEMICOLON)
	case char == '*':
		sc.addToken(STAR)
	case char == '!' && sc.match('='):
		sc.addToken(BANG_EQUAL)
	case char == '!':
		sc.addToken(BANG)
	case char == '=' && sc.match('='):
		sc.addToken(EQUAL_EQUAL)
	case char == '=':
		sc.addToken(EQUAL)
	case char == '<' && sc.match('='):
		sc.addToken(LESS_EQUAL)
	case char == '<':
		sc.addToken(LESS)
	case char == '>' && sc.match('='):
		sc.addToken(GREATER_EQUAL)
	case char == '>':
		sc.addToken(GREATER)
	case char == '/' && sc.match('/'):
		// Discard a comment
		for sc.peek() != '\n' && !sc.isAtEnd() {
			sc.advance()
		}
	case char == '/':
		sc.addToken(SLASH)
	case char == ' ' || char == '\r' || char == '\t':
		// Ignore whitespace
	case char == '\n':
		sc.line++
	case char == '"':
		return sc.addStringToken()
	case isDigit(char):
		return sc.addNumberToken()
	case isAlpha(char):
		return sc.addIdentifier()
	default:
		return &Error{sc.line, "Unexpected character."}
	}
	// unreachable code
	return nil
}

func (sc *Scanner) addToken(tokenType TokenType) {
	sc.addTokenLiteral(tokenType, nil)
}

func (sc *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
	lexeme := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, Token{TokenType: tokenType, Lexeme: &lexeme, Literal: literal, Line: sc.line})
}

func (sc *Scanner) addIdentifier() error {
	for isAlphaNum(sc.peek()) {
		sc.advance()
	}

	text := sc.source[sc.start:sc.current]
	if t, exist := keywords[text]; exist {
		sc.addToken(t)
	} else {
		sc.addToken(IDENTIFIER)
	}

	return nil
}

func (sc *Scanner) addNumberToken() error {
	for isDigit(sc.peek()) {
		sc.advance()
	}
	// Number with a floating point
	if sc.peek() == '.' && isDigit(sc.peekNext()) {
		// Consume the "."
		sc.advance()
		for isDigit(sc.peek()) {
			sc.advance()
		}
	}
	value, err := strconv.ParseFloat(sc.source[sc.start:sc.current], 64)
	if err != nil {
		return &Error{sc.line, "Unexpected character."}
	}
	sc.addTokenLiteral(NUMBER, value)
	return nil
}

func (sc *Scanner) addStringToken() error {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()
	}
	if sc.isAtEnd() {
		return &Error{sc.line, "Unterminated string."}
	}
	// The closing " .
	sc.advance()

	value := sc.source[sc.start+1 : sc.current-1] //
	sc.addTokenLiteral(STRING, value)

	return nil
}

// Return current character and move one character forward
func (sc *Scanner) advance() byte {
	ch := sc.source[sc.current]
	sc.current++
	return ch
}

// Checks if current character matches expected and advances in case it is true.
func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != expected {
		return false
	}
	sc.current++
	return true
}

// Return current character
func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return 0
	}
	return sc.source[sc.current]
}

// Return next character
func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.source) {
		return 0
	}
	return sc.source[sc.current+1]
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func isAlphaNum(char byte) bool {
	return isAlpha(char) || isDigit(char)
}

func isAlpha(char byte) bool {
	return char == '_' || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

type Error struct {
	line    int
	message string
}

func (e *Error) Error() string {
	return formatter.ReportError(e.line, "", e.message)
}
