package interpreter

import (
	"fmt"

	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/internal/stmt"
	"github.com/lidanielm/glox/src/pkg/internal/tool"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Resolver struct {
	ip *Interpreter
	scopes tool.Stack[map[string]bool]
	currFunc FunctionType
}

func NewResolver(ip *Interpreter) *Resolver {
	scopes := tool.NewStack[map[string]bool]()
	return &Resolver{ip: ip, scopes: *scopes, currFunc: NONE}
}

func (r *Resolver) VisitBlockStmt(stmt stmt.Block) error {
	r.beginScope()
	r.ResolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitBreakStmt(stmt stmt.Break) error {
	return nil
}

func (r *Resolver) VisitContinueStmt(stmt stmt.Continue) error {
	return nil
}

func (r *Resolver) VisitVarStmt(stmt stmt.Var) error {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}

	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt stmt.Function) error {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt stmt.Expression) error {
	_, err := r.resolveExpr(stmt.Expr)
	return err
}

func (r *Resolver) VisitIfStmt(stmt stmt.If) error {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	
	err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return err
	}

	if stmt.ElseBranch != nil {
		err = r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitPrintStmt(stmt stmt.Print) error {
	_, err := r.resolveExpr(stmt.Expr)
	return err
}

func (r *Resolver) VisitWhileStmt(stmt stmt.While) error {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}

	err = r.resolveStmt(stmt.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitReturnStmt(stmt stmt.Return) error {
	if r.currFunc == NONE {
		return lox_error.NewParseError(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value == nil {
		return nil
	}

	_, err := r.resolveExpr(stmt.Value)
	return err
}

func (r *Resolver) VisitVariableExpr(expr ast.Variable) (any, error) {
	// If variable exists in current scope but value is false,
	// that means we have declared it but not yet defined it
	if !r.scopes.IsEmpty() {
		val, ok := r.scopes.Peek()[expr.Name.Lexeme]
		if ok && !val {
			return nil, lox_error.NewParseError(expr.Name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr ast.Assign) (any, error) {
	_, err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}
	
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr ast.Binary) (any, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr ast.Call) (any, error) {
	_, err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		_, err = r.resolveExpr(argument)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr ast.Grouping) (any, error) {
	_, err := r.resolveExpr(expr.Expression)
	return nil, err
}

func (r *Resolver) VisitLiteralExpr(expr ast.Literal) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr ast.Logical) (any, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitTernaryExpr(expr ast.Ternary) (any, error) {
	_, err := r.resolveExpr(expr.Condition)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr ast.Unary) (any, error) {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) ResolveStmts(stmts []stmt.Stmt) (any, error) {
	for _, stmt := range stmts {
		err := r.resolveStmt(stmt)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) resolveStmt(stmt stmt.Stmt) error {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) (any, error) {
	return expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := r.scopes.Length(); i >= 0; i-- {
		if _, ok := r.scopes.Get(i)[name.Lexeme]; ok {
			r.ip.resolve(expr, r.scopes.Length() - 1 - i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function stmt.Function, ftype FunctionType) {
	enclosingFunc := r.currFunc
	r.currFunc = ftype

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmts(function.Body)
	r.endScope()

	r.currFunc = enclosingFunc
}

func (r *Resolver) beginScope() {
	newScope := make(map[string]bool, 0)
	r.scopes.Push(newScope)
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek()
	_, ok := scope[name.Lexeme]
	if ok {
		fmt.Println(lox_error.NewError(name, "Already a variable with this name in this scope.").Error())
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek()
	scope[name.Lexeme] = true
}