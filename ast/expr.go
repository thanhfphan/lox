package ast

type Expr interface {
	Accept(v ExprVisitor) any
}

// BinaryExpr ...
type BinaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(e)
}

// UnaryExpr ...
type UnaryExpr struct {
	Op    Token
	Right Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(e)
}

// LiteralExpr ...
type LiteralExpr struct {
	Val any
}

func (e *LiteralExpr) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(e)
}

// GroupingExpr ...
type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(e)
}

// VariableExpr ...
type VariableExpr struct {
	Name *Token
}

func (e *VariableExpr) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(e)
}

// AssignExpr ...
type AssignExpr struct {
	Name  *Token
	Value Expr
}

func (e *AssignExpr) Accept(v ExprVisitor) any {
	return v.VisitAssignExpr(e)
}

// LogicalExpr ...
type LogicalExpr struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func (s *LogicalExpr) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpr(s)
}

// CallExpr ...
type CallExpr struct {
	Callee    Expr
	Paren     *Token
	Arguments []Expr
}

func (s *CallExpr) Accept(v ExprVisitor) any {
	return v.VisitCallExpr(s)
}
