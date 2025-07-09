package parser

import (
	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/internal/stmt"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

/** DECLARING TYPE DEFINITIONS AND CONSTRUCTORS **/
type Parser struct {
	tokens []token.Token
	curr int
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
	if p.match(token.PRINT) {
		return p.printStatement()
	} else if p.match(token.LEFT_BRACE) {
		return stmt.NewBlock(), nil
	}

	return p.expressionStatement()
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


// Evaluate the expression recursively
func (p *Parser) expression() (ast.Expr, error) {
	expr, err := p.assignment()
	if err != nil {
		return nil, err
	}
	return expr, nil
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

		if variable, ok := expr.(ast.Variable); ok {
			name := variable.Name
			return ast.NewAssign(name, value), nil
		}

		return nil, lox_error.NewRuntimeError(equals, "Invalid assignment target.")
	}

	return expr, nil
}


// Evaluate ternary operation
func (p *Parser) ternary() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.INTERRO) {
		operator1 := p.previous()
		left, err := p.equality()
		if err != nil {
			return nil, err
		}

		operator2 := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = ast.NewTernary(expr, operator1, left, operator2, right)
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
	
	// If there isn't a unary operator, parse it as a primary operation
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	return expr, nil
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

	if p.match(token.IDENTIFIER) {
		return ast.NewVariable(p.previous()), nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(token.RIGHT_PAREN) {
			return nil, lox_error.NewError(p.peek(), "Expect ')' after expression.")
		}
		return ast.NewGrouping(expr), nil
	}

	return nil, lox_error.NewError(p.peek(), "Expecting expression.")
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
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
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
