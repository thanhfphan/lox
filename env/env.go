package env

import (
	"fmt"
	"lox/ast"
)

type Env struct {
	values    map[string]any
	enclosing *Env
}

func New() *Env {
	return &Env{
		values:    map[string]any{},
		enclosing: nil,
	}
}

func (e *Env) SetEnclosing(environemnt *Env) {
	e.enclosing = environemnt
}

func (e *Env) Define(name string, val any) {
	e.values[name] = val
}

func (e *Env) Get(token *ast.Token) (any, error) {
	if val, has := e.values[token.Lexeme()]; has {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(token)
	}

	return nil, fmt.Errorf("Undefined variable '%s'", token.Lexeme())
}

func (e *Env) Assign(token *ast.Token, val any) error {
	if _, has := e.values[token.Lexeme()]; has {
		e.values[token.Lexeme()] = val
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(token, val)
	}

	return fmt.Errorf("Undefined variable '%s'", token.Lexeme())
}
