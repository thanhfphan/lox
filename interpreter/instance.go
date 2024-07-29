package interpreter

import (
	"fmt"
	"lox/token"
)

var _ Callable = (*Instance)(nil)

type Instance struct {
	class  *Class
	fields map[string]any
}

func NewInstance(c *Class) *Instance {
	return &Instance{
		class:  c,
		fields: make(map[string]any),
	}
}

func (i *Instance) Arity() int {
	return 0
}

func (i *Instance) Call(interpreter *Interpreter, arguments []any) any {
	panic("unimplemented")
}

func (i *Instance) String() string {
	return i.class.name + " instance"
}

func (i *Instance) Set(name *token.Token, value any) {
	i.fields[name.Lexeme()] = value
}

func (i *Instance) Get(name *token.Token) any {
	val, has := i.fields[name.Lexeme()]
	if has {
		return val
	}

	method := i.class.FindMethod(name.Lexeme())
	if method != nil {
		return method.Bind(i)
	}

	panic(fmt.Errorf("%s Undefined property '%s'.", name, name.Lexeme()))
}
