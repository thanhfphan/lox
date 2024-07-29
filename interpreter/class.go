package interpreter

var _ Callable = (*Class)(nil)

type Class struct {
	name string
}

func NewClass(name string) *Class {
	return &Class{
		name: name,
	}
}

func (c *Class) Arity() int {
	return 0
}

func (c *Class) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewInstance(c)
	return instance
}

func (c *Class) String() string {
	return c.name
}
