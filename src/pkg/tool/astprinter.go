package tool

import (
	"fmt"

	"github.com/lidanielm/glox/src/pkg/internal/ast"
)

type AstPrinter struct {}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a AstPrinter) Print(expr ast.Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitBinaryExpr(expr ast.Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a AstPrinter) VisitGroupingExpr(expr ast.Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a AstPrinter) VisitLiteralExpr(expr ast.Literal) any {
	return fmt.Sprintf("%v", expr.Value)
}

func (a AstPrinter) VisitUnaryExpr(expr ast.Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a AstPrinter) VisitTernaryExpr(expr ast.Ternary) any {
	return a.parenthesize(expr.Operator1.Lexeme, expr.Operator2.Lexeme, expr.Condition, expr.Left, expr.Right)
}

func (a AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	str := "(" + name
	for _, expr := range exprs {
		str += " "
		str += expr.Accept(a).(string)
	}
	str += ")"

	return str
}