package main

import 
(
	"fmt",
	"bufio",
	"os"
)

type Interpreter struct {
	hadError bool
}

func (iptr *Interpreter) main() {
	if len(os.Args) > 1 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 1 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func (iptr *Interpreter) runFile(string path) error {
	// Wrapper for run if given file path
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	run(string(data))

	if (hadError) os.Exit(64)
	return nil
}

func (iptr *Interpreter) runPrompt() error {
	// Wrapper for run in repl environment
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _ := reader.ReadLine()
		if line == nil {
			break
		}
		run(line)
	}
	return nil
}

func (iptr *Interpreter) run(string source) error {
	// Run interpreter
	scanner := Scanner(source)
	tokens := scanner.scanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}