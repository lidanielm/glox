package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lidanielm/glox/src/pkg/scanner"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/parser"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) error {
	// Wrapper for run if given file path
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	run(string(data))
	return nil
}

func runPrompt() error {
	// Wrapper for run in repl environment
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _, _ := reader.ReadLine()
		if line == nil {
			break
		}
		run(string(line))
	}
	return nil
}

func run(source string) *lox_error.Error {
	// Run interpreter
	scan := scanner.NewScanner(source)
	tokens, err := scan.ScanTokens()

	parser := parser.NewParser(tokens)
	expr, err := parser.Parse()

	// Stop if there was a syntax error
	if err != nil {
		return err
	}

	fmt.Println(expr)
	
	return nil
}
