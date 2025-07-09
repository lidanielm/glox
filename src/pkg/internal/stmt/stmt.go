package stmt

import (
	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Stmt interface {
	Accept(visitor Visitor[any]) error
}

type Visitor[R any] interface {
	VisitExpressionStmt(stmt Expression) error
	VisitPrintStmt(stmt Print) error
	VisitVarStmt(stmt Var) error
	VisitBlockStmt(stmt Block) error
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

type Var struct {
	Name token.Token
	Initializer ast.Expr
}

func NewVar(name token.Token, initializer ast.Expr) *Var {
	return &Var{Name: name, Initializer: initializer}
}

func (v Var) Accept(visitor Visitor[any]) error {
	return visitor.VisitVarStmt(v)
}

type Block struct {
	Statements []Stmt
}

func NewBlock() *Block {
	return &Block{Statements: []Stmt{}}
}

func (b Block) Accept(visitor Visitor[any]) error {
	return visitor.VisitBlockStmt(b)
}