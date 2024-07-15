package ast

type Expr interface {
	Accept(v Visitor) any
}

type BinaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (e *BinaryExpr) Accept(v Visitor) any {
	return v.VisitBinaryExpr(e)
}

type UnaryExpr struct {
	Op   Token
	Expr Expr
}

func (e *UnaryExpr) Accept(v Visitor) any {
	return v.VisitUnaryExpr(e)
}

type LiteralExpr struct {
	Val any
}

func (e *LiteralExpr) Accept(v Visitor) any {
	return v.VisitLiteralExpr(e)
}

type GroupingExpr struct {
	Expr Expr
}

func (e *GroupingExpr) Accept(v Visitor) any {
	return v.VisitGroupingExpr(e)
}
