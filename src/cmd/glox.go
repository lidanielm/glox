package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lidanielm/glox/src/pkg/scanner"
)

type Interpreter struct {
	hadError bool
}

func (iptr *Interpreter) main() {
	if len(os.Args) > 1 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 1 {
		iptr.runFile(os.Args[1])
	} else {
		iptr.runPrompt()
	}
}

func (iptr *Interpreter) runFile(path string) error {
	// Wrapper for run if given file path
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	iptr.run(string(data))

	if iptr.hadError {
		os.Exit(64)
	}
	return nil
}

func (iptr *Interpreter) runPrompt() error {
	// Wrapper for run in repl environment
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _, _ := reader.ReadLine()
		if line == nil {
			break
		}
		iptr.run(string(line))
	}
	return nil
}

func (iptr *Interpreter) run(source string) error {
	// Run interpreter
	scan := scanner.NewScanner(source)
	tokens, err := scan.ScanTokens()

	// TODO: Handle custom error
	if err != nil {
		return err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return nil
}