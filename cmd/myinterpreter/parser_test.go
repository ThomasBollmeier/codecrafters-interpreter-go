package main

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	code := "(68 - 11) >= -(17 / 54 + 34)"
	parser := NewParser(code)

	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser.Parse() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}
