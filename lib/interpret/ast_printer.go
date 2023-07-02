package interpret

import (
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
	"strings"
)

type AstPrinter struct {
	isReverseNotation bool
}

func (v *AstPrinter) VisitSuperExpr(expr *parse.SuperExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (v *AstPrinter) VisitGetExpr(expr *parse.GetExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (v *AstPrinter) VisitSetExpr(expr *parse.SetExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (v *AstPrinter) VisitThisExpr(expr *parse.ThisExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (v *AstPrinter) VisitCallExpr(expr *parse.CallExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (v *AstPrinter) VisitLogicalExpr(expr *parse.LogicalExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (v *AstPrinter) VisitVariableExpr(expr *parse.VariableExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (v *AstPrinter) VisitAssignExpr(expr *parse.AssignExpression) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func NewAstPrinter(isReverseNotation bool) *AstPrinter {
	return &AstPrinter{
		isReverseNotation: isReverseNotation,
	}
}

func (v *AstPrinter) Print(expr parse.Expression) (string, error) {
	acpt, err := expr.Accept(v)
	if err != nil {
		return "", err
	}
	return acpt.(string), nil
}

func (v *AstPrinter) parenthesize(name string, expressions ...parse.Expression) (string, error) {
	var sb strings.Builder
	sb.WriteString("(")
	if !v.isReverseNotation {
		sb.WriteString(name)
	}
	for _, expr := range expressions {
		if !v.isReverseNotation {
			sb.WriteString(" ")
		}
		acpt, err := expr.Accept(v)
		if err != nil {
			return "", err
		}
		switch actt := acpt.(type) {
		case *scan.Literal:
			sb.WriteString(actt.Value.String())
		case string:
			sb.WriteString(actt)
		}

		if v.isReverseNotation {
			sb.WriteString(" ")
		}
	}
	if v.isReverseNotation {
		sb.WriteString(name)
	}
	sb.WriteString(")")
	return sb.String(), nil
}

func (v *AstPrinter) VisitBinaryExpr(expr *parse.BinaryExpression) (interface{}, error) {
	return v.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (v *AstPrinter) VisitGroupingExpr(expr *parse.GroupingExpression) (interface{}, error) {
	return v.parenthesize("group", expr.Expr)
}

func (v *AstPrinter) VisitLiteralExpr(expr *parse.LiteralExpression) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return expr.Value, nil
}

func (v *AstPrinter) VisitUnaryExpr(expr *parse.UnaryExpression) (interface{}, error) {
	return v.parenthesize(expr.Operator.Lexeme, expr.Right)
}
