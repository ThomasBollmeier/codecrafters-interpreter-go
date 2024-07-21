package main

import "fmt"

type TokenType string

const (
	LeftParen  TokenType = "LEFT_PAREN"
	RightParen TokenType = "RIGHT_PAREN"
	LeftBrace  TokenType = "LEFT_BRACE"
	RightBrace TokenType = "RIGHT_BRACE"
	EOF        TokenType = "EOF"
)

var singleCharTokenTypes = map[string]TokenType{
	"(": TokenType(LeftParen),
	")": TokenType(RightParen),
	"{": TokenType(LeftBrace),
	"}": TokenType(RightBrace),
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
