package parse

import "github.com/hrumst/gox-lox/lib/scan"

type ExpressionVisitor interface {
	VisitBinaryExpr(expr *BinaryExpression) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpression) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpression) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpression) (interface{}, error)
	VisitVariableExpr(expr *VariableExpression) (interface{}, error)
	VisitAssignExpr(expr *AssignExpression) (interface{}, error)
	VisitLogicalExpr(expr *LogicalExpression) (interface{}, error)
	VisitCallExpr(expr *CallExpression) (interface{}, error)
	VisitGetExpr(expr *GetExpression) (interface{}, error)
	VisitSetExpr(expr *SetExpression) (interface{}, error)
	VisitThisExpr(expr *ThisExpression) (interface{}, error)
	VisitSuperExpr(expr *SuperExpression) (interface{}, error)
}

type Expression interface {
	Accept(visitor ExpressionVisitor) (interface{}, error)
}

type BinaryExpression struct {
	Left     Expression
	Operator scan.Token
	Right    Expression
}

func NewBinaryExpression(left Expression, operator scan.Token, right Expression) *BinaryExpression {
	return &BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (be *BinaryExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(be)
}

type GroupingExpression struct {
	Expr Expression
}

func NewGroupingExpression(expr Expression) *GroupingExpression {
	return &GroupingExpression{
		Expr: expr,
	}
}

func (ge *GroupingExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(ge)
}

type LiteralExpression struct {
	Value *scan.Literal
}

func NewLiteralExpression(value *scan.Literal) *LiteralExpression {
	return &LiteralExpression{
		Value: value,
	}
}

func (le *LiteralExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(le)
}

type UnaryExpression struct {
	Operator scan.Token
	Right    Expression
}

func NewUnaryExpression(operator scan.Token, right Expression) *UnaryExpression {
	return &UnaryExpression{
		Operator: operator,
		Right:    right,
	}
}

func (ue *UnaryExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(ue)
}

type VariableExpression struct {
	Name scan.Token
}

func NewVariableExpression(name scan.Token) *VariableExpression {
	return &VariableExpression{
		Name: name,
	}
}

func (ve *VariableExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(ve)
}

type AssignExpression struct {
	Name  scan.Token
	Value Expression
}

func NewAssignExpression(name scan.Token, value Expression) *AssignExpression {
	return &AssignExpression{
		Name:  name,
		Value: value,
	}
}

func (ae *AssignExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitAssignExpr(ae)
}

type LogicalExpression struct {
	Left, Right Expression
	Operator    scan.Token
}

func NewLogicalExpression(
	left Expression,
	operator scan.Token,
	right Expression,
) *LogicalExpression {
	return &LogicalExpression{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (le *LogicalExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(le)
}

type CallExpression struct {
	Callee    Expression
	Paren     scan.Token
	Arguments []Expression
}

func NewCallExpression(callee Expression, paren scan.Token, arguments []Expression) *CallExpression {
	return &CallExpression{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (ce *CallExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitCallExpr(ce)
}

type GetExpression struct {
	Object Expression
	Name   scan.Token
}

func NewGetExpression(object Expression, name scan.Token) *GetExpression {
	return &GetExpression{
		Object: object,
		Name:   name,
	}
}

func (ge *GetExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitGetExpr(ge)
}

type SetExpression struct {
	Object Expression
	Name   scan.Token
	Value  Expression
}

func NewSetExpression(object Expression, name scan.Token, value Expression) *SetExpression {
	return &SetExpression{
		Object: object,
		Name:   name,
		Value:  value,
	}
}
func (se *SetExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitSetExpr(se)
}

type ThisExpression struct {
	Keyword scan.Token
}

func NewThisExpression(keyword scan.Token) *ThisExpression {
	return &ThisExpression{
		Keyword: keyword,
	}
}

func (te *ThisExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitThisExpr(te)
}

type SuperExpression struct {
	Keyword scan.Token
	Method  scan.Token
}

func NewSuperExpression(keyword scan.Token, method scan.Token) *SuperExpression {
	return &SuperExpression{
		Keyword: keyword,
		Method:  method,
	}
}

func (se *SuperExpression) Accept(visitor ExpressionVisitor) (interface{}, error) {
	return visitor.VisitSuperExpr(se)
}
