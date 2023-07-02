package interpret

import (
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
	"io"
)

type Interpreter struct {
	writer      io.Writer
	environment *Environment
	globals     *Environment
	locals      map[parse.Expression]int
}

func NewInterpreter(writer io.Writer) *Interpreter {
	globalFuncs := NewEnvironment(nil)
	globalFuncs.Define("clock", scan.NewCallableLoxValue(NewClockFunction()))

	return &Interpreter{
		writer:      writer,
		environment: NewEnvironment(nil),
		globals:     globalFuncs,
		locals:      make(map[parse.Expression]int),
	}
}

func (i *Interpreter) Interpret(stmts []parse.Statement) error {
	for _, stmt := range stmts {
		if _, err := i.execute(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Evaluate(expr parse.Expression) (*scan.LoxValue, error) {
	result, err := expr.Accept(i)
	if err != nil {
		return nil, err
	}
	switch rt := result.(type) {
	case *scan.Literal:
		return rt.Value, nil
	case *scan.LoxValue:
		return rt, nil
	case nil:
		return scan.NewNilLoxValue(), nil
	default:
		return nil, NewRuntimeError("unsupported expression evaluate result type", nil)
	}
}

func (i *Interpreter) execute(stmt parse.Statement) (interface{}, error) {
	res, err := stmt.Accept(i)
	return res, err
}

func (i *Interpreter) executeBlock(stmts []parse.Statement, nextEnv *Environment) (interface{}, error) {
	prevEnv := i.environment
	i.environment = nextEnv
	defer func() {
		i.environment = prevEnv
	}()

	for _, stmt := range stmts {
		if res, err := i.execute(stmt); err != nil {
			return nil, err
		} else if control, ok := res.(executeControl); ok {
			return control, nil
		}
	}
	return nil, nil
}

func (i *Interpreter) lookUpVariable(token scan.Token, expr parse.Expression) (*scan.LoxValue, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.getAt(distance, token)
	}
	return i.globals.Get(token)
}

func (i *Interpreter) resolve(expr parse.Expression, depth int) {
	i.locals[expr] = depth
}
