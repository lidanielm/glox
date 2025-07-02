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
		err.Error()
	} else {
		err.Message = " at '" + tok.Lexeme + "': " + msg
		err.Error()
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
	err.Error()
	return err
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Error at [line %d]: %s", e.Token.Line, e.Message)
}