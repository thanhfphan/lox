package interpreter

type Class struct {
	name string
}

func NewClass(name string) *Class {
	return &Class{
		name: name,
	}
}

func (c *Class) String() string {
	return c.name
}
