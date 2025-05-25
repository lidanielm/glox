package parser

import (
	"fmt"
	"github.com/lidanielm/glox/src/pkg/token"
	"github.com/lidanielm/glox/src/pkg/internal/ast"
)


type Parser {
	tokens []token.TokenType
	curr int
}


// Constructor for Parser
func NewParser(tokens []token.TokenType) *Parser {
	return &Parser{tokens: tokens, curr: 0}
}


// Evaluate the expression recursively
func (p Parser) expression() ast.Expr {
	return p.equality()
}


// Evaluate equality operation recursively
func (p Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr := ast.NewBinary(expr, operator, right)
	}

	return expr
}


// Evaluate comparison operation recursively
func (p Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr := ast.NewBinary(expr, operator, right)
	}

	return expr
}

// Evaluate an addition/subtraction operation recursively
func (p Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right := p.factor()
		expr := ast.NewBinary(expr, operator, right)
	}

	return expr
}

// Evaluate a multiplication/division operation recursively
func (p Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right := p.unary()
		expr := ast.NewBinary(expr, operator, right)
	}

	return expr
}

// Evaluate a unary operation recursively
func (p Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}
	
	// If there isn't a unary operator, parse it as a primary operation
	return primary()
}

func (p Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return ast.NewLiteral(false)
	}

	if p.match(token.TRUE) {
		return ast.NewLiteral(true)
	}

	if p.match(token.NIL) {
		return ast.NewLiteral(nil)
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteral(p.previous().literal)
	}

	if p.match(token.LEFT_PAREN) {
		expr := expression()
		consume(RIGHT_PAREN, "Expect ')' after expression.")
		return ast.NewGrouping(expr)
	}

	// TODO: should never get here, so throw some syntax error if it does
	return nil
}




/* HELPERS */

// Check if the current token has any of the given type. If so, it consumes the token
// and returns true. Otherwise, it returns false and leaves the token alone.
func (p Parser) match(tokenTypes ...token.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

// Check if the current token is of the given type
func (p Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	
	return p.peek().typ == tokenType
}

// Advance to the next token in the sequence
func (p Parser) advance() token.Token {
	if !p.isAtEnd() {
		current++
	}

	return p.previous()
}

// Check if all tokens are parsed
func (p Parser) isAtEnd() bool {
	return peek().typ == token.EOF
}

// Return the current token that hasn't yet been consumed
func (p Parser) peek() token.Token {
	return tokens[current]
}

// Return the most recently consumed token
func (p Parser) previous() token.Token {
	return tokens[current - 1]
}
