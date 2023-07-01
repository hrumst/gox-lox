package interpret

import "github.com/hrumst/gox-lox/lib/parse"

func (r *Resolver) VisitStmtExpression(stmt *parse.StmtExpression) (interface{}, error) {
	return nil, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitStmtPrint(stmt *parse.StmtPrint) (interface{}, error) {
	return nil, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitStmtVar(stmt *parse.StmtVar) (interface{}, error) {
	if err := r.declare(stmt.Name); err != nil {
		return nil, err
	}
	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) VisitStmtBlock(stmt *parse.StmtBlock) (interface{}, error) {
	r.beginScope()
	err := r.resolveStmts(stmt.Stmts)
	r.endScope()
	return nil, err
}

func (r *Resolver) VisitStmtIf(stmt *parse.StmtIf) (interface{}, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) VisitStmtWhile(stmt *parse.StmtWhile) (interface{}, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}
	return nil, r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitStmtExecuteControl(stmt *parse.StmtExecuteControl) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitStmtFunction(stmt *parse.StmtFunction) (interface{}, error) {
	if err := r.declare(stmt.Name); err != nil {
		return nil, err
	}
	r.define(stmt.Name)
	return nil, r.resolveFunction(stmt, inFunctionType)
}

func (r *Resolver) VisitStmtReturn(stmt *parse.StmtReturn) (interface{}, error) {
	if r.currentFuncType == noneFunctionType {
		return nil, NewRuntimeError("can't return from top-level code", &stmt.Keyword)
	}
	if stmt.Value != nil {
		return nil, r.resolveExpr(stmt.Value)
	}
	return nil, nil
}
