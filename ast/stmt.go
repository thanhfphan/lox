package ast

type Stmt interface {
	Accept(v Visitor)
}

type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) Accept(v Visitor) {
	v.VisitPrintStmt(s)
}

type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) Accept(v Visitor) {
	v.VisitExpressionStmt(s)
}
