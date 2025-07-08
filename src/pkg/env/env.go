package env

import (
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Env struct {
	values map[string]any
}

func NewEnv() *Env {
	values := make(map[string]any)
	return &Env{values: values}
}

func (e Env) Define(name string, value any) {
	e.values[name] = value
}

func (e Env) Get(name token.Token) (any, error) {
	value, exists := e.values[name.Lexeme]
	if exists {
		return value, nil
	} else {
		return nil, lox_error.NewRuntimeError(name, "Undefined variable '" + name.Lexeme + "'.")
	}
}

func (e Env) Assign(name token.Token, value any) error {
	_, exists := e.values[name.Lexeme]
	if exists {
		e.values[name.Lexeme] = value
		return nil
	}

	return lox_error.NewRuntimeError(name, "Undefined variable '" + name.Lexeme + "'.")
}