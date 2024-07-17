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

func (e *Env) Get(token *ast.Token) (any, error) {
	if val, has := e.values[token.Lexeme()]; has {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(token)
	}

	return nil, fmt.Errorf("Get: Undefined variable '%s'", token.Lexeme())
}

func (e *Env) Assign(token *ast.Token, val any) error {
	if _, has := e.values[token.Lexeme()]; has {
		e.values[token.Lexeme()] = val
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(token, val)
	}

	return fmt.Errorf("Assign: Undefined variable '%s'", token.Lexeme())
}
