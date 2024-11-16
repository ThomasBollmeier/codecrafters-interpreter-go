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

func TestParser_ParseLogicalExpression(t *testing.T) {
	code := "a or b and a == \"hello\""
	parser := NewParser(code)

	ast, err := parser.ParseExpression()
	if err != nil {
		t.Fatalf("parser.ParseExpression() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}

func TestParser_ParseForStmt(t *testing.T) {
	code := `for (var foo = 0; foo < 3;) 
		print foo = foo + 1;`
	parser := NewParser(code)

	ast, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("parser.ParseExpression() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}

func TestParser_ParseReturn(t *testing.T) {
	code := `
		fun fib(n) {
			if (n < 2) return n;
			return fib(n - 2) + fib(n - 1);
		}
		var start = clock();
		print fib(10) == 55;
		print (clock() - start) < 5; // 5 seconds`

	parser := NewParser(code)

	ast, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("parser.ParseExpression() error = %v", err)
	}

	ast.accept(NewAstPrinter())
}
