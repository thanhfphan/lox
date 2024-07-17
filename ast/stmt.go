package ast

type Stmt interface {
	Accept(v StmtVisitor)
}

// PrintStmt ...
type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) Accept(v StmtVisitor) {
	v.VisitPrintStmt(s)
}

// ExpressionStmt ...
type ExpressionStmt struct {
	Expression Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) {
	v.VisitExpressionStmt(s)
}

// VarStmt ...
type VarStmt struct {
	Name *Token
	Expr Expr
}

func (s *VarStmt) Accept(v StmtVisitor) {
	v.VisitVarStmt(s)
}

// BlockStmt ...
type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) Accept(v StmtVisitor) {
	v.VisitBlockStmt(s)
}
