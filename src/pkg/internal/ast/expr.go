package ast

import (
	"github.com/lidanielm/glox/src/pkg/token"
)

type Expr interface {
	Accept(visitor Visitor[any]) (any, error)
}

type Visitor[R any] interface {
	VisitBinaryExpr(expr Binary) (R, error)
	VisitGroupingExpr(expr Grouping) (R, error)
	VisitLiteralExpr(expr Literal) (R, error)
	VisitUnaryExpr(expr Unary) (R, error)
	VisitTernaryExpr(expr Ternary) (R, error)
}

type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinary(Left Expr, Operator token.Token, Right Expr) *Binary {
	return &Binary{Left: Left, Operator: Operator, Right: Right}
}

func (b Binary) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) *Grouping {
	return &Grouping{Expression: Expression}
}

func (g Grouping) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value any
}

func NewLiteral(Value any) *Literal {
	return &Literal{Value: Value}
}

func (l Literal) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	Operator token.Token
	Right Expr
}

func NewUnary(Operator token.Token, Right Expr) *Unary {
	return &Unary{Operator: Operator, Right: Right}
}

func (u Unary) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitUnaryExpr(u)
}

type Ternary struct {
	Condition Expr
	Operator1 token.Token
	Left Expr
	Operator2 token.Token
	Right Expr
}

func NewTernary(Condition Expr, Operator1 token.Token, Left Expr, Operator2 token.Token, Right Expr) *Ternary {
	return &Ternary{Condition: Condition, Operator1: Operator1, Left: Left, Operator2: Operator2, Right: Right}
}

func (t Ternary) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitTernaryExpr(t)
}

