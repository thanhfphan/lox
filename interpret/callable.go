package interpret

type Callable interface {
	Call(interpreter *Interpreter, arguments []any) any
	Arity() int
	ToString() string
}
