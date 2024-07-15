package ast

type Visitor interface {
	VisitLiteralExpr(expr *LiteralExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitBinaryExpr(expr *BinaryExpr) any

	VisitPrintStmt(stmt *PrintStmt)
	VisitExpressionStmt(stmt *ExpressionStmt)
}
