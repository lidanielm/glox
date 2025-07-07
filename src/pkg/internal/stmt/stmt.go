package stmt

import (
	"github.com/lidanielm/glox/src/pkg/internal/ast"
)

type Stmt interface {
	Accept(visitor Visitor[any]) error
}

type Visitor[R any] interface {
	VisitExpressionStmt(stmt Expression) error
	VisitPrintStmt(stmt Print) error
}

type Expression struct {
	Expr ast.Expr
}

func NewExpression(expr ast.Expr) *Expression {
	return &Expression{Expr: expr}
}

func (e Expression) Accept(visitor Visitor[any]) error {
	return visitor.VisitExpressionStmt(e)
}

type Print struct {
	Expr ast.Expr
}

func NewPrint(expr ast.Expr) *Print {
	return &Print{Expr: expr}
}

func (p Print) Accept(visitor Visitor[any]) error {
	return visitor.VisitPrintStmt(p)
}

