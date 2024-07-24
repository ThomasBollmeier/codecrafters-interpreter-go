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

func (s *Scanner) AdvanceToken() (TokenIntf, error) {
	if s.done {
		return nil, NewScannerError("no more tokens")
	}

	for {

		cInfo, err := s.skipWhitespace()
		if err != nil {
			s.done = true
			return newToken(EOF, "", 0, 0), nil
		}

		tokenType, ok := singleCharTokenTypes[cInfo.char]
		if ok {
			return newToken(
				tokenType,
				cInfo.char,
				cInfo.line,
				cInfo.column,
			), nil
		}

		switch cInfo.char {
		case "=":
			return s.scanEqual(cInfo), nil
		case "!":
			return s.scanBang(cInfo), nil
		case "<", ">":
			return s.scanRelationOp(cInfo), nil
		case "/":
			token := s.scanSlash(cInfo)
			if token != nil {
				return token, nil
			} else {
				continue
			}
		case "\"":
			return s.scanString(cInfo), nil
		}

		return newErrorToken(
			cInfo.char,
			fmt.Sprintf("Unexpected character: %s", cInfo.char),
			cInfo.line,
			cInfo.column,
		), nil

	}
}

func (s *Scanner) scanString(cInfo charInfo) TokenIntf {
	lexeme := cInfo.char
	for {
		ci, err := s.advanceChar()
		if err != nil {
			return newErrorToken(
				lexeme,
				"Unterminated string.",
				cInfo.line,
				cInfo.column,
			)
		}
		lexeme += ci.char
		if ci.char == "\"" {
			return newToken(
				String,
				lexeme,
				cInfo.line,
				cInfo.column,
			)
		}
	}
}

func (s *Scanner) scanSlash(cInfo charInfo) *Token {
	nextChar := s.peekChar()

	if nextChar != "/" {
		return newToken(
			Slash,
			cInfo.char,
			cInfo.line,
			cInfo.column,
		)
	} else { // a line comment
		for {
			cInfo, err := s.advanceChar()
			if err != nil || cInfo.char == "\n" {
				break
			}
		}
		return nil
	}
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

	return newToken(
		tokenType,
		lexeme,
		cInfo.line,
		cInfo.column,
	)
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
	return newToken(
		tokenType,
		lexeme,
		cInfo.line,
		cInfo.column,
	)

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
	return newToken(
		tokenType,
		lexeme,
		cInfo.line,
		cInfo.column,
	)
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
