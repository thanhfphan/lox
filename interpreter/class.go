package interpreter

var _ Callable = (*Class)(nil)

type Class struct {
	superClass *Class
	name       string
	methods    map[string]*Function
}

func NewClass(name string, methods map[string]*Function, superClass *Class) *Class {
	return &Class{
		name:       name,
		methods:    methods,
		superClass: superClass,
	}
}

func (c *Class) Arity() int {
	initalizer := c.FindMethod("init")
	if initalizer == nil {
		return 0
	}

	return initalizer.Arity()
}

func (c *Class) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewInstance(c)
	initMethod := c.FindMethod("init")
	if initMethod != nil {
		initMethod.Bind(instance).Call(interpreter, arguments)
	}
	return instance
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) FindMethod(name string) *Function {
	if val, ok := c.methods[name]; ok {
		return val
	}

	if c.superClass != nil {
		return c.superClass.FindMethod(name)
	}

	return nil
}
