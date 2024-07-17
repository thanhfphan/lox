package ast

type ExprVisitor interface {
	VisitLiteralExpr(expr *LiteralExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitBinaryExpr(expr *BinaryExpr) any
	VisitVariableExpr(expr *VariableExpr) any
	VisitAssignExpr(expr *AssignExpr) any
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitVarStmt(stmt *VarStmt)
	VisitBlockStmt(stmt *BlockStmt)
	VisitIfStmt(stmt *IfStmt)
}
