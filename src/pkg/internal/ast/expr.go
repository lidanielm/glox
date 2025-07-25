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
	VisitVariableExpr(expr Variable) (R, error)
	VisitAssignExpr(expr Assign) (R, error)
	VisitLogicalExpr(expr Logical) (R, error)
	VisitCallExpr(expr Call) (R, error)
	VisitGetExpr(expr Get) (R, error)
	VisitSetExpr(expr Set) (R, error)
	VisitThisExpr(expr This) (R, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
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

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewLogical(left Expr, operator token.Token, right Expr) *Logical {
	return &Logical{Left: left, Operator: operator, Right: right}
}

func (l Logical) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitLogicalExpr(l)
}

type Unary struct {
	Operator token.Token
	Right    Expr
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
	Left      Expr
	Operator2 token.Token
	Right     Expr
}

func NewTernary(Condition Expr, Operator1 token.Token, Left Expr, Operator2 token.Token, Right Expr) *Ternary {
	return &Ternary{Condition: Condition, Operator1: Operator1, Left: Left, Operator2: Operator2, Right: Right}
}

func (t Ternary) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitTernaryExpr(t)
}

type Variable struct {
	Name token.Token
}

func NewVariable(name token.Token) *Variable {
	return &Variable{Name: name}
}

func (v Variable) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func NewAssign(name token.Token, value Expr) *Assign {
	return &Assign{Name: name, Value: value}
}

func (a Assign) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitAssignExpr(a)
}

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func NewCall(callee Expr, paren token.Token, arguments []Expr) *Call {
	return &Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (c Call) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitCallExpr(c)
}

type Get struct {
	Object Expr
	Name token.Token
}

func NewGet(object Expr, name token.Token) *Get {
	return &Get{Object: object, Name: name}
}

func (g Get) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitGetExpr(g)
}

type Set struct {
	Object Expr
	Name token.Token
	Value Expr
}

func NewSet(object Expr, name token.Token, value Expr) *Set {
	return &Set{Object: object, Name: name, Value: value}
}

func (s Set) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitSetExpr(s)
}

type This struct {
	Keyword token.Token
}

func NewThis(keyword token.Token) *This {
	return &This{Keyword: keyword}
}

func (t This) Accept(visitor Visitor[any]) (any, error) {
	return visitor.VisitThisExpr(t)
}