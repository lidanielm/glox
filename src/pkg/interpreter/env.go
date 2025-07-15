package interpreter

import (
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Env struct {
	parent *Env           // enclosing environment
	values map[string]any // variable-value map
}

func NewEnv() *Env {
	values := make(map[string]any)
	return &Env{values: values}
}

func (e *Env) WithParent(parent *Env) *Env {
	e.parent = parent
	return e
}

func (e *Env) Define(name string, value any) {
	e.values[name] = value
}

func (e *Env) Get(name token.Token) (any, error) {
	value, exists := e.values[name.Lexeme]
	if exists {
		return value, nil
	}

	// If no parent environment, then this variable doesn't exist
	if e.parent == nil {
		return nil, lox_error.NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	}

	return e.parent.Get(name)
}

func (e *Env) GetAt(distance int, name token.Token) (any, error) {
	value, exists := e.ancestor(distance).values[name.Lexeme]
	if exists {
		return value, nil
	} else {
		return nil, lox_error.NewError(name, "Undefined variable '"+name.Lexeme+"'.")
	}
}

func (e *Env) Assign(name token.Token, value any) error {
	_, exists := e.values[name.Lexeme]
	if exists {
		e.values[name.Lexeme] = value
		return nil
	}

	// If no parent environment, then this variable doesn't exist
	if e.parent == nil {
		return lox_error.NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	}


	return e.parent.Assign(name, value)
}

func (e *Env) AssignAt(distance int, name token.Token, value any) error {
	e.ancestor(distance).values[name.Lexeme] = value
	return nil
}

func (e *Env) ancestor(distance int) *Env {
	env := e
	for i := 0; i < distance; i++ {
		env = env.parent
	}

	return env
}