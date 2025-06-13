package lox_error

import (
	"fmt"
	"os"

	"github.com/lidanielm/glox/src/pkg/token"
)

type Error struct {
	token token.Token
	message string
}

func NewError(tok token.Token, msg string) *Error {
	err := &Error{token: tok, message: msg}
	if tok.Type == token.EOF {
		Report(tok.Line, " at end: " + msg)
	} else {
		Report(tok.Line, " at '" + tok.Lexeme + "': " + msg)
	}
	return err
}

func Report(line int, message string) {
	fmt.Fprintf(os.Stderr, "Error at [line %d]: %s", line, message)
}

// func (err *Error) throwError() {
// 	// Throw error
// 	reportError(line, "", message)
// }

// func (iptr *Interpreter) reportError(int line, string where, string message) {
// 	fmt.Println("Error at [line " + line + "]" + where + ": " + message)
// 	iptr.hadError = true
// }
