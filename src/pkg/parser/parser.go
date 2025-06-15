package parser

import (
	"github.com/lidanielm/glox/src/pkg/token"
	"github.com/lidanielm/glox/src/pkg/internal/ast"
	"github.com/lidanielm/glox/src/pkg/lox_error"
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


type ParseError struct {
	message string
}


func (p *Parser) Parse() (ast.Expr, *lox_error.Error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	
	return expr, nil
}


// Evaluate the expression recursively
func (p *Parser) expression() (ast.Expr, *lox_error.Error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	return expr, nil
}


// Evaluate equality operation recursively
func (p *Parser) equality() (ast.Expr, *lox_error.Error) {
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


// Evaluate ternary operation
// func (p *Parser) ternary() (ast.Expr, *lox_error.Error) {
// 	expr, err := p.comparison()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for p.match(token.INTERRO) {
// 		operator := p.comparison()
		
// 	}
// }


// Evaluate comparison operation recursively
func (p *Parser) comparison() (ast.Expr, *lox_error.Error) {
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
func (p *Parser) term() (ast.Expr, *lox_error.Error) {
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
func (p *Parser) factor() (ast.Expr, *lox_error.Error) {
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
func (p *Parser) unary() (ast.Expr, *lox_error.Error) {
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

func (p *Parser) primary() (ast.Expr, *lox_error.Error) {
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
