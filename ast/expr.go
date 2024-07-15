package ast

type ExprVisitor interface {
	VisitLiteralExpr(expr *LiteralExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitBinaryExpr(expr *BinaryExpr) any
}

type Expr interface {
	Accept(v ExprVisitor) any
}

type BinaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(e)
}

type UnaryExpr struct {
	Op   Token
	Expr Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(e)
}

type LiteralExpr struct {
	Val any
}

func (e *LiteralExpr) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(e)
}
