package parser

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/lidanielm/glox/src/pkg/internal/ast"
// 	"github.com/lidanielm/glox/src/pkg/internal/stmt"
// 	"github.com/lidanielm/glox/src/pkg/scanner"
// 	"github.com/lidanielm/glox/src/pkg/token"
// 	"github.com/lidanielm/glox/src/pkg/tool"
// )

// func TestArithmetic(t *testing.T) {
// 	expr := ast.NewBinary(
// 		ast.NewUnary(
// 			*token.NewToken(token.MINUS, "-", nil, 1),
// 			ast.NewLiteral(123),
// 		),
// 		*token.NewToken(token.STAR, "*", nil, 1),
// 		ast.NewGrouping(
// 			ast.NewLiteral(45.67),
// 		),
// 	)

// 	printer := tool.NewAstPrinter()
// 	result := printer.Print(expr)

// 	fmt.Println(result)

// 	expected := "(* (- 123) (group 45.67))"
// 	if result != expected {
// 		t.Errorf("Expected %q, got %q", expected, result)
// 	}
// }

// func TestPrint(t *testing.T) {
// 	// Test script: "var x = 42; print x;"
// 	script := "var x = 42; print x;"

// 	// Scan the tokens
// 	scanner := scanner.NewScanner(script)
// 	tokens, err := scanner.ScanTokens()
// 	if err != nil {
// 		t.Fatalf("Failed to scan tokens: %v", err)
// 	}

// 	// Parse the statements
// 	parser := NewParser(tokens)
// 	statements, err := parser.Parse()
// 	if err != nil {
// 		t.Fatalf("Failed to parse statements: %v", err)
// 	}

// 	// Verify we have exactly 2 statements
// 	if len(statements) != 2 {
// 		t.Fatalf("Expected 2 statements, got %d", len(statements))
// 	}

// 	// Test the first statement (variable declaration)
// 	varStmt, ok := statements[0].(*stmt.Var)
// 	if !ok {
// 		t.Fatalf("First statement should be a variable declaration")
// 	}

// 	if varStmt.Name.Lexeme != "x" {
// 		t.Errorf("Expected variable name 'x', got '%s'", varStmt.Name.Lexeme)
// 	}

// 	// Verify the initializer is a literal with value 42
// 	literal, ok := varStmt.Initializer.(*ast.Literal)
// 	if !ok {
// 		t.Fatalf("Variable initializer should be a literal, got %T", varStmt.Initializer)
// 	}

// 	if literal.Value != 42.0 {
// 		t.Errorf("Expected literal value 42, got %v", literal.Value)
// 	}

// 	// Test the second statement (print statement)
// 	printStmt, ok := statements[1].(*stmt.Print)
// 	if !ok {
// 		t.Fatalf("Second statement should be a print statement")
// 	}

// 	// Verify the print statement contains a variable reference
// 	variable, ok := printStmt.Expr.(*ast.Variable)
// 	if !ok {
// 		t.Fatalf("Print statement should contain a variable reference, got %T", printStmt.Expr)
// 	}

// 	if variable.Name.Lexeme != "x" {
// 		t.Errorf("Expected variable name 'x' in print statement, got '%s'", variable.Name.Lexeme)
// 	}

// 	// Verify the structure is correct
// 	fmt.Printf("Successfully parsed: var x = 42; print x;\n")
// 	fmt.Printf("Variable declaration: var %s = %v\n", varStmt.Name.Lexeme, literal.Value)
// 	fmt.Printf("Print statement: print %s\n", variable.Name.Lexeme)
// }