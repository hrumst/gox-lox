package parse

import "github.com/hrumst/gox-lox/lib/scan"

type StatementInterpreter interface {
	VisitStmtExpression(stmt *StmtExpression) (interface{}, error)
	VisitStmtPrint(stmt *StmtPrint) (interface{}, error)
	VisitStmtVar(stmt *StmtVar) (interface{}, error)
	VisitStmtBlock(stmt *StmtBlock) (interface{}, error)
	VisitStmtIf(stmt *StmtIf) (interface{}, error)
	VisitStmtWhile(stmt *StmtWhile) (interface{}, error)
	VisitStmtExecuteControl(stmt *StmtExecuteControl) (interface{}, error)
	VisitStmtFunction(stmt *StmtFunction) (interface{}, error)
	VisitStmtReturn(stmt *StmtReturn) (interface{}, error)
	VisitStmtClass(stmt *StmtClass) (interface{}, error)
}

type Statement interface {
	Accept(interpreter StatementInterpreter) (interface{}, error)
}

type StmtExpression struct {
	Expression Expression
}

func NewStmtExpression(expr Expression) *StmtExpression {
	return &StmtExpression{
		Expression: expr,
	}
}

func (s *StmtExpression) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtExpression(s)
}

type StmtPrint struct {
	Expression Expression
}

func NewStmtPrint(expr Expression) *StmtPrint {
	return &StmtPrint{
		Expression: expr,
	}
}

func (s *StmtPrint) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtPrint(s)
}

type StmtVar struct {
	Name        scan.Token
	Initializer Expression
}

func NewStmtVar(name scan.Token, expression Expression) *StmtVar {
	return &StmtVar{Name: name, Initializer: expression}
}

func (s *StmtVar) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtVar(s)
}

type StmtBlock struct {
	Stmts []Statement
}

func NewStmtBlock(stmts []Statement) *StmtBlock {
	return &StmtBlock{
		Stmts: stmts,
	}
}

func (s *StmtBlock) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtBlock(s)
}

type StmtIf struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func NewStmtIf(
	condition Expression,
	thenBranch, elseBranch Statement,
) *StmtIf {
	return &StmtIf{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (s *StmtIf) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtIf(s)
}

type StmtWhile struct {
	Condition Expression
	Body      Statement
}

func NewStmtWhile(
	condition Expression,
	body Statement,
) *StmtWhile {
	return &StmtWhile{Condition: condition, Body: body}
}

func (s *StmtWhile) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtWhile(s)
}

type StmtExecuteControl struct {
	Control scan.Token
}

func NewExecuteControlStmt(control scan.Token) *StmtExecuteControl {
	return &StmtExecuteControl{
		Control: control,
	}
}

func (s *StmtExecuteControl) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtExecuteControl(s)
}

type StmtFunction struct {
	Name   scan.Token
	Params []scan.Token
	Body   []Statement
}

func NewStmtFunction(name scan.Token, params []scan.Token, body []Statement) *StmtFunction {
	return &StmtFunction{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (s *StmtFunction) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtFunction(s)
}

type StmtReturn struct {
	Keyword scan.Token
	Value   Expression
}

func NewStmtReturn(keyword scan.Token, value Expression) *StmtReturn {
	return &StmtReturn{
		Keyword: keyword,
		Value:   value,
	}
}

func (s *StmtReturn) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtReturn(s)
}

type StmtClass struct {
	Name       scan.Token
	Methods    []Statement
	SuperClass *VariableExpression
}

func NewStmtClass(name scan.Token, superClass *VariableExpression, methods []Statement) *StmtClass {
	return &StmtClass{
		Name:       name,
		Methods:    methods,
		SuperClass: superClass,
	}
}

func (s *StmtClass) Accept(interpreter StatementInterpreter) (interface{}, error) {
	return interpreter.VisitStmtClass(s)
}
