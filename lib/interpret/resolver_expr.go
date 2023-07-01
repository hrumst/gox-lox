package interpret

import "github.com/hrumst/gox-lox/lib/parse"

func (r *Resolver) VisitBinaryExpr(expr *parse.BinaryExpression) (interface{}, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitGroupingExpr(expr *parse.GroupingExpression) (interface{}, error) {
	return nil, r.resolveExpr(expr.Expr)
}

func (r *Resolver) VisitLiteralExpr(expr *parse.LiteralExpression) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *parse.UnaryExpression) (interface{}, error) {
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *parse.VariableExpression) (interface{}, error) {
	if len(r.scopes) > 0 {
		if res, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && res == false {
			return nil, NewRuntimeError(
				"can't read local variable in its own initializer",
				&expr.Name,
			)
		}
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *parse.AssignExpression) (interface{}, error) {
	if err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *parse.LogicalExpression) (interface{}, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitCallExpr(expr *parse.CallExpression) (interface{}, error) {
	if err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}
	for _, arg := range expr.Arguments {
		if err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
