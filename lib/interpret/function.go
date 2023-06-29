package interpret

import (
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
)

type LoxFunction struct {
	interpreter *Interpreter
	declaration *parse.StmtFunction
	closure     *Environment
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxFunction) Call(args []*scan.LoxValue) (interface{}, error) {
	environment := NewEnvironment(l.closure)
	for i, param := range l.declaration.Params {
		environment.Define(param.Lexeme, args[i])
	}
	result, err := l.interpreter.executeBlock(l.declaration.Body, environment)
	if err != nil {
		return nil, err
	}
	if resControl, ok := result.(executeControl); ok {
		if resControl.value != nil {
			return resControl.value, nil
		}
	}
	return scan.NewNilLoxValue(), nil
}

func NewLoxFunction(
	interpreter *Interpreter,
	declaration *parse.StmtFunction,
	closure *Environment,
) *LoxFunction {
	return &LoxFunction{
		interpreter: interpreter,
		declaration: declaration,
		closure:     closure,
	}
}
