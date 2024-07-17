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
	Op   Token
	Expr Expr
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
	Expr Expr
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
