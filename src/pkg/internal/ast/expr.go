package ast

import (
	"github.com/lidanielm/glox/src/pkg/token"
)

type Expr interface {
	Accept(visitor Visitor[any]) any
}

type Visitor[R any] interface {
	VisitBinaryExpr(expr Binary) R
	VisitGroupingExpr(expr Grouping) R
	VisitLiteralExpr(expr Literal) R
	VisitUnaryExpr(expr Unary) R
}

type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinary(Left Expr, Operator token.Token, Right Expr) *Binary {
	return &Binary{Left: Left, Operator: Operator, Right: Right}
}

func (b Binary) Accept(visitor Visitor[any]) any {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}

func (g Grouping) Accept(visitor Visitor[any]) any {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{Value: value}
}

func (l Literal) Accept(visitor Visitor[any]) any {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	Operator token.Token
	Right Expr
}

func NewUnary(Operator token.Token, Right Expr) *Unary {
	return &Unary{Operator: Operator, Right: Right}
}

func (u Unary) Accept(visitor Visitor[any]) any {
	return visitor.VisitUnaryExpr(u)
}

