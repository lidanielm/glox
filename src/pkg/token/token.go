package token

import "fmt"

type Token struct {
	Type TokenType
	Lexeme string
	Literal interface{}
	Line int
}

func NewToken(typ TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{Type: typ, Lexeme: lexeme, Literal: literal, Line: line}
}

func (token *Token) ToString() string {
	return "Type: " + string(byte(token.Type)) + ", Lexeme: " + token.Lexeme + ", Literal: " + fmt.Sprintf("%v", token.Literal)
}
