package scanner

import "testing"

func TestScanTokens(t *testing.T) {
	for _, test := range []struct {
		str   string
		token Token
	}{
		{"(", Token{LEFT_PAREN, getStrPtr("("), nil, 1}},
		{")", Token{RIGHT_PAREN, getStrPtr(")"), nil, 1}},
		{"{", Token{LEFT_BRACE, getStrPtr("{"), nil, 1}},
		{"}", Token{RIGHT_BRACE, getStrPtr("}"), nil, 1}},
		{",", Token{COMMA, getStrPtr(","), nil, 1}},
		{".", Token{DOT, getStrPtr("."), nil, 1}},
		{"-", Token{MINUS, getStrPtr("-"), nil, 1}},
		{"+", Token{PLUS, getStrPtr("+"), nil, 1}},
		{";", Token{SEMICOLON, getStrPtr(";"), nil, 1}},
		{"*", Token{STAR, getStrPtr("*"), nil, 1}},
		{"/", Token{SLASH, getStrPtr("/"), nil, 1}},
		{"!", Token{BANG, getStrPtr("!"), nil, 1}},
		{"!=", Token{BANG_EQUAL, getStrPtr("!="), nil, 1}},
		{"=", Token{EQUAL, getStrPtr("="), nil, 1}},
		{"==", Token{EQUAL_EQUAL, getStrPtr("=="), nil, 1}},
		{"<", Token{LESS, getStrPtr("<"), nil, 1}},
		{"<=", Token{LESS_EQUAL, getStrPtr("<="), nil, 1}},
		{">", Token{GREATER, getStrPtr(">"), nil, 1}},
		{">=", Token{GREATER_EQUAL, getStrPtr(">="), nil, 1}},
		{"123", Token{NUMBER, getStrPtr("123"), 123, 1}},
		{"\"123\"", Token{STRING, getStrPtr("\"123\""), "123", 1}},
		{"\"abc\"", Token{STRING, getStrPtr("\"abc\""), "abc", 1}},
		{"and", Token{AND, getStrPtr("and"), nil, 1}},
		{"class", Token{CLASS, getStrPtr("class"), nil, 1}},
		{"else", Token{ELSE, getStrPtr("else"), nil, 1}},
		{"false", Token{FALSE, getStrPtr("false"), nil, 1}},
		{"for", Token{FOR, getStrPtr("for"), nil, 1}},
		{"fun", Token{FUN, getStrPtr("fun"), nil, 1}},
		{"if", Token{IF, getStrPtr("if"), nil, 1}},
		{"nil", Token{NIL, getStrPtr("nil"), nil, 1}},
		{"or", Token{OR, getStrPtr("or"), nil, 1}},
		{"print", Token{PRINT, getStrPtr("print"), nil, 1}},
		{"return", Token{RETURN, getStrPtr("return"), nil, 1}},
		{"super", Token{SUPER, getStrPtr("super"), nil, 1}},
		{"this", Token{THIS, getStrPtr("this"), nil, 1}},
		{"true", Token{TRUE, getStrPtr("true"), nil, 1}},
		{"var", Token{VAR, getStrPtr("var"), nil, 1}},
		{"while", Token{WHILE, getStrPtr("while"), nil, 1}},
	} {
		sc := NewScanner(test.str)
		got, err := sc.ScanTokens()
		if err != nil {
			t.Error(err)
		}
		if len(got) != 2 {
			t.Fatalf("expect len(1), got: %d\n", len(got))
		}
		gotToken := got[0]
		if gotToken.TokenType != test.token.TokenType {
			t.Fatalf("expect: %v, got: %v, string: %v\n", test.token.TokenType, gotToken.TokenType, test.str)
		}
		if gotToken.Line != test.token.Line {
			t.Fatalf("expect: %v, got: %v, string: %v\n", test.token.Line, gotToken.Line, test.str)
		}
		if gotToken.Lexeme == nil && test.token.Lexeme != nil || gotToken.Lexeme != nil && test.token.Lexeme == nil || *gotToken.Lexeme != *test.token.Lexeme {
			if gotToken.Lexeme != nil && test.token.Lexeme != nil {
				t.Fatalf("expect: %v, got: %v, string: %v\n", *test.token.Lexeme, *gotToken.Lexeme, test.str)
			} else {
				t.Fatalf("expect: %v, got: %v, string: %v\n", test.token.Lexeme, gotToken.Lexeme, test.str)
			}
		}
		// TODO : add compare literals

		eofToken := got[1]
		if eofToken.TokenType != EOF {
			t.Fatalf("expect: %v, got: %v, string: %v\n", EOF, gotToken.TokenType, test.str)
		}
	}
}

func TestScanTokensIgnored(t *testing.T) {
	for _, test := range []string{
		"// this is a comment", " ", "\r", "\t", "\n",
	} {
		sc := NewScanner(test)
		got, err := sc.ScanTokens()
		if err != nil {
			t.Error(err)
		}
		if len(got) != 1 {
			t.Fatalf("expect len(1), got: %d\n", len(got))
		}
		gotToken := got[0]
		if gotToken.TokenType != EOF {
			t.Fatalf("expect: %v, got: %v\n", EOF, gotToken.TokenType)
		}
	}
}

func TestScanTokensError(t *testing.T) {
	for _, test := range []struct {
		text      string
		expectErr string
	}{
		{"\"missing quote", "[line 1] Error: Unterminated string."},
		{"~", "[line 1] Error: Unexpected character."},
	} {
		sc := NewScanner(test.text)
		got, err := sc.ScanTokens()
		if err == nil {
			t.Fatalf("expect err\n")
		}
		if err.Error() != test.expectErr {
			t.Fatalf("expect: %v, got: %v\n", test.expectErr, err.Error())
		}
		if len(got) != 1 {
			t.Fatalf("expect len(1), got: %d\n", len(got))
		}
	}
}

func TestScanTokenErrors(t *testing.T) {
	for _, test := range []struct {
		text      string
		expectErr string
	}{
		{"@ > \n\"2", "[line 1] Error: Unexpected character.\n[line 2] Error: Unterminated string."},
	} {
		sc := NewScanner(test.text)
		got, err := sc.ScanTokens()
		if err == nil {
			t.Fatalf("expect err\n")
		}
		if err.Error() != test.expectErr {
			t.Fatalf("expect: %v, got: %v\n", test.expectErr, err.Error())
		}
		if len(got) != 2 {
			t.Fatalf("expect len(1), got: %d\n", len(got))
		}
	}
}

func getStrPtr(s string) *string {
	return &s
}
