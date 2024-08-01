package main

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	code := "72 * 63 / 48"
	parser := NewParser(code)

	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser.Parse() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}
