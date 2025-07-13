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
	VisitIfStmt(stmt If) error
	VisitWhileStmt(stmt While) error
	VisitBreakStmt(stmt Break) error
	VisitContinueStmt(stmt Continue) error
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

func NewBlock(statements []Stmt) *Block {
	return &Block{Statements: statements}
}

func (b Block) Accept(visitor Visitor[any]) error {
	return visitor.VisitBlockStmt(b)
}

type If struct {
	Condition ast.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition ast.Expr, thenBranch Stmt, elseBranch Stmt) *If {
	return &If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (i If) Accept(visitor Visitor[any]) error {
	return visitor.VisitIfStmt(i)
}

type While struct {
	Condition ast.Expr
	Body Stmt
	Increment ast.Expr // optional, for for-loops
}

func NewWhile(condition ast.Expr) *While {
	return &While{Condition: condition}
}

func (w *While) WithBody(body Stmt) *While {
	w.Body = body
	return w
}

func (w *While) WithIncrement(increment ast.Expr) *While {
	w.Increment = increment
	return w
}

func (w While) Accept(visitor Visitor[any]) error {
	return visitor.VisitWhileStmt(w)
}

type Break struct {
	Loop *While
}

func NewBreak(loop *While) *Break {
	return &Break{Loop: loop}
}

func (b Break) Accept(visitor Visitor[any]) error {
	return visitor.VisitBreakStmt(b)
}

type Continue struct {
	Loop *While
}

func NewContinue(loop *While) *Continue {
	return &Continue{Loop: loop}
}

func (c Continue) Accept(visitor Visitor[any]) error {
	return visitor.VisitContinueStmt(c)
}