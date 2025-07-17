package interpreter

import (
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type ClassType int

const (
	NONE_CLASS ClassType = iota
	CLASS
)

type IClass interface {
	FindMethod(name string) (*Function, error)
}

type Class struct {
	name string
	methods map[string]*Function
}

func NewClass(name string, methods map[string]*Function) *Class {
	return &Class{name: name, methods: methods}
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) Arity() int {
	initializer, err := c.FindMethod(c.name)
	if err != nil {
		return 0;
	}
	
	return initializer.Arity()
}

func (c *Class) Call(ip *Interpreter, arguments []any) (any, error) {
	instance := NewInstance(c)
	initializer, err := instance.FindMethod("init")
	if err == nil {
		initializer.Bind(instance).Call(ip, arguments)
	}
	return instance, nil
}

func (c *Class) AddMethod(name string, method *Function) {
	c.methods[name] = method
}

func (c *Class) FindMethod(name string) (*Function, error) {
	fn, exists := c.methods[name]
	if !exists {
		return nil, lox_error.NewRuntimeError(token.Token{}, "Undefined method '"+name+"' for class '"+c.String()+"'.")
	}

	return fn, nil
}