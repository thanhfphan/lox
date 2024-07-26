package env

import (
	"fmt"
	"lox/ast"
)

type Env struct {
	values    map[string]any
	enclosing *Env
}

func New(encolsing *Env) *Env {
	return &Env{
		values:    map[string]any{},
		enclosing: encolsing,
	}
}

func (e *Env) Define(name string, val any) {
	e.values[name] = val
}

func (e *Env) GetAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Env) ancestor(distance int) *Env {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}

func (e *Env) Get(token *ast.Token) (any, error) {
	if val, has := e.values[token.Lexeme()]; has {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(token)
	}

	return nil, fmt.Errorf("Get: Undefined variable '%s'", token.Lexeme())
}

func (e *Env) Assign(name *ast.Token, val any) error {
	if _, has := e.values[name.Lexeme()]; has {
		e.values[name.Lexeme()] = val
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, val)
	}

	return fmt.Errorf("Assign: Undefined variable '%s'", name.Lexeme())
}

func (e *Env) AssignAt(distance int, name *ast.Token, val any) error {
	env := e.ancestor(distance)
	env.values[name.Lexeme()] = val
	return nil
}
