package main

import (
	"fmt"
)

type Scanner struct {
	characters []string
	index      int
	line       int
	column     int
	done       bool
}

type charInfo struct {
	char   string
	line   int
	column int
}

func NewScanner(content string) *Scanner {
	var characters []string
	for _, char := range content {
		characters = append(characters, string(char))
	}
	return &Scanner{
		characters: characters,
		index:      0,
		line:       1,
		column:     1,
		done:       false,
	}
}

func (s *Scanner) AdvanceToken() (*Token, error) {
	if s.done {
		return nil, NewScannerError("no more tokens")
	}

	cInfo, err := s.skipWhitespace()
	if err != nil {
		s.done = true
		return &Token{Type: EOF}, nil
	}

	tokenType, ok := singleCharTokenTypes[cInfo.char]
	if ok {
		return &Token{
			Type:   tokenType,
			Lexeme: cInfo.char,
			Line:   cInfo.line,
			Column: cInfo.column,
		}, nil
	}

	switch cInfo.char {
	case "=":
		return s.scanEqual(cInfo), nil
	case "!":
		return s.scanBang(cInfo), nil
	case "<", ">":
		return s.scanRelationOp(cInfo), nil
	}

	return &Token{
		Type:   Error,
		Lexeme: cInfo.char,
		Line:   cInfo.line,
		Column: cInfo.column,
	}, nil
}

func (s *Scanner) scanRelationOp(cInfo charInfo) *Token {
	var tokenType TokenType
	var lexeme string
	nextChar := s.peekChar()

	if nextChar != "=" {
		if cInfo.char == ">" {
			tokenType = Greater
		} else {
			tokenType = Less
		}
		lexeme = cInfo.char
	} else {
		_, _ = s.advanceChar()
		if cInfo.char == ">" {
			tokenType = GreaterEqual
		} else {
			tokenType = LessEqual
		}
		lexeme = cInfo.char + nextChar
	}

	return &Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   cInfo.line,
		Column: cInfo.column,
	}
}

func (s *Scanner) scanBang(cInfo charInfo) *Token {
	var tokenType TokenType
	var lexeme string
	nextChar := s.peekChar()
	if nextChar != "=" {
		tokenType = Bang
		lexeme = cInfo.char
	} else {
		_, _ = s.advanceChar()
		tokenType = BangEqual
		lexeme = cInfo.char + nextChar
	}
	return &Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   cInfo.line,
		Column: cInfo.column,
	}

}

func (s *Scanner) scanEqual(cInfo charInfo) *Token {
	var tokenType TokenType
	var lexeme string
	nextChar := s.peekChar()
	if nextChar != "=" {
		tokenType = Equal
		lexeme = cInfo.char
	} else {
		_, _ = s.advanceChar()
		tokenType = EqualEqual
		lexeme = cInfo.char + nextChar
	}
	return &Token{
		Type:   tokenType,
		Lexeme: lexeme,
		Line:   cInfo.line,
		Column: cInfo.column,
	}
}

func (s *Scanner) advanceChar() (charInfo, error) {
	if s.index >= len(s.characters) {
		return charInfo{}, NewScannerError("no more characters")
	}
	char := s.characters[s.index]
	line, column := s.line, s.column
	if char != "\n" {
		s.column++
	} else {
		s.line++
		s.column = 1
	}
	s.index++
	return charInfo{char, line, column}, nil
}

func (s *Scanner) peekChar() string {
	if s.index >= len(s.characters) {
		return ""
	}
	return s.characters[s.index]
}

func (s *Scanner) getLocation() (int, int) {
	return s.line, s.column
}

func (s *Scanner) skipWhitespace() (charInfo, error) {
	for {
		cInfo, err := s.advanceChar()
		if err != nil {
			return charInfo{}, err
		}
		switch cInfo.char {
		case " ", "\t", "\n", "\r":
			continue
		default:
			return cInfo, nil
		}
	}
}

type ScannerError struct {
	message string
}

func NewScannerError(format string, args ...any) *ScannerError {
	return &ScannerError{fmt.Sprintf(format, args...)}
}

func (e *ScannerError) Error() string {
	return e.message
}
