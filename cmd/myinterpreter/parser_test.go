package main

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	code := "40 + 2"
	parser := NewParser(code)

	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser.Parse() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}
