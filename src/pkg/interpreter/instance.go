package interpreter

import (
	"github.com/lidanielm/glox/src/pkg/lox_error"
	"github.com/lidanielm/glox/src/pkg/token"
)

type Instance struct {
	class *Class
	fields map[string]any
}

func NewInstance(class *Class) *Instance {
	fields := make(map[string]any)
	return &Instance{class: class, fields: fields}
}

func (i *Instance) String() string {
	return i.class.name + " instance"
}

func (i *Instance) Get(name string) (any, error) {
	value, exists := i.fields[name]
	if !exists {
		return nil, lox_error.NewRuntimeError(token.Token{}, "Undefined property '"+name+"'.")
	}
	
	return value, nil
}

func (i *Instance) Set(name string, property any) {
	i.fields[name] = property
}

func (i *Instance) FindMethod(name string) (*Function, error) {
	return i.class.FindMethod(name)
}