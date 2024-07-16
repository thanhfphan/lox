package ast

type Stmt interface {
	Accept(v Visitor)
}

// PrintStmt ...
type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) Accept(v Visitor) {
	v.VisitPrintStmt(s)
}

// ExpressionStmt ...
type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) Accept(v Visitor) {
	v.VisitExpressionStmt(s)
}

// VarStmt ...
type VarStmt struct {
	Name *Token
	Expr Expr
}

func (s *VarStmt) Accept(v Visitor) {
	v.VisitVarStmt(s)
}

// BlockStmt ...
type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) Accept(v Visitor) {
	v.VisitBlockStmt(s)
}
