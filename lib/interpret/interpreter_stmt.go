package interpret

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
)

func (i *Interpreter) VisitStmtVar(stmt *parse.StmtVar) (interface{}, error) {
	value := scan.NewNilLoxValue()
	if stmt.Initializer != nil {
		initValue, err := i.Evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
		value = initValue
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitStmtBlock(stmt *parse.StmtBlock) (interface{}, error) {
	res, err := i.executeBlock(stmt.Stmts, NewEnvironment(i.environment))
	return res, err
}

func (i *Interpreter) VisitStmtIf(stmt *parse.StmtIf) (interface{}, error) {
	conditionValue, err := i.Evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	if conditionValue.Bool() {
		if res, err := i.execute(stmt.ThenBranch); err != nil {
			return nil, err
		} else if control, ok := res.(executeControl); ok {
			return control, nil
		}
	} else if stmt.ElseBranch != nil {
		if _, err := i.execute(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitStmtWhile(stmt *parse.StmtWhile) (interface{}, error) {
	for {
		conditionValue, err := i.Evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !conditionValue.Bool() {
			break
		}

		if res, err := i.execute(stmt.Body); err != nil {
			return nil, err
		} else if control, ok := res.(executeControl); ok {
			if control.isBreak() {
				break
			}
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitStmtExecuteControl(stmt *parse.StmtExecuteControl) (interface{}, error) {
	switch stmt.Control.Type {
	case scan.BREAK:
		return newBreakExecuteControl(), nil
	default:
		return newContinueExecuteControl(), nil
	}
}

func (i *Interpreter) VisitStmtFunction(stmt *parse.StmtFunction) (interface{}, error) {
	function := scan.NewCallableLoxValue(NewLoxFunction(i, stmt, i.environment, false))
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitStmtReturn(stmt *parse.StmtReturn) (interface{}, error) {
	var value *scan.LoxValue
	if stmt.Value != nil {
		var err error
		value, err = i.Evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return newReturnExecuteControl(value), nil
}

func (i *Interpreter) VisitStmtPrint(stmt *parse.StmtPrint) (interface{}, error) {
	value, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintln(i.writer, value.String()); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitStmtClass(stmt *parse.StmtClass) (interface{}, error) {
	i.environment.Define(stmt.Name.Lexeme, nil)
	environment := i.environment

	var superClass *LoxClass
	if stmt.SuperClass != nil {
		superClassLoxValue, err := i.Evaluate(stmt.SuperClass)
		if err != nil {
			return nil, err
		}
		superClassLox, err := superClassLoxValue.Callable()
		if err != nil {
			return nil, NewRuntimeError(
				fmt.Sprintf("superclass must be a class. error: %s", err.Error()),
				&stmt.SuperClass.Name,
			)
		}

		environment = NewEnvironment(environment)
		environment.Define("super", superClassLoxValue)

		superClass = superClassLox.(*LoxClass)
	}

	methods := make(map[string]scan.LoxCallable)
	for _, method := range stmt.Methods {
		stmtFunc := method.(*parse.StmtFunction)
		loxFunc := NewLoxFunction(i, stmtFunc, environment, stmtFunc.Name.Lexeme == "init")
		methods[stmtFunc.Name.Lexeme] = loxFunc
	}
	class := NewLoxClass(i, stmt, superClass, methods)
	if err := i.environment.Assign(stmt.Name, scan.NewClassLoxValue(scan.LoxClass(class))); err != nil {
		return nil, err
	}
	return nil, nil
}
