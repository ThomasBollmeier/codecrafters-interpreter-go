package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "tokenize":
		tokenize(os.Args[2])
	case "parse":
		parse(os.Args[2])
	case "evaluate":
		evaluate(os.Args[2])
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

}

func evaluate(filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	interpreter := NewInterpreter(string(fileContents))
	value, err, isRuntimeError := interpreter.Eval()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error evaluating file: %v\n", err)
		if isRuntimeError {
			os.Exit(70)
		} else {
			os.Exit(65)
		}
	}

	fmt.Printf("%s\n", value)
}

func parse(filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	parser := NewParser(string(fileContents))
	ast, err := parser.Parse()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(65)
	}

	astPrinter := NewAstPrinter()
	ast.accept(astPrinter)
}

func tokenize(filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := NewScanner(string(fileContents))
	var tokens []TokenInfo
	var errorTokens []TokenInfo

	for {
		token, err := scanner.AdvanceToken()
		if err != nil {
			break
		}
		if token.GetTokenType() != Error {
			tokens = append(tokens, token)
		} else {
			errorTokens = append(errorTokens, token)
		}
	}

	for _, token := range errorTokens {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"%s\n",
			token)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	if errorTokens != nil {
		os.Exit(65)
	}
}
