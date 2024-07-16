package ast

type Expr interface {
	Accept(v Visitor) any
}

// BinaryExpr ...
type BinaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (e *BinaryExpr) Accept(v Visitor) any {
	return v.VisitBinaryExpr(e)
}

// UnaryExpr ...
type UnaryExpr struct {
	Op   Token
	Expr Expr
}

func (e *UnaryExpr) Accept(v Visitor) any {
	return v.VisitUnaryExpr(e)
}

// LiteralExpr ...
type LiteralExpr struct {
	Val any
}

func (e *LiteralExpr) Accept(v Visitor) any {
	return v.VisitLiteralExpr(e)
}

// GroupingExpr ...
type GroupingExpr struct {
	Expr Expr
}

func (e *GroupingExpr) Accept(v Visitor) any {
	return v.VisitGroupingExpr(e)
}

// VariableExpr ...
type VariableExpr struct {
	Name *Token
}

func (e *VariableExpr) Accept(v Visitor) any {
	return v.VisitVariableExpr(e)
}

// AssignExpr ...
type AssignExpr struct {
	Name  *Token
	Value Expr
}

func (e *AssignExpr) Accept(v Visitor) any {
	return v.VisitAssignExpr(e)
}
