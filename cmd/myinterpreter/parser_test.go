package main

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	code := "(68 - 11) >= -(17 / 54 + 34)"
	parser := NewParser(code)

	ast, err := parser.ParseExpression()
	if err != nil {
		t.Fatalf("parser.ParseExpression() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}

func TestParser_ParseAssignment(t *testing.T) {
	code := "a = b = 42"
	parser := NewParser(code)

	ast, err := parser.ParseExpression()
	if err != nil {
		t.Fatalf("parser.ParseExpression() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}
