package lox_error

import (
	"fmt"

	"github.com/lidanielm/glox/src/pkg/token"
)

type LoxError struct {
	Token token.Token
	Message string
}

func NewError(tok token.Token, msg string) *LoxError {
	err := &LoxError{Token: tok}
	if tok.Type == token.EOF {
		err.Message = " at end: " + msg
		fmt.Println(err.Error())
	} else {
		err.Message = " at '" + tok.Lexeme + "': " + msg
		fmt.Println(err.Error())
	}
	return err
}

func (e *LoxError) Error() string {
	return fmt.Sprintf("Error at [line %d]: %s", e.Token.Line, e.Message)
}

type RuntimeError struct {
	Token token.Token
	Message string
}

func NewRuntimeError(tok token.Token, msg string) *RuntimeError {
	err := &RuntimeError{Token: tok, Message: msg}
	// fmt.Println(err.Error())
	return err
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError at [line %d]: %s", e.Token.Line, e.Message)
}

type ParseError struct {
	Token token.Token
	Message string
}

func NewParseError(tok token.Token, message string) *ParseError {
	err := &ParseError{Token: tok, Message: message}
	// fmt.Println(err.Error())
	return err
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ParseError at [line %d]: %s", e.Token.Line, e.Message)
}