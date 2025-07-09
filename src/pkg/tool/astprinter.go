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
	s, err := expr.Accept(a)
	if err != nil {
		return "Unable to parse string"
	}
	return s.(string)
}

func (a AstPrinter) VisitBinaryExpr(expr ast.Binary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a AstPrinter) VisitGroupingExpr(expr ast.Grouping) (any, error) {
	return a.parenthesize("group", expr.Expression)
}

func (a AstPrinter) VisitLiteralExpr(expr ast.Literal) (any, error) {
	return fmt.Sprintf("%v", expr.Value), nil
}

func (a AstPrinter) VisitUnaryExpr(expr ast.Unary) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a AstPrinter) VisitTernaryExpr(expr ast.Ternary) (any, error) {
	return a.parenthesize(expr.Operator1.Lexeme, expr.Condition, expr.Left, expr.Right)
}

func (a AstPrinter) VisitVariableExpr(expr ast.Variable) (any, error) {
	return a.parenthesize("var "+expr.Name.Lexeme)
}

func (a AstPrinter) VisitAssignExpr(expr ast.Assign) (any, error) {
	return a.parenthesize("var"+expr.Name.Lexeme+" = ", expr.Value)
}

func (a AstPrinter) parenthesize(name string, exprs ...ast.Expr) (string, error) {
	str := "(" + name
	for _, expr := range exprs {
		str += " "
		s, err := expr.Accept(a)
		if err != nil {
			return "", err
		}
		str += s.(string)
	}
	str += ")"

	return str, nil
}