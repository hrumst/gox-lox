package interpret

import (
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
)

type functionType int

const (
	noneFunctionType functionType = iota
	inFunctionType
)

type Resolver struct {
	scopes          []map[string]bool
	interpreter     *Interpreter
	currentFuncType functionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      make([]map[string]bool, 0),
	}
}

func (r *Resolver) Resolve(stmts []parse.Statement) error {
	r.beginScope()
	if err := r.resolveStmts(stmts); err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) resolveStmts(stmts []parse.Statement) error {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt parse.Statement) error {
	_, err := stmt.Accept(r)
	return err
}

func (r *Resolver) resolveExpr(expr parse.Expression) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) resolveLocal(expr parse.Expression, name scan.Token) {
	for i := len(r.scopes) - 1; i >= 0; i -= 1 {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) declare(name scan.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		return NewRuntimeError("already variable with this name in this scope", &name)
	}

	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name scan.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveFunction(function *parse.StmtFunction, funcType functionType) error {
	enclosingFuncType := r.currentFuncType
	r.currentFuncType = funcType

	r.beginScope()
	for _, param := range function.Params {
		if err := r.declare(param); err != nil {
			return err
		}
		r.define(param)
	}
	if err := r.resolveStmts(function.Body); err != nil {
		return err
	}
	r.endScope()
	r.currentFuncType = enclosingFuncType
	return nil
}
