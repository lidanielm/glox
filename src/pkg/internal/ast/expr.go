package ast

type expr interface {
	Accept(visitor Visitor[any]) any
}

type Binary struct {
	left Expr
	operator Token
	right Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{left: Expr, operator: Token, right: Expr}
}

func (b Binary) Accept(visitor Visitor[any]) any {
	return visitor.visitBinaryExpr(b)
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression: Expr}
}

func (g Grouping) Accept(visitor Visitor[any]) any {
	return visitor.visitGroupingExpr(g)
}

type Literal struct {
	value Object
}

func NewLiteral(value Object) *Literal {
	return &Literal{value: Object}
}

func (l Literal) Accept(visitor Visitor[any]) any {
	return visitor.visitLiteralExpr(l)
}

type Unary struct {
	operator Token
	right Expr
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{operator: Token, right: Expr}
}

func (u Unary) Accept(visitor Visitor[any]) any {
	return visitor.visitUnaryExpr(u)
}

