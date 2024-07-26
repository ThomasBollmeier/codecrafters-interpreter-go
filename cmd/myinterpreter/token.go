package main

import (
	"fmt"
	"math"
	"strconv"
)

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
	Number       TokenType = "NUMBER"
	Identifier   TokenType = "IDENTIFIER"
	And          TokenType = "AND"
	Class        TokenType = "CLASS"
	Else         TokenType = "ELSE"
	False        TokenType = "FALSE"
	For          TokenType = "FOR"
	Fun          TokenType = "FUN"
	If           TokenType = "IF"
	Nil          TokenType = "NIL"
	Or           TokenType = "OR"
	Print        TokenType = "PRINT"
	Return       TokenType = "RETURN"
	Super        TokenType = "SUPER"
	This         TokenType = "THIS"
	True         TokenType = "TRUE"
	Var          TokenType = "VAR"
	While        TokenType = "WHILE"
	Error        TokenType = "ERROR"
	EOF          TokenType = "EOF"
)

var reservedWords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

var singleCharTokenTypes = map[rune]TokenType{
	'(': LeftParen,
	')': RightParen,
	'{': LeftBrace,
	'}': RightBrace,
	'+': Plus,
	'-': Minus,
	'*': Star,
	'.': Dot,
	',': Comma,
	';': Semicolon,
}

type TokenInfo interface {
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
	case Number:
		floatValue, err := strconv.ParseFloat(t.lexeme, 64)
		if err != nil {
			floatValue = 0.0
		}
		if math.Abs(floatValue-math.Round(floatValue)) > 1e-5 {
			return fmt.Sprintf("%s %s %.4f", t.tokenType, t.lexeme, floatValue)
		} else {
			return fmt.Sprintf("%s %s %.1f", t.tokenType, t.lexeme, floatValue)
		}
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
