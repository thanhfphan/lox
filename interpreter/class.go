package interpreter

var _ Callable = (*Class)(nil)

type Class struct {
	name    string
	methods map[string]*Function
}

func NewClass(name string, methods map[string]*Function) *Class {
	return &Class{
		name:    name,
		methods: methods,
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

func (c *Class) FindMethod(name string) *Function {
	val, ok := c.methods[name]
	if !ok {
		return nil
	}

	return val
}
