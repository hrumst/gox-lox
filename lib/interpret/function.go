package interpret

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
)

type LoxFunction struct {
	interpreter   *Interpreter
	declaration   *parse.StmtFunction
	closure       *Environment
	isInitializer bool
}

func NewLoxFunction(
	interpreter *Interpreter,
	declaration *parse.StmtFunction,
	closure *Environment,
	isInitializer bool,
) *LoxFunction {
	return &LoxFunction{
		interpreter:   interpreter,
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (l *LoxFunction) String() string {
	return fmt.Sprintf("[function] %s", l.declaration.Name.Lexeme)
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxFunction) Call(args []*scan.LoxValue) (*scan.LoxValue, error) {
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

	if l.isInitializer {
		//todo make refactor token -> string
		this := scan.NewToken(scan.THIS, "this", nil, l.declaration.Name.Line)
		return l.closure.getAt(0, this)
	}
	return scan.NewNilLoxValue(), nil
}

func (l *LoxFunction) bind(instance *LoxClassInstance) *LoxFunction {
	environment := NewEnvironment(l.closure)
	environment.Define("this", scan.NewClassInstanceLoxValue(instance))
	return NewLoxFunction(l.interpreter, l.declaration, environment, l.isInitializer)
}
