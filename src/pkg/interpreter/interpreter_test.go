package interpreter

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/lidanielm/glox/src/pkg/internal/ast"
// 	"github.com/lidanielm/glox/src/pkg/internal/stmt"
// 	"github.com/lidanielm/glox/src/pkg/internal/tool"
// 	"github.com/lidanielm/glox/src/pkg/parser"
// 	"github.com/lidanielm/glox/src/pkg/scanner"
// )

// func TestInterpreter(t *testing.T) {
// 	ip := NewInterpreter()
// 	script := "var x = 42; print x;"
// 	scanner := scanner.NewScanner(script)
// 	tokens, err := scanner.ScanTokens()
// 	if err != nil {
// 		t.Fatalf("Failed to scan tokens: %v", err)
// 	}

// 	// Parse the statements
// 	parser := parser.NewParser(tokens)
// 	statements, err := parser.Parse()
// 	if err != nil {
// 		t.Fatalf("Failed to parse statements: %v", err)
// 	}

// 	err = ip.Interpret([]stmt.Stmt{statements[0]})
// 	if err != nil {
// 		t.Fatalf("Failed to interpret statements: %v", err)
// 	}

// 	varStmt, ok := statements[0].(*stmt.Var)
// 	if !ok {
// 		t.Fatalf("First statement should be a variable declaration")
// 	}

// 	literal, ok := varStmt.Initializer.(*ast.Literal)
// 	if !ok {
// 		t.Fatalf("Variable initializer should be a literal, got %T", varStmt.Initializer)
// 	}

// 	value, err := ip.env.Get(varStmt.Name)
// 	if err != nil {
// 		t.Fatalf("First statement not stored in environment: %v", err)
// 	}

// 	fmt.Printf("Successfully parsed variable declaration")
// 	fmt.Printf("Variable declaration: var %s = %v\n", varStmt.Name.Lexeme, literal.Value)
// 	fmt.Printf("Environment value: map[%s] = %v\n", varStmt.Name.Lexeme, value)

// 	printStmt, ok := statements[1].(*stmt.Print)
// 	if !ok {
// 		t.Fatalf("Second statement should be a print statement")
// 	}

// 	_, err = ip.env.Get(varStmt.Name)
// 	if err != nil {
// 		t.Fatalf("Variable not stored in environment: %v", err)
// 	}

// 	expr := printStmt.Expr

// 	printer := tool.NewAstPrinter()
// 	result := printer.Print(expr)

// 	fmt.Printf("Print statement expression: %s\n", result)

// 	value, err = ip.evaluate(expr)
// 	if err != nil {
// 		t.Fatalf("Variable not stored in environment: %v", err)
// 	}

// 	fmt.Printf("Value: %v", value)
// }