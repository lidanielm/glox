package parser

import (
	"fmt"
	"testing"

	"github.com/lidanielm/glox/src/pkg/token"
	"github.com/lidanielm/glox/src/pkg/internal/ast"
)

func TestAstPrinter(t *testing.T) {
	expr := ast.NewBinary(
		ast.NewUnary(
			*token.NewToken(token.MINUS, "-", nil, 1),
			ast.NewLiteral(123),
		),
		*token.NewToken(token.STAR, "*", nil, 1),
		ast.NewGrouping(
			ast.NewLiteral(45.67),
		),
	)

	printer := ast.NewAstPrinter()
	result := printer.Print(expr)

	fmt.Println(result)

	expected := "(* (- 123) (group 45.67))"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}