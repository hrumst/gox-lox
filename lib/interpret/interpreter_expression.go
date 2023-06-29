package interpret

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
)

func (i *Interpreter) VisitLiteralExpr(expr *parse.LiteralExpression) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *parse.BinaryExpression) (interface{}, error) {
	leftVal, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, ConvertToRuntimeError("evaluate expression error", err, &expr.Operator)
	}
	rightVal, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, ConvertToRuntimeError("evaluate expression error", err, &expr.Operator)
	}

	if leftVal.IsBoolean() && rightVal.IsBoolean() {
		leftBool, rightBool := leftVal.Bool(), rightVal.Bool()
		switch expr.Operator.Type {
		case scan.BANG_EQUAL:
			return scan.NewBooleanLoxValue(leftBool != rightBool), nil
		case scan.EQUAL_EQUAL:
			return scan.NewBooleanLoxValue(leftBool == rightBool), nil
		}
	}

	if leftVal.IsString() && rightVal.IsString() {
		leftStr, rightStr := leftVal.String(), rightVal.String()
		switch expr.Operator.Type {
		case scan.BANG_EQUAL:
			return scan.NewBooleanLoxValue(leftStr != rightStr), nil
		case scan.EQUAL_EQUAL:
			return scan.NewBooleanLoxValue(leftStr == rightStr), nil
		}
	}

	if leftVal.IsString() || rightVal.IsString() {
		leftStr, rightStr := leftVal.String(), rightVal.String()
		switch expr.Operator.Type {
		case scan.PLUS:
			return scan.NewStringLoxValue(leftStr + rightStr), nil
		}
	}

	leftNum, err := leftVal.Number()
	if err != nil {
		return nil, ConvertToRuntimeError("evaluate expression error", err, &expr.Operator)
	}
	rightNum, err := rightVal.Number()
	if err != nil {
		return nil, ConvertToRuntimeError("evaluate expression error", err, &expr.Operator)
	}

	switch expr.Operator.Type {
	case scan.GREATER:
		return scan.NewBooleanLoxValue(leftNum > rightNum), nil
	case scan.GREATER_EQUAL:
		return scan.NewBooleanLoxValue(leftNum >= rightNum), nil
	case scan.LESS:
		return scan.NewBooleanLoxValue(leftNum < rightNum), nil
	case scan.LESS_EQUAL:
		return scan.NewBooleanLoxValue(leftNum <= rightNum), nil

	case scan.BANG_EQUAL:
		return scan.NewBooleanLoxValue(leftNum != rightNum), nil
	case scan.EQUAL_EQUAL:
		return scan.NewBooleanLoxValue(leftNum == rightNum), nil

	case scan.MINUS:
		return scan.NewFloatLoxValue(leftNum - rightNum), nil
	case scan.SLASH:
		if rightNum == 0. {
			return nil,
				ConvertToRuntimeError(
					"evaluate expression error",
					fmt.Errorf("zero division error"),
					&expr.Operator,
				)
		}
		return scan.NewFloatLoxValue(leftNum / rightNum), nil
	case scan.STAR:
		return scan.NewFloatLoxValue(leftNum * rightNum), nil
	case scan.PLUS:
		return scan.NewFloatLoxValue(leftNum + rightNum), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *parse.GroupingExpression) (interface{}, error) {
	return i.Evaluate(expr.Expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *parse.UnaryExpression) (interface{}, error) {
	rightVal, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, ConvertToRuntimeError("evaluate expression error", err, &expr.Operator)
	}

	switch expr.Operator.Type {
	case scan.MINUS:
		numVal, err := rightVal.Number()
		if err != nil {
			return nil, ConvertToRuntimeError("evaluate expression error", err, &expr.Operator)
		}
		return scan.NewFloatLoxValue(-1. * numVal), nil
	case scan.BANG:
		boolVal := rightVal.Bool()
		return scan.NewBooleanLoxValue(!boolVal), nil
	}

	// unreachable
	return nil, nil
}

func (i *Interpreter) VisitStmtExpression(stmt *parse.StmtExpression) (interface{}, error) {
	return i.Evaluate(stmt.Expression)
}

func (i *Interpreter) VisitLogicalExpr(expr *parse.LogicalExpression) (interface{}, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == scan.OR {
		if left.Bool() {
			return left, nil
		} else {
			if !left.Bool() {
				return left, nil
			}
		}
	}

	return i.Evaluate(expr.Right)
}

func (i *Interpreter) VisitVariableExpr(expr *parse.VariableExpression) (interface{}, error) {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *parse.AssignExpression) (interface{}, error) {
	value, err := i.Evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	if err := i.environment.Assign(expr.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitCallExpr(expr *parse.CallExpression) (interface{}, error) {
	callee, err := i.Evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := make([]*scan.LoxValue, 0)
	for _, arg := range expr.Arguments {
		evalArg, err := i.Evaluate(arg)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, evalArg)
	}

	calleeFunc, err := callee.Callable()
	if err != nil {
		return nil, ConvertToRuntimeError("can only call functions or classes", err, &expr.Paren)
	}

	if len(arguments) != calleeFunc.Arity() {
		return nil, NewRuntimeError(
			fmt.Sprintf("expected %d arguments but got %d", calleeFunc.Arity(), len(arguments)),
			&expr.Paren,
		)
	}

	return calleeFunc.Call(arguments)
}
