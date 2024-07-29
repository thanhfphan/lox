package interpreter

type Callable interface {
	Call(interpreter *Interpreter, arguments []any) any
	Arity() int
	String() string
}
