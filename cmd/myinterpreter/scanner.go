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

	char, err := s.skipWhitespace()
	if err != nil {
		s.done = true
		return &Token{Type: EOF}, nil
	}

	tokenType, ok := singleCharTokenTypes[char]
	if ok {
		return &Token{Type: tokenType, Lexeme: char}, nil
	}

	return nil, NewScannerError("unexpected character %s", char)
}

func (s *Scanner) advanceChar() (string, error) {
	if s.index >= len(s.characters) {
		return "", NewScannerError("no more characters")
	}
	char := s.characters[s.index]
	s.index++
	return char, nil
}

func (s *Scanner) skipWhitespace() (string, error) {
	for {
		char, err := s.advanceChar()
		if err != nil {
			return "", err
		}
		switch char {
		case " ", "\t", "\n", "\r":
			continue
		default:
			return char, nil
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
