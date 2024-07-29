package interpreter

var _ Callable = (*Instance)(nil)

type Instance struct {
	class *Class
}

func NewInstance(c *Class) *Instance {
	return &Instance{
		class: c,
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
