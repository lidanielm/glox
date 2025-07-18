package parser

import (
	"slices"

	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/internal/stmt"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

/** DECLARING TYPE DEFINITIONS AND CONSTRUCTORS **/
type Parser struct {
	tokens []token.Token
	curr int
	enclosingLoop *stmt.While
}

// Constructor for Parser
func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, curr: 0}
}


func (p *Parser) Parse() ([]stmt.Stmt, error) {
	statements := []stmt.Stmt{}
	for !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, declaration)
	}
	
	return statements, nil
}


func (p *Parser) declaration() (stmt.Stmt, error) {
	if p.match(token.CLASS) {
		class, err := p.classDeclaration()
		if err != nil {
			return nil, err
		}

		return class, nil
	}
	if p.match(token.FUN) {
		fn, err := p.function("function")
		if err != nil {
			return nil, err
		}

		return fn, nil
	}
	if p.match(token.VAR) {
		stmt, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}

		return stmt, nil
	}

	return p.statement()
}

func (p *Parser) classDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_BRACE, "Expect '{' after class name.")
	if err != nil {
		return nil, err
	}

	methods := make([]stmt.Function, 0)
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		fn, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, fn)
	}

	_, err = p.consume(token.RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return stmt.NewClass(name, methods), nil
}

func (p *Parser) function(kind string) (stmt.Function, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return stmt.Function{}, err
	}

	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return stmt.Function{}, err
	}

	params := []token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(params) >= 255 {
				return stmt.Function{}, lox_error.NewParseError(p.peek(), "Can't have more than 255 parameters.")
			}
	
			identifier, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return stmt.Function{}, err
			}
	
			params = append(params, identifier)

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return stmt.Function{}, err
	}

	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return stmt.Function{}, err
	}

	body, err := p.block()
	if err != nil {
		return stmt.Function{}, err
	}

	return *stmt.NewFunction(name, params, body), nil
}

func (p *Parser) varDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return stmt.NewVar(name, initializer), nil
}


func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(token.IF) {
		return p.ifStatement()
	} else if p.match(token.PRINT) {
		return p.printStatement()
	} else if p.match(token.WHILE) {
		return p.whileStatement()
	} else if p.match(token.FOR) {
		return p.forStatement()
	} else if p.match(token.BREAK) {
		return p.breakStatement()
	} else if p.match(token.CONTINUE) {
		return p.continueStatement()
	} else if p.match(token.RETURN) {
		return p.returnStatement()
	} else if p.match(token.LEFT_BRACE) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return stmt.NewBlock(block), nil
	}

	return p.expressionStatement()
}

func (p *Parser) ifStatement() (stmt.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after 'if'.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(token.ELSE) {
		elseBranch, err := p.statement()
		if err != nil {
			return nil, err
		}

		return stmt.NewIf(condition, thenBranch, elseBranch), nil
	} else {
		return stmt.NewIf(condition, thenBranch, nil), nil
	}
}

func (p *Parser) whileStatement() (stmt.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	prevLoop := p.enclosingLoop
	whileStmt := stmt.NewWhile(condition)
	p.enclosingLoop = whileStmt

	body, err := p.statement()
	p.enclosingLoop = prevLoop
	if err != nil {
		return nil, err
	}

	whileStmt = whileStmt.WithBody(body)

	return whileStmt, nil
}

func (p *Parser) forStatement() (stmt.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer stmt.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	prevLoop := p.enclosingLoop
	whileStmt := stmt.NewWhile(condition)
	p.enclosingLoop = whileStmt
	
	body, err := p.statement()
	p.enclosingLoop = prevLoop
	if err != nil {
		return nil, err
	}

	if increment != nil {
		whileStmt = whileStmt.WithIncrement(increment)
	}

	if condition == nil {
		condition = ast.NewLiteral(true)
	}

	body = whileStmt.WithBody(body)

	if initializer != nil {
		body = stmt.NewBlock([]stmt.Stmt{
			initializer,
			body,
		})
	}

	return body, nil
}

func (p *Parser) breakStatement() (stmt.Stmt, error) {
	if p.enclosingLoop == nil {
		return nil, lox_error.NewParseError(p.peek(), "'break' statement has no enclosing loop.")
	}

	_, err := p.consume(token.SEMICOLON, "Expect semicolon after 'break'.")
	if err != nil {
		return nil, err
	}

	return stmt.NewBreak(p.enclosingLoop), nil
}

func (p *Parser) continueStatement() (stmt.Stmt, error) {
	if p.enclosingLoop == nil {
		return nil, lox_error.NewParseError(p.peek(), "'continue' statement has no enclosing loop.")
	}

	_, err := p.consume(token.SEMICOLON, "Expect semicolon after 'continue'.")
	if err != nil {
		return nil, err
	}

	return stmt.NewContinue(p.enclosingLoop), nil
}

func (p *Parser) returnStatement() (stmt.Stmt, error) {
	keyword := p.previous()

	var expr ast.Expr
	if !p.check(token.SEMICOLON) {
		var err error
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	
	_, err := p.consume(token.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}
	
	return stmt.NewReturn(keyword, expr), nil
}

func (p *Parser) printStatement() (stmt.Stmt, error) {
	// Evaluate argument
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	// Check if statement is terminated by semicolon
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return stmt.NewPrint(value), nil
}

func (p *Parser) expressionStatement() (stmt.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return stmt.NewExpression(expr), nil
}


func (p *Parser) block() ([]stmt.Stmt, error) {
	statements := []stmt.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, decl)
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	return statements, nil
}


// Evaluate the expression recursively
func (p *Parser) expression() (ast.Expr, error) {
	expr, err := p.assignment()
	if err != nil {
		return nil, err
	}
	return expr, nil
}


func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*ast.Variable); ok {
			name := variable.Name
			return ast.NewAssign(name, value), nil
		}

		return nil, lox_error.NewParseError(equals, "Invalid assignment target.")
	}

	return expr, nil
}


// Evaluate ternary operation
func (p *Parser) ternary() (ast.Expr, error) {
	expr, err := p.logical_or()
	if err != nil {
		return nil, err
	}

	for p.match(token.INTERRO) {
		operator1 := p.previous()
		left, err := p.logical_or()
		if err != nil {
			return nil, err
		}

		operator2 := p.previous()
		right, err := p.logical_or()
		if err != nil {
			return nil, err
		}

		expr = ast.NewTernary(expr, operator1, left, operator2, right)
	}

	return expr, nil
}


func (p *Parser) logical_or() (ast.Expr, error) {
	expr, err := p.logical_and()
	if err != nil {
		return nil, err
	}

	if p.match(token.OR) {
		operator := p.previous()
		right, err := p.logical_and()
		if err != nil {
			return nil, err
		}
		return ast.NewLogical(expr, operator, right), nil
	}

	return expr, nil
}


func (p *Parser) logical_and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		return ast.NewLogical(expr, operator, right), nil
	}

	return expr, nil
}


// Evaluate equality operation recursively
func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr, nil
}


// Evaluate comparison operation recursively
func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr, nil
}

// Evaluate an addition/subtraction operation recursively
func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr, nil
}

// Evaluate a multiplication/division operation recursively
func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr, nil
}

// Evaluate a unary operation recursively
func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnary(operator, right), nil
	}

	return p.call()
}

func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(token.DOT) {
			name, err := p.consume(token.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}

			expr = ast.NewGet(expr, name)
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := []ast.Expr{}
	if !p.check(token.RIGHT_PAREN) {
		expr, err := p.expression()
		for {
			if err != nil {
				return nil, err
			}

			if len(arguments) >= 255 {
				return nil, lox_error.NewParseError(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, expr)

			if !p.match(token.COMMA) {
				break
			}

			expr, err = p.expression()
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return ast.NewCall(callee, paren, arguments), nil
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return ast.NewLiteral(false), nil
	}

	if p.match(token.TRUE) {
		return ast.NewLiteral(true), nil
	}

	if p.match(token.NIL) {
		return ast.NewLiteral(nil), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteral(p.previous().Literal), nil
	} 

	if p.match(token.THIS) {
		return ast.NewThis(p.previous()), nil
	}

	if p.match(token.IDENTIFIER) {
		return ast.NewVariable(p.previous()), nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(token.RIGHT_PAREN) {
			return nil, lox_error.NewParseError(p.peek(), "Expect ')' after expression.")
		}
		return ast.NewGrouping(expr), nil
	}

	return nil, lox_error.NewParseError(p.peek(), "Expecting expression.")
}


// Continues parsing tokens until it reaches a statement boundary
// Used after ParseError is thrown
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
			case token.CLASS:
			case token.FUN:
			case token.VAR:
			case token.FOR:
			case token.IF:
			case token.WHILE:
			case token.PRINT:
			case token.RETURN:
				return
		}

		p.advance()
	}
}


/* HELPERS */

// Check if the current token has any of the given type. If so, it consumes the token
// and returns true. Otherwise, it returns false and leaves the token alone.
func (p *Parser) match(tokenTypes ...token.TokenType) bool {
	if slices.ContainsFunc(tokenTypes, p.check) {
			p.advance()
			return true
		}

	return false
}

// Consume the current token
func (p *Parser) consume(tokenType token.TokenType, message string) (token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return token.Token{}, lox_error.NewParseError(p.peek(), message)
}

// Check if the current token is of the given type
func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	
	return p.peek().Type == tokenType
}

// Advance to the next token in the sequence
func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.curr++
	}

	return p.previous()
}

// Check if all tokens are parsed
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// Return the current token that hasn't yet been consumed
func (p *Parser) peek() token.Token {
	return p.tokens[p.curr]
}

// Return the most recently consumed token
func (p *Parser) previous() token.Token {
	return p.tokens[p.curr - 1]
}
