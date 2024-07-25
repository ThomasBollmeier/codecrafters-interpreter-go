package main

import (
	"fmt"
	"unicode"
)

type Scanner struct {
	characters []rune
	index      int
	line       int
	column     int
	done       bool
}

type charInfo struct {
	char   rune
	line   int
	column int
}

func NewScanner(content string) *Scanner {
	var characters []rune
	for _, char := range content {
		characters = append(characters, char)
	}
	return &Scanner{
		characters: characters,
		index:      0,
		line:       1,
		column:     1,
		done:       false,
	}
}

func (s *Scanner) AdvanceToken() (TokenInfo, error) {
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
				string(cInfo.char),
				cInfo.line,
				cInfo.column,
			), nil
		}

		switch cInfo.char {
		case '=':
			return s.scanEqual(cInfo), nil
		case '!':
			return s.scanBang(cInfo), nil
		case '<', '>':
			return s.scanRelationOp(cInfo), nil
		case '/':
			token := s.scanSlash(cInfo)
			if token != nil {
				return token, nil
			} else {
				continue
			}
		case '"':
			return s.scanString(cInfo), nil
		}

		if unicode.IsDigit(cInfo.char) {
			return s.scanNumber(cInfo), nil
		}

		return newErrorToken(
			string(cInfo.char),
			fmt.Sprintf("Unexpected character: %c", cInfo.char),
			cInfo.line,
			cInfo.column,
		), nil

	}
}

func (s *Scanner) scanNumber(cInfo charInfo) TokenInfo {
	lexeme := string(cInfo.char)
	foundDot := false
loop:
	for {
		nextChars := s.peekNChars(2)
		switch len(nextChars) {
		case 0:
			break loop
		case 1:
			ch := nextChars[0]
			if unicode.IsDigit(ch) {
				lexeme += string(ch)
				_, _ = s.advanceChar()
			}
			break loop
		case 2:
			ch := nextChars[0]
			if unicode.IsDigit(ch) {
				lexeme += string(ch)
				_, _ = s.advanceChar()
			} else if ch == '.' && !foundDot {
				foundDot = true
				ch2 := nextChars[1]
				if unicode.IsDigit(ch2) {
					lexeme += string(ch)
					_, _ = s.advanceChar()
				} else {
					break loop
				}
			} else {
				break loop
			}
		}
	}

	return newToken(
		Number,
		lexeme,
		cInfo.line,
		cInfo.column,
	)
}

func (s *Scanner) scanString(cInfo charInfo) TokenInfo {
	lexeme := string(cInfo.char)
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
		lexeme += string(ci.char)
		if ci.char == '"' {
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
	nextChar, err := s.peekChar()

	if nextChar != '/' || err != nil {
		return newToken(
			Slash,
			string(cInfo.char),
			cInfo.line,
			cInfo.column,
		)
	} else { // a line comment
		for {
			cInfo, err = s.advanceChar()
			if err != nil || cInfo.char == '\n' {
				break
			}
		}
		return nil
	}
}

func (s *Scanner) scanRelationOp(cInfo charInfo) *Token {
	var tokenType TokenType
	var lexeme string
	nextChar, err := s.peekChar()

	if nextChar != '=' || err != nil {
		if cInfo.char == '>' {
			tokenType = Greater
		} else {
			tokenType = Less
		}
		lexeme = string(cInfo.char)
	} else {
		_, _ = s.advanceChar()
		if cInfo.char == '>' {
			tokenType = GreaterEqual
		} else {
			tokenType = LessEqual
		}
		lexeme = string(cInfo.char) + string(nextChar)
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
	nextChar, err := s.peekChar()
	if nextChar != '=' || err != nil {
		tokenType = Bang
		lexeme = string(cInfo.char)
	} else {
		_, _ = s.advanceChar()
		tokenType = BangEqual
		lexeme = string(cInfo.char) + string(nextChar)
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
	nextChar, err := s.peekChar()
	if nextChar != '=' || err != nil {
		tokenType = Equal
		lexeme = string(cInfo.char)
	} else {
		_, _ = s.advanceChar()
		tokenType = EqualEqual
		lexeme = string(cInfo.char) + string(nextChar)
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
	if char != '\n' {
		s.column++
	} else {
		s.line++
		s.column = 1
	}
	s.index++
	return charInfo{char, line, column}, nil
}

func (s *Scanner) peekChar() (rune, error) {
	if s.index >= len(s.characters) {
		return ' ', NewScannerError("no more characters")
	}
	return s.characters[s.index], nil
}

func (s *Scanner) peekNChars(n int) []rune {
	exclIndex := s.index + n
	if exclIndex < len(s.characters) {
		return s.characters[s.index:exclIndex]
	} else {
		return s.characters[s.index:]
	}
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
		case ' ', '\t', '\n', '\r':
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
