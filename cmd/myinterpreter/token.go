package main

import "fmt"

type TokenType string

const (
	LeftParen    TokenType = "LEFT_PAREN"
	RightParen   TokenType = "RIGHT_PAREN"
	LeftBrace    TokenType = "LEFT_BRACE"
	RightBrace   TokenType = "RIGHT_BRACE"
	Plus         TokenType = "PLUS"
	Minus        TokenType = "MINUS"
	Star         TokenType = "STAR"
	Slash        TokenType = "SLASH"
	Dot          TokenType = "DOT"
	Comma        TokenType = "COMMA"
	Semicolon    TokenType = "SEMICOLON"
	Equal        TokenType = "EQUAL"
	EqualEqual   TokenType = "EQUAL_EQUAL"
	Bang         TokenType = "BANG"
	BangEqual    TokenType = "BANG_EQUAL"
	Less         TokenType = "LESS"
	LessEqual    TokenType = "LESS_EQUAL"
	Greater      TokenType = "GREATER"
	GreaterEqual TokenType = "GREATER_EQUAL"
	String       TokenType = "STRING"
	Error        TokenType = "ERROR"
	EOF          TokenType = "EOF"
)

var singleCharTokenTypes = map[string]TokenType{
	"(": LeftParen,
	")": RightParen,
	"{": LeftBrace,
	"}": RightBrace,
	"+": Plus,
	"-": Minus,
	"*": Star,
	".": Dot,
	",": Comma,
	";": Semicolon,
}

type TokenIntf interface {
	GetTokenType() TokenType
	GetLexeme() string
	GetPosition() (int, int)
}

type Token struct {
	tokenType TokenType
	lexeme    string
	line      int
	column    int
}

func newToken(tokenType TokenType, lexeme string, line int, column int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		line:      line,
		column:    column,
	}
}

func (t Token) String() string {
	switch t.tokenType {
	case String:
		length := len(t.lexeme)
		strValue := t.lexeme[1 : length-1]
		return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, strValue)
	default:
		return fmt.Sprintf("%s %s null", t.tokenType, t.lexeme)
	}

}

func (t Token) GetTokenType() TokenType {
	return t.tokenType
}

func (t Token) GetLexeme() string {
	return t.lexeme
}

func (t Token) GetPosition() (int, int) {
	return t.line, t.column
}

type ErrorToken struct {
	lexeme  string
	message string
	line    int
	column  int
}

func newErrorToken(lexeme string, message string, line int, column int) *ErrorToken {
	return &ErrorToken{
		lexeme:  lexeme,
		message: message,
		line:    line,
		column:  column,
	}
}

func (e ErrorToken) String() string {
	return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}

func (e ErrorToken) GetTokenType() TokenType {
	return Error
}

func (e ErrorToken) GetLexeme() string {
	return e.lexeme
}

func (e ErrorToken) GetPosition() (int, int) {
	return e.line, e.column
}
