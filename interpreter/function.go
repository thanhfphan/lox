package interpreter

import (
	"lox/ast"
	"lox/env"
)

var _ Callable = (*Function)(nil)

type Function struct {
	declaration   *ast.FunctionStmt
	closure       *env.Env
	isInitializer bool
}

func NewFunction(declaration *ast.FunctionStmt, closure *env.Env, isInitializer bool) *Function {
	return &Function{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) any {
	env := env.New(f.closure)
	for i := 0; i < len(f.declaration.Params); i++ {
		env.Define(f.declaration.Params[i].Lexeme(), arguments[i])
	}

	var returnValue any

	// hack for back to top of the Stack
	func() {
		defer func() {
			if r := recover(); r != nil {
				tmp, ok := r.(*Return)
				if !ok {
					return
				}

				if f.isInitializer {
					returnValue = f.closure.GetAt(0, "this")
					return
				}

				// return statement
				returnValue = tmp.Value
			}
		}()

		interpreter.executeBlock(f.declaration.Body, env)
	}()

	if f.isInitializer {
		return f.closure.GetAt(0, "this")
	}

	return returnValue
}

func (f *Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme() + ">"
}

func (f *Function) Bind(ins *Instance) *Function {
	env := env.New(f.closure)
	env.Define("this", ins)
	return NewFunction(f.declaration, env, f.isInitializer)
}
