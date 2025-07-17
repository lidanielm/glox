package interpreter

import (
	"fmt"
	"time"

	"github.com/lidanielm/glox/src/pkg/internal/stmt"
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Callable interface {
	Arity() int
	Call(ip *Interpreter, arguments []any) (any, error)
}

type FunctionType int

const (
	NONE_FUNC FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

type Function struct {
	declaration stmt.Function
	closure *Env
	isInitializer bool
}

func NewFunction(declaration stmt.Function, closure *Env, isInitializer bool) *Function {
	return &Function{declaration: declaration, closure: closure, isInitializer: isInitializer}
}

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

func (f *Function) Call(ip *Interpreter, arguments []any) (any, error) {
	env := NewEnv().WithParent(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	err := ip.executeBlock(f.declaration.Body, env)
	if err != nil {
		if returnError, ok := err.(lox_error.ReturnError); ok {
			if f.isInitializer {
				return f.closure.GetAt(0, *token.NewToken(token.THIS, "this", nil, 0))
			}
			return returnError.Value, nil
		}
		return nil, err
	}

	return nil, nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

func (f *Function) Bind(instance *Instance) *Function {
	env := NewEnv().WithParent(f.closure)
	env.Define("this", instance)
	return NewFunction(f.declaration, env, f.isInitializer)
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