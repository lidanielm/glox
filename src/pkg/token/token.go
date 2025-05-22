package token

import "fmt"

type Token struct {
	typ TokenType
	lexeme string
	literal interface{}
	line int
}

func NewToken(typ TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{typ: typ, lexeme: lexeme, literal: literal, line: line}
}

func (token *Token) ToString() string {
	return "Type: " + string(byte(token.typ)) + ", Lexeme: " + token.lexeme + ", Literal: " + fmt.Sprintf("%v", token.literal)
}
