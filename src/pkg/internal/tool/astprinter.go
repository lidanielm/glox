package tool

import (
	"fmt"

	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
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
	return a.parenthesize("var "+expr.Name.Lexeme+" = ", expr.Value)
}

func (a AstPrinter) VisitLogicalExpr(expr ast.Logical) (any, error) {
	switch expr.Operator.Type {
	case token.OR:
		return a.parenthesize("or", expr.Left, expr.Right)
	case token.AND:
		return a.parenthesize("and", expr.Left, expr.Right)
	default:
		return nil, lox_error.NewError(expr.Operator, "Unrecognizable logical expression.")
	}
}

func (a AstPrinter) VisitCallExpr(expr ast.Call) (any, error) {
	return a.parenthesize("call "+expr.Paren.Lexeme, expr.Arguments...)
}

func (a AstPrinter) VisitGetExpr(expr ast.Get) (any, error) {
	return a.parenthesize("get "+expr.Name.Lexeme, expr.Object)
}

func (a AstPrinter) VisitSetExpr(expr ast.Set) (any, error) {
	return a.parenthesize("set "+expr.Name.Lexeme, expr.Object, expr.Value)
}

func (a AstPrinter) VisitThisExpr(expr ast.This) (any, error) {
	return fmt.Sprintf("this.%s", expr.Keyword.Lexeme), nil
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