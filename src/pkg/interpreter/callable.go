package interpreter

import (
	"fmt"
	"time"

	"github.com/lidanielm/glox/src/pkg/env"
	"github.com/lidanielm/glox/src/pkg/internal/stmt"
	"github.com/lidanielm/glox/src/pkg/lox_error"
)

type Callable interface {
	Arity() int
	Call(ip *Interpreter, arguments []any) (any, error)
}

type Function struct {
	declaration stmt.Function
	closure *env.Env
}

func NewFunction(declaration stmt.Function, closure *env.Env) *Function {
	return &Function{declaration: declaration, closure: closure}
}

func (f Function) Arity() int {
	return len(f.declaration.Params)
}

func (f Function) Call(ip *Interpreter, arguments []any) (any, error) {
	env := env.NewEnv().WithParent(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	err := ip.executeBlock(f.declaration.Body, env)
	if err != nil {
		if returnError, ok := err.(lox_error.ReturnError); ok {
			return returnError.Value, nil
		}
		return nil, err
	}

	return nil, nil
}

func (f Function) toString() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

type ClockFn struct{}

func (c *ClockFn) Arity() int {
	return 0
}

func (c *ClockFn) Call(ip *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}

func (c *ClockFn) String() string {
	return "<native fn>"
}