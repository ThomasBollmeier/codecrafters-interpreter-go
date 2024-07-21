package main

import "fmt"

type TokenType string

const (
	LeftParen  TokenType = "LEFT_PAREN"
	RightParen TokenType = "RIGHT_PAREN"
	EOF        TokenType = "EOF"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s null", t.Type, t.Lexeme)
}
