package ast

type Expr interface {
}

type BinaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

type UnaryExpr struct {
	Op   Token
	Expr Expr
}

type LiteralExpr struct {
	Val any
}

type GroupingExpr struct {
	Expr Expr
}
