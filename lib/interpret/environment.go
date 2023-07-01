package interpret

import (
	"github.com/hrumst/gox-lox/lib/scan"
)

type Environment struct {
	values    map[string]*scan.LoxValue
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]*scan.LoxValue),
	}
}

func (e *Environment) Define(name string, value *scan.LoxValue) {
	e.values[name] = value
}

// todo add support for not yet initialised variable
func (e *Environment) Get(token scan.Token) (*scan.LoxValue, error) {
	if _, exists := e.values[token.Lexeme]; exists {
		return e.values[token.Lexeme], nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(token)
	}
	return nil, NewRuntimeError("undefined variable", &token)
}

func (e *Environment) Assign(token scan.Token, value *scan.LoxValue) error {
	if _, exists := e.values[token.Lexeme]; exists {
		e.values[token.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(token, value)
	}
	return NewRuntimeError("undefined variable", &token)
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i += 1 {
		env = env.enclosing
	}
	return env
}

func (e *Environment) getAt(distance int, token scan.Token) (*scan.LoxValue, error) {
	return e.ancestor(distance).Get(token)
}

func (e *Environment) assignAt(distance int, name scan.Token, value *scan.LoxValue) {
	e.ancestor(distance).values[name.Lexeme] = value
}
