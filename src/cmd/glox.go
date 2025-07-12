package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lidanielm/glox/src/pkg/interpreter"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/parser"
	"github.com/lidanielm/glox/src/pkg/scanner"
)

func main() {
	// fmt.Println(os.Args)
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
	interpreter := interpreter.NewInterpreter()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	run(string(data), interpreter)
	return nil
}

func runPrompt() error {
	// Wrapper for run in repl environment
	reader := bufio.NewReader(os.Stdin)

	interpreter := interpreter.NewInterpreter()

	for {
		fmt.Print("> ")
		line, _, _ := reader.ReadLine()
		if line == nil {
			break
		}
		err := run(string(line), interpreter)
		if err != nil {
			if runtimeError, ok := err.(*lox_error.RuntimeError); ok {
				fmt.Println(runtimeError.Error())
				os.Exit(1)
			} else if loxError, ok := err.(*lox_error.LoxError); ok {
				fmt.Println(loxError.Error())
				os.Exit(1)
			} else if parseError, ok := err.(*lox_error.ParseError); ok {
				fmt.Println(parseError.Error())
				os.Exit(1)
			} else {
				os.Exit(1)
			}
		}
	}
	return nil
}

func run(source string, interpreter *interpreter.Interpreter) error {
	// Run interpreter
	scan := scanner.NewScanner(source)
	tokens, err := scan.ScanTokens()

	if err != nil {
		return err
	}

	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()

	// Stop if there was a syntax error
	if err != nil {
		return err
	}
	interpreter.Interpret(statements)
	
	return nil
}