package main

import "fmt"

type TokenType string

const (
	LeftParen  TokenType = "LEFT_PAREN"
	RightParen TokenType = "RIGHT_PAREN"
	LeftBrace  TokenType = "LEFT_BRACE"
	RightBrace TokenType = "RIGHT_BRACE"
	Plus       TokenType = "PLUS"
	Minus      TokenType = "MINUS"
	Star       TokenType = "STAR"
	Dot        TokenType = "DOT"
	Comma      TokenType = "COMMA"
	Semicolon  TokenType = "SEMICOLON"
	EOF        TokenType = "EOF"
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

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s null", t.Type, t.Lexeme)
}
