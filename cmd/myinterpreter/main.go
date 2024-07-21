package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	tokenize(os.Args[2])

}

func tokenize(filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := NewScanner(string(fileContents))
	var tokens []*Token
	var errorTokens []*Token

	for {
		token, err := scanner.AdvanceToken()
		if err != nil {
			break
		}
		if token.Type != Error {
			tokens = append(tokens, token)
		} else {
			errorTokens = append(errorTokens, token)
		}
	}

	for _, token := range errorTokens {
		fmt.Fprintf(
			os.Stderr,
			"[line %d] Error: Unexpected character: %s\n",
			token.Line,
			token.Lexeme)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	if errorTokens != nil {
		os.Exit(65)
	}
}
