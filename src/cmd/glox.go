package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lidanielm/glox/src/pkg/interpreter"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/parser"
	"github.com/lidanielm/glox/src/pkg/scanner"
	"github.com/lidanielm/glox/src/pkg/tool"
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
		err := run(string(line))
		if err != nil {
			if _, ok := err.(*lox_error.RuntimeError); ok {
				os.Exit(70)
			} else if _, ok := err.(*lox_error.LoxError); ok {
				os.Exit(65)
			} else {
				os.Exit(1)
			}
		}
	}
	return nil
}

func run(source string) error {
	// Run interpreter
	scan := scanner.NewScanner(source)
	tokens, err := scan.ScanTokens()

	parser := parser.NewParser(tokens)
	expr, err := parser.Parse()

	// Stop if there was a syntax error
	if err != nil {
		return err
	}

	interpreter := interpreter.NewInterpreter()
	interpreter.Interpret(expr)

	fmt.Println(tool.NewAstPrinter().Print(expr))
	
	return nil
}