package main

import (
	"testing"
)

func TestNewScanner(t *testing.T) {
	scanner := NewScanner("")
	if scanner == nil {
		t.Errorf("scanner should not be nil")
	}
}

func TestScanner_AdvanceToken(t *testing.T) {
	scanner := NewScanner("  ((){}\n")
	var tokens []TokenInfo
	for {
		token, err := scanner.AdvanceToken()
		if err != nil {
			break
		}
		tokens = append(tokens, token)
	}

	assertEq(6, len(tokens), t)
	assertEq(LeftParen, tokens[0].GetTokenType(), t)
	assertEq(LeftParen, tokens[1].GetTokenType(), t)
	assertEq(RightParen, tokens[2].GetTokenType(), t)
	assertEq(LeftBrace, tokens[3].GetTokenType(), t)
	assertEq(RightBrace, tokens[4].GetTokenType(), t)
	assertEq(EOF, tokens[5].GetTokenType(), t)

}

func assertEq(expected any, actual any, t *testing.T) {
	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}
