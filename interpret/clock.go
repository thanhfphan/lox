package interpret

import "time"

var (
	_ Callable = (*Clock)(nil)
)

type Clock struct {
}

func NewClock() *Clock {
	return &Clock{}
}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Call(interpreter *Interpreter, arguments []any) any {
	seconds := time.Now().Second()
	return seconds
}

func (c *Clock) ToString() string {
	return "<native fn>"
}
