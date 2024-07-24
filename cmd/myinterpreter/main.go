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

	if command != "tokenize" {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	tokenize(os.Args[2])

}

func tokenize(filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := NewScanner(string(fileContents))
	var tokens []TokenIntf
	var errorTokens []TokenIntf

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
