package main

import (
	"github.com/lidanielm/glox/src/pkg/tool"
)

func main() {
	outputDir := "src/pkg/internal/ast"
	baseName := "Expr"
	types := []string{
		"Binary : Left Expr, Operator token.Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal : Value interface{}",
		"Unary : Operator token.Token, Right Expr",
		"Ternary : Condition Expr, Operator1 token.Token, Left Expr, Operator2 token.Token, Right Expr",
		"Variable : Name token.Token",
	}
	err := tool.DefineAst(outputDir, baseName, types)
	if err != nil {
		panic(err)
	}

	// outputDir := "src/pkg/internal/stmt"
	// baseName := "Stmt"
	// types := []string{"Expression : Expression Expr", "Print : Expression Expr"}
	// err := tool.DefineAst(outputDir, baseName, types)
	// if err != nil {
	// 	panic(err)
	// }
}