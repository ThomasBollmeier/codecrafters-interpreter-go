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
	scanner := NewScanner("  (()\n")
	var tokens []Token
	for {
		token, err := scanner.AdvanceToken()
		if err != nil {
			break
		}
		tokens = append(tokens, *token)
	}

	assertEq(4, len(tokens), t)
	assertEq(LeftParen, tokens[0].Type, t)
	assertEq(LeftParen, tokens[1].Type, t)
	assertEq(RightParen, tokens[2].Type, t)
	assertEq(EOF, tokens[3].Type, t)

}

func assertEq(expected any, actual any, t *testing.T) {
	if expected != actual {
		t.Fatalf("expected: %v, actual: %v", expected, actual)
	}
}
